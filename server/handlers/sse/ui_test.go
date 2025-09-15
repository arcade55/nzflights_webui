//go:build manual_test

package sse

/*
	To run this :   go test -v -tags=manual_test -run TestInteractiveUI . -manual

*/

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/arcade55/logging"
	"github.com/arcade55/nzflights-models"
	"github.com/arcade55/nzflights_webui/webui/pages"
	"github.com/google/uuid"
)

// Define a flag to ensure this test only runs when you explicitly ask for it.
var manualTest = flag.Bool("manual", false, "run manual, long-running tests")

// testVisitorIDMiddleware always sets the SAME visitor ID for predictable test keys.
func testVisitorIDMiddleware(userID string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie := &http.Cookie{
				Name:  "visitor_id",
				Value: userID,
				Path:  "/",
			}
			http.SetCookie(w, cookie)
			r.AddCookie(cookie)
			next.ServeHTTP(w, r)
		})
	}
}

// TestInteractiveUI starts a full web server for manual browser testing.
func TestInteractiveUI(t *testing.T) {
	if !*manualTest {
		t.Skip("Skipping manual test. To run, use: go test -v -tags=manual_test -run TestInteractiveUI . -manual")
	}

	// Initialize the logger
	logger, _, _ := logging.Init(context.Background(), logging.Config{
		Format:    logging.FormatPretty,
		Level:     logging.LevelDebug,
		AddSource: true,
	})
	log := logger.WithContext(context.Background())
	log.Info("Starting interactive UI test server...")

	// Set up the isolated NATS environment and a predictable User ID.
	kv, cleanup := setupTestEnvironment(t)
	defer cleanup()
	testUserID := uuid.NewString()
	log.Info("Test user ID generated", slog.String("userID", testUserID))

	// Backend Simulator Goroutine
	go func() {
		flightsToSimulate := []nzflights.FlightValue{
			{ElementId: "NZ527", Flight: nzflights.Flight{Operator: "ANZ", Ident: "NZ527", OriginCity: "Auckland", DestinationCity: "Wellington", Status: "Scheduled"}},
			{ElementId: "QF144", Flight: nzflights.Flight{Operator: "JST", Ident: "QF144", Origin: "Sydney", Destination: "Auckland", Status: "En Route"}},
			{ElementId: "JET345", Flight: nzflights.Flight{Ident: "JET345", Origin: "Christchurch", Destination: "Queenstown", Status: "Landed"}},
			{ElementId: "XXX", Flight: nzflights.Flight{Ident: "XXX", Origin: "Christchurch", Destination: "Queenstown", Status: "Landed"}},
		}
		var flightIndex int

		time.Sleep(5 * time.Second)
		for {

			time.Sleep(10 * time.Second)
			if flightIndex >= len(flightsToSimulate) {
				log.Info("SIMULATOR: All flights have been added.")
				return
			}

			flight := flightsToSimulate[flightIndex]
			key := fmt.Sprintf("users.%s.flights.owned.%s", testUserID, flight.Flight.Ident)
			data, _ := json.Marshal(flight)

			if _, err := kv.Put(context.Background(), key, data); err != nil {
				log.Error(fmt.Errorf("SIMULATOR: failed to put flight: %w", err))
			} else {
				log.Info("SIMULATOR: Wrote update for flight", slog.String("flightID", flight.Flight.Ident))
			}
			flightIndex++
		}
	}()

	// Set up the HTTP server mux.
	mux := http.NewServeMux()

	sseHandler := &FlightSSEHandler{KV: kv}
	mux.Handle("/sse/flights", testVisitorIDMiddleware(testUserID)(sseHandler))

	// Load flights for search
	flights, err := loadFlightsForSearch("testdata/scheduled_departures.json")
	if err != nil {
		t.Fatalf("Failed to load flights for search: %v", err)
	}
	searchHandler := &SearchSSEHandler{Flights: flights}
	mux.Handle("/sse/search-flights", testVisitorIDMiddleware(testUserID)(searchHandler))

	// The Home page handler.
	homeHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pages.HomePage().RenderStream(w)
	})
	mux.Handle("/", testVisitorIDMiddleware(testUserID)(homeHandler))

	// Handler for static files.
	staticDir := http.Dir("../../static")
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(staticDir)))

	// Start the server on port 8080.
	port := "8080"
	log.Info("🚀 Server is running! Open this URL in your browser", slog.String("url", "http://localhost:"+port))
	log.Warn("Press Ctrl+C in the terminal to stop the server.")

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Error(fmt.Errorf("server failed to start: %w", err))
		t.FailNow()
	}
}

type ScheduledDepartures struct {
	ScheduledDepartures []nzflights.Flight `json:"scheduled_departures"`
}

func loadFlightsForSearch(path string) (map[string]nzflights.FlightValue, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var departures ScheduledDepartures
	if err := json.Unmarshal(file, &departures); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	flights := make(map[string]nzflights.FlightValue)
	for _, f := range departures.ScheduledDepartures {
		flights[f.Ident] = nzflights.FlightValue{
			ElementId: f.Ident,
			Flight:    f,
		}
	}
	return flights, nil
}
