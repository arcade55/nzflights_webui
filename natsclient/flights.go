package natsclient

import (
	"context"
	"errors"
	"sync"

	"github.com/nats-io/nats.go/jetstream"
)

// FlightStore defines the interface for all flight data operations.
// It abstracts away the underlying caching layers (in-memory, cloud).
type FlightStore interface {
	// GetMultiple retrieves the latest values for a slice of flight keys.
	// It automatically checks the fast in-memory cache first, then the cloud cache for each key.
	// Returns a map of keys to their found entries. Keys not found are omitted.
	GetMultiple(ctx context.Context, keys []string) (map[string]jetstream.KeyValueEntry, error)

	// WatchMultiple creates a unified watcher for a given slice of keys.
	// It intelligently merges updates from both the in-memory and cloud caches for all keys,
	// providing a single channel of updates to the caller.
	// The returned Watcher must be stopped by the caller when no longer needed.
	WatchMultiple(ctx context.Context, keys []string) (Watcher, error)

	// --- In-Memory Only Methods for Development ---

	// GetMultipleInMemory retrieves values only from the fast in-memory cache.
	GetMultipleInMemory(ctx context.Context, keys []string) (map[string]jetstream.KeyValueEntry, error)
	// WatchMultipleInMemory creates a watcher that only listens to the in-memory cache.
	WatchMultipleInMemory(ctx context.Context, keys []string) (Watcher, error)
}

// Watcher is a simplified interface for a Key-Value watcher.
type Watcher interface {
	Updates() <-chan jetstream.KeyValueEntry
	Stop()
}

// flightStore is the concrete implementation of our FlightStore interface.
type flightStore struct {
	inMemoryKV jetstream.KeyValue
	cloudKV    jetstream.KeyValue
}

// newFlightStore is a private constructor for our flight store.
func newFlightStore(inMemoryKV, cloudKV jetstream.KeyValue) FlightStore {
	return &flightStore{
		inMemoryKV: inMemoryKV,
		cloudKV:    cloudKV,
	}
}

// GetMultiple fetches multiple keys in parallel using a WaitGroup and a Mutex.
func (s *flightStore) GetMultiple(ctx context.Context, keys []string) (map[string]jetstream.KeyValueEntry, error) {
	results := make(map[string]jetstream.KeyValueEntry)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, key := range keys {
		wg.Add(1)
		go func(k string) {
			defer wg.Done()
			var entry jetstream.KeyValueEntry
			var err error

			// Try the fast in-memory mirror first.
			entry, err = s.inMemoryKV.Get(ctx, k)

			// If not there, check the cloud KV store.
			if err != nil {
				entry, err = s.cloudKV.Get(ctx, k)
			}

			// If we found an entry in either store, add it to the results map.
			if err == nil {
				mu.Lock()
				results[k] = entry
				mu.Unlock()
			}
		}(key)
	}

	// Wait for all the concurrent fetches to complete.
	wg.Wait()

	return results, nil
}

// WatchMultiple creates and manages watchers for multiple keys, merging their updates.
func (s *flightStore) WatchMultiple(ctx context.Context, keys []string) (Watcher, error) {
	if len(keys) == 0 {
		return nil, errors.New("WatchMultiple requires at least one key")
	}

	mergedUpdates := make(chan jetstream.KeyValueEntry, 64)
	done := make(chan struct{})
	var allWatchers []jetstream.KeyWatcher

	// Goroutine to forward updates from any watcher to the merged channel.
	forwarder := func(w jetstream.KeyWatcher) {
		for {
			select {
			case entry := <-w.Updates():
				if entry != nil {
					mergedUpdates <- entry
				}
			case <-done:
				return
			}
		}
	}

	for _, key := range keys {
		// Watch in-memory store for this key
		memWatcher, err := s.inMemoryKV.Watch(ctx, key, jetstream.IgnoreDeletes())
		if err != nil {
			// Stop any watchers we've already created
			for _, w := range allWatchers {
				w.Stop()
			}
			return nil, err
		}
		allWatchers = append(allWatchers, memWatcher)
		go forwarder(memWatcher)

		// Watch cloud store for this key
		cloudWatcher, err := s.cloudKV.Watch(ctx, key, jetstream.IgnoreDeletes())
		if err != nil {
			for _, w := range allWatchers {
				w.Stop()
			}
			return nil, err
		}
		allWatchers = append(allWatchers, cloudWatcher)
		go forwarder(cloudWatcher)
	}

	return &mergedWatcher{
		updates:     mergedUpdates,
		allWatchers: allWatchers,
		done:        done,
	}, nil
}

// --- In-Memory Only Implementations ---

// GetMultipleInMemory fetches multiple keys in parallel, only from the in-memory store.
func (s *flightStore) GetMultipleInMemory(ctx context.Context, keys []string) (map[string]jetstream.KeyValueEntry, error) {
	results := make(map[string]jetstream.KeyValueEntry)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, key := range keys {
		wg.Add(1)
		go func(k string) {
			defer wg.Done()
			// Only try the in-memory mirror.
			if entry, err := s.inMemoryKV.Get(ctx, k); err == nil {
				mu.Lock()
				results[k] = entry
				mu.Unlock()
			}
		}(key)
	}

	wg.Wait()
	return results, nil
}

// WatchMultipleInMemory creates watchers that only listen to the in-memory store.
func (s *flightStore) WatchMultipleInMemory(ctx context.Context, keys []string) (Watcher, error) {
	if len(keys) == 0 {
		return nil, errors.New("WatchMultipleInMemory requires at least one key")
	}

	mergedUpdates := make(chan jetstream.KeyValueEntry, 64)
	done := make(chan struct{})
	var allWatchers []jetstream.KeyWatcher

	forwarder := func(w jetstream.KeyWatcher) {
		for {
			select {
			case entry := <-w.Updates():
				if entry != nil {
					mergedUpdates <- entry
				}
			case <-done:
				return
			}
		}
	}

	for _, key := range keys {
		// Only watch the in-memory store for this key.
		memWatcher, err := s.inMemoryKV.Watch(ctx, key, jetstream.IgnoreDeletes())
		if err != nil {
			for _, w := range allWatchers {
				w.Stop()
			}
			return nil, err
		}
		allWatchers = append(allWatchers, memWatcher)
		go forwarder(memWatcher)
	}

	return &mergedWatcher{
		updates:     mergedUpdates,
		allWatchers: allWatchers,
		done:        done,
	}, nil
}

// mergedWatcher implements the Watcher interface for multiple keys.
type mergedWatcher struct {
	updates     <-chan jetstream.KeyValueEntry
	allWatchers []jetstream.KeyWatcher
	done        chan struct{}
}

func (w *mergedWatcher) Updates() <-chan jetstream.KeyValueEntry {
	return w.updates
}

func (w *mergedWatcher) Stop() {
	for _, watcher := range w.allWatchers {
		watcher.Stop()
	}
	close(w.done)
}
