package sse

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/arcade55/nzflights-models"
	"github.com/arcade55/nzflights_webui/server/handlers/middleware"
	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats-server/v2/test"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// setupTestEnvironment creates a clean, isolated NATS environment for each test.
func setupTestEnvironment(t *testing.T) (jetstream.KeyValue, func()) {
	t.Helper()
	opts := &server.Options{Port: -1, JetStream: true}
	s := test.RunServer(opts)
	nc, err := nats.Connect(s.ClientURL())
	if err != nil {
		t.Fatalf("NATS connect failed: %v", err)
	}
	js, err := jetstream.New(nc)
	if err != nil {
		t.Fatalf("JetStream context failed: %v", err)
	}

	// Using a unique bucket name guarantees test isolation.
	bucketName := fmt.Sprintf("flights_%d", time.Now().UnixNano())
	kv, err := js.CreateKeyValue(context.Background(), jetstream.KeyValueConfig{Bucket: bucketName})
	if err != nil {
		t.Fatalf("KV creation failed: %v", err)
	}

	cleanup := func() {
		nc.Close()
		s.Shutdown()
	}
	return kv, cleanup
}

// readEvents reads a specific number of SSE events from a response within a timeout.
func readEvents(t *testing.T, scanner *bufio.Scanner, count int, timeout time.Duration) []string {
	t.Helper()
	var events []string
	eventChan := make(chan string)

	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "data: elements") {
				eventChan <- line
			}
		}
	}()

	for i := 0; i < count; i++ {
		select {
		case event := <-eventChan:
			events = append(events, event)
		case <-time.After(timeout):
			t.Errorf("timed out waiting for event %d of %d", i+1, count)
			return events
		}
	}
	return events
}

// TestFlightSSE_InitialState verifies that N events are sent for N initial flights.
func TestFlightSSE_InitialState(t *testing.T) {
	kv, cleanup := setupTestEnvironment(t)
	defer cleanup()

	const expectedUserID = "user123"
	const sessionToken = "a-valid-session-uuid-from-a-cookie"

	initialFlights := []nzflights.FlightValue{
		{ElementId: "NZ500", Flight: nzflights.Flight{Ident: "NZ500"}},
		{ElementId: "QF144", Flight: nzflights.Flight{Ident: "QF144"}},
	}
	for _, flight := range initialFlights {
		key := fmt.Sprintf("users.%s.flights.owned.%s", expectedUserID, flight.Flight.Ident)
		data, _ := json.Marshal(flight)
		kv.Put(context.Background(), key, data)
	}

	mockStore := middleware.NewMockSessionStore()
	protectedHandler := middleware.Auth(mockStore)(&FlightSSEHandler{KV: kv})
	server := httptest.NewServer(protectedHandler)
	defer server.Close()

	req, _ := http.NewRequest("GET", server.URL, nil)
	req.AddCookie(&http.Cookie{Name: "session_token", Value: sessionToken})

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer res.Body.Close()

	scanner := bufio.NewScanner(res.Body)
	events := readEvents(t, scanner, len(initialFlights), 1*time.Second)

	if len(events) != len(initialFlights) {
		t.Fatalf("Expected %d initial events, but got %d", len(initialFlights), len(events))
	}
	t.Logf("Successfully received %d initial SSE events.", len(events))
}

// TestFlightSSE_Unauthorized verifies the middleware blocks unauthorized requests.
func TestFlightSSE_Unauthorized(t *testing.T) {
	kv, cleanup := setupTestEnvironment(t)
	defer cleanup()

	mockStore := middleware.NewMockSessionStore()
	protectedHandler := middleware.Auth(mockStore)(&FlightSSEHandler{KV: kv})
	server := httptest.NewServer(protectedHandler)
	defer server.Close()

	// --- Request without cookie ---
	res, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status %d for request without cookie, but got %d", http.StatusUnauthorized, res.StatusCode)
	}
	t.Log("Successfully blocked request without cookie.")

	// --- Request with invalid cookie ---
	req, _ := http.NewRequest("GET", server.URL, nil)
	req.AddCookie(&http.Cookie{Name: "session_token", Value: "invalid-token"})
	res, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status %d for request with invalid cookie, but got %d", http.StatusUnauthorized, res.StatusCode)
	}
	t.Log("Successfully blocked request with invalid cookie.")
}
