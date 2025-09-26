package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/arcade55/logging"
	"github.com/arcade55/nzflights-models"
	"github.com/arcade55/nzflights_webui/server/handlers/sse"
	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

var (
	flightsMap map[string]nzflights.FlightValue
	cloudJS    jetstream.JetStream
	inMemoryJS jetstream.JetStream
	mirrorKV   jetstream.KeyValue
)

func TestMain(m *testing.M) {
	// Initialize logger
	logger, _, _ := logging.Init(context.Background(), logging.Config{
		Format:    logging.FormatPretty,
		Level:     logging.LevelDebug,
		AddSource: true,
	})
	log := logger.WithContext(context.Background())

	// Read credentials
	creds, err := credsFile.ReadFile("nats.cred")
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	// Connect to Synadia Cloud
	cloudNC, err := nats.Connect("tls://connect.ngs.global:4222", nats.UserCredentialBytes(creds))
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	defer cloudNC.Close()

	cloudJS, err = jetstream.New(cloudNC)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	// Setup in-memory NATS server
	ns, err := server.NewServer(&server.Options{
		JetStream: true,
	})
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	go ns.Start()
	if !ns.ReadyForConnections(10 * time.Second) {
		log.Error(fmt.Errorf("NATS server failed to start"))
		os.Exit(1)
	}
	defer ns.Shutdown()

	nc, err := nats.Connect(ns.ClientURL())
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	defer nc.Close()

	inMemoryJS, err = jetstream.New(nc)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	// Create in-memory KV store mirroring Synadia
	ctx := context.Background()
	sourceKV, err := cloudJS.KeyValue(ctx, "flights")
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	mirrorKV, err = inMemoryJS.CreateKeyValue(ctx, jetstream.KeyValueConfig{
		Bucket:  "flights",
		Storage: jetstream.MemoryStorage,
		Mirror: &jetstream.StreamSource{
			Name:   "flights",
			Domain: "ngs",
		},
	})
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	// Initialize map and watcher
	flightsMap = make(map[string]nzflights.FlightValue)
	watcher, err := mirrorKV.WatchAll(ctx)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	defer watcher.Stop()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		for entry := range watcher.Updates() {
			if entry == nil {
				continue
			}
			var fv nzflights.FlightValue
			if err := json.Unmarshal(entry.Value(), &fv); err != nil {
				log.Error(err)
				continue
			}
			flightsMap[entry.Key()] = fv
			log.Info("Updated flight in map", slog.String("key", entry.Key()))
			wg.Done()
		}
	}()

	// Copy existing data from source to mirror
	keys, err := sourceKV.Keys(ctx)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	for _, key := range keys {
		entry, err := sourceKV.Get(ctx, key)
		if err != nil {
			log.Error(err)
			continue
		}
		mirrorKV.Put(ctx, key, entry.Value())
	}

	wg.Wait()

	// Run tests
	exitCode := m.Run()

	os.Exit(exitCode)
}

func TestSearchFlights(t *testing.T) {
	h := &sse.SearchSSEHandler{Flights: flightsMap}

	reqBody := strings.NewReader(`{"searchTerm": "NZ"}`)
	req := httptest.NewRequest(http.MethodPost, "/search-flights", reqBody)
	w := httptest.NewRecorder()

	h.Search(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}

	if !strings.Contains(string(body), "NZ") {
		t.Errorf("expected response to contain 'NZ'; got %s", string(body))
	}
}
