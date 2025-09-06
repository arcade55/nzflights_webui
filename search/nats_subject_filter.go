package search

/*
import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/nats-io/nats.go/jetstream"
)

// --- NEW: KeyFilterBuilder for creating readable and error-free filters ---

const (
	// These constants map to the token positions in your key structure.
	// flights.{fa_flight_id}.{year}.{month}.{day}.{hour}.{min}.{sec}.{ident}.{origin}.{dest}
	tokenCount         = 11
	tokenIdxFAFlightID = 1
	tokenIdxYear       = 2
	tokenIdxMonth      = 3
	tokenIdxDay        = 4
	tokenIdxHour       = 5
	tokenIdxMinute     = 6
	tokenIdxSecond     = 7
	tokenIdxIdent      = 8
	tokenIdxOrigin     = 9
	tokenIdxDest       = 10
)

// KeyFilterBuilder provides a fluent interface to build NATS KV filter strings.
type KeyFilterBuilder struct {
	tokens []string
}

// NewFilterBuilder initializes a builder for the 'flights' bucket.
func NewFilterBuilder() *KeyFilterBuilder {
	tokens := make([]string, tokenCount)
	tokens[0] = "flights"
	for i := 1; i < tokenCount; i++ {
		tokens[i] = "*" // Default to wildcard
	}
	return &KeyFilterBuilder{tokens: tokens}
}

// WithDate sets the year, month, and day for the filter.
func (b *KeyFilterBuilder) WithDate(year int, month time.Month, day int) *KeyFilterBuilder {
	b.tokens[tokenIdxYear] = fmt.Sprintf("%04d", year)
	b.tokens[tokenIdxMonth] = fmt.Sprintf("%02d", month)
	b.tokens[tokenIdxDay] = fmt.Sprintf("%02d", day)
	return b
}

// WithHour sets the hour for the filter. It requires WithDate to be set first.
func (b *KeyFilterBuilder) WithHour(hour int) *KeyFilterBuilder {
	b.tokens[tokenIdxHour] = fmt.Sprintf("%02d", hour)
	return b
}

// WithIdent sets the airline identifier (e.g., "ANZ5023").
func (b *KeyFilterBuilder) WithIdent(ident string) *KeyFilterBuilder {
	b.tokens[tokenIdxIdent] = ident
	return b
}

// WithRoute sets both the origin and destination.
func (b *KeyFilterBuilder) WithRoute(origin, destination string) *KeyFilterBuilder {
	b.tokens[tokenIdxOrigin] = origin
	b.tokens[tokenIdxDest] = destination
	return b
}

// WithDestination sets the destination airport.
func (b *KeyFilterBuilder) WithDestination(destination string) *KeyFilterBuilder {
	b.tokens[tokenIdxDest] = destination
	return b
}

// Build constructs the final filter string.
// Use 'useTrailingWildcard' as 'true' for date/time ranges (e.g., all flights in a day).
func (b *KeyFilterBuilder) Build(useTrailingWildcard bool) string {
	if useTrailingWildcard {
		return strings.Join(b.tokens, ".") + ".>"
	}
	// Trim trailing wildcards for more specific queries like destination only
	lastTokenIndex := tokenCount - 1
	for lastTokenIndex > 0 && b.tokens[lastTokenIndex] == "*" {
		lastTokenIndex--
	}

	// If everything is a wildcard, we need to add a > to get all flights
	if lastTokenIndex == 0 {
		return "flights.>"
	}

	return strings.Join(b.tokens[:lastTokenIndex+1], ".")
}

// --- UPDATED: Filter functions using the new KeyFilterBuilder ---

// FilterByDate retrieves all flight keys for a specific year, month, and day.
func FilterByDate(ctx context.Context, kv jetstream.KeyValue, year int, month time.Month, day int) ([]string, error) {
	filter := NewFilterBuilder().WithDate(year, month, day).Build(true)
	log.Printf("Using filter: %s", filter)
	return kv.Keys(ctx, jetstream.KeysFilter(filter))
}

// FilterByHour retrieves all flight keys for a specific hour on a specific date.
func FilterByHour(ctx context.Context, kv jetstream.KeyValue, year int, month time.Month, day int, hour int) ([]string, error) {
	filter := NewFilterBuilder().WithDate(year, month, day).WithHour(hour).Build(true)
	log.Printf("Using filter: %s", filter)
	return kv.Keys(ctx, jetstream.KeysFilter(filter))
}

// FilterByAirline retrieves all flight keys for a specific airline identifier.
func FilterByAirline(ctx context.Context, kv jetstream.KeyValue, ident string) ([]string, error) {
	filter := NewFilterBuilder().WithIdent(ident).Build(false)
	log.Printf("Using filter: %s", filter)
	return kv.Keys(ctx, jetstream.KeysFilter(filter))
}

// FilterByRoute retrieves all flight keys for a specific route.
func FilterByRoute(ctx context.Context, kv jetstream.KeyValue, origin, destination string) ([]string, error) {
	filter := NewFilterBuilder().WithRoute(origin, destination).Build(false)
	log.Printf("Using filter: %s", filter)
	return kv.Keys(ctx, jetstream.KeysFilter(filter))
}

// FilterByDestination retrieves all flight keys to a specific destination.
func FilterByDestination(ctx context.Context, kv jetstream.KeyValue, destination string) ([]string, error) {
	filter := NewFilterBuilder().WithDestination(destination).Build(false)
	log.Printf("Using filter: %s", filter)
	return kv.Keys(ctx, jetstream.KeysFilter(filter))
}

// FilterByAirlineAndDate retrieves flights for a specific airline on a given date.
func FilterByAirlineAndDate(ctx context.Context, kv jetstream.KeyValue, ident string, year int, month time.Month, day int) ([]string, error) {
	filter := NewFilterBuilder().WithDate(year, month, day).WithIdent(ident).Build(true)
	log.Printf("Using filter: %s", filter)
	return kv.Keys(ctx, jetstream.KeysFilter(filter))
}
*/
