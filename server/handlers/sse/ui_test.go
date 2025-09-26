//go:build manual_test

package sse

/*
	To run this: go test -v -tags=manual_test -run TestInteractiveUI . -manual

	Manual Testing Instructions:
	1. Run the test command above to start the server on http://localhost:8080.
	2. Open the URL in your browser.
	3. The home page should load with a flight list (updated via SSE as the simulator adds flights).
	4. To test search: Assuming the home page has a search input (bound to POST to /sse/search-flights via DataStar or similar),
	   type a partial flight identifier (e.g., "NZ", "QF", "JET", or "XXX") and observe the results rendering in #search-results.
	   The search queries the static flights loaded from testdata/scheduled_departures.json.
	5. Verify SSE updates: Watch the flight list update every ~10s as the simulator adds flights.
	6. Stop the server with Ctrl+C.
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

// generateSampleFlights creates a map with two representative FlightValue entries.
func generateSampleFlights() map[string]nzflights.FlightValue {
	flights := make(map[string]nzflights.FlightValue)

	// Flight 1: NZ123 from Auckland to Wellington
	f1 := nzflights.Flight{
		Ident:           "NZ123",
		IdentICAO:       "ANZ123",
		IdentIATA:       "NZ123",
		Operator:        "ANZ",
		FAFlightID:      "ANZ123-1756532266-airline-1234",
		Origin:          "NZAA",
		OriginIATA:      "AKL",
		OriginCity:      "Auckland",
		Destination:     "NZWN",
		DestinationIATA: "WLG",
		DestinationCity: "Wellington",
		AircraftType:    "A320",
		ScheduledOut:    "2025-09-15T09:00:00Z",
		ScheduledIn:     "2025-09-15T10:15:00Z",
		ActualOff:       "",
		ActualOn:        "",
		Status:          "Scheduled",
		GateOrigin:      "24",
		GateDestination: "12",
		Alerts:          nil,
	}
	flights["NZ123"] = nzflights.FlightValue{
		CorrelationID: "webhook-uuid-98765",
		NatsKey:       "users.test.flights.owned.NZ123",
		ElementId:     "NZ123",
		LastUpdated:   time.Now(),
		Flight:        f1,
	}

	// Flight 2: QF456 from Sydney to Christchurch
	f2 := nzflights.Flight{
		Ident:           "QF456",
		IdentICAO:       "QFA456",
		IdentIATA:       "QF456",
		Operator:        "QFA",
		FAFlightID:      "QFA456-1756536446-airline-5678",
		Origin:          "YSSY",
		OriginIATA:      "SYD",
		OriginCity:      "Sydney",
		Destination:     "NZCH",
		DestinationIATA: "CHC",
		DestinationCity: "Christchurch",
		AircraftType:    "B738",
		ScheduledOut:    "2025-09-15T12:30:00Z",
		ScheduledIn:     "2025-09-15T17:45:00Z",
		ActualOff:       "2025-09-15T12:35:00Z",
		ActualOn:        "",
		Status:          "En Route",
		GateOrigin:      "T1-15",
		GateDestination: "9",
		Alerts: []nzflights.Alert{
			{
				LongDescription:  "Minor departure delay due to air traffic control.",
				ShortDescription: "Delayed departure",
				Summary:          "Flight QF456 delayed by 5 minutes.",
				EventCode:        "DELAY",
				AlertID:          1001,
			},
		},
	}
	flights["QF456"] = nzflights.FlightValue{
		CorrelationID: "webhook-uuid-54321",
		NatsKey:       "users.test.flights.owned.QF456",
		ElementId:     "QF456",
		LastUpdated:   time.Now(),
		Flight:        f2,
	}

	return flights
}

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
			{ElementId: "NZ999", Flight: nzflights.Flight{Operator: "ANZ", Ident: "NZ999", OriginCity: "Wellington", DestinationCity: "Dunedin", Status: "Delayed"}},
			{ElementId: "QF200", Flight: nzflights.Flight{Operator: "QFA", Ident: "QF200", Origin: "Melbourne", Destination: "Christchurch", Status: "Boarding"}},
		}
		var flightIndex int

		time.Sleep(1 * time.Second)
		for {
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
				log.Info("SIMULATOR: Wrote update for flight", slog.String("flightID", flight.Flight.Ident), slog.String("status", flight.Flight.Status))
			}
			flightIndex++

			time.Sleep(1 * time.Second) // Stagger additions for easier observation.
		}
	}()

	// Set up the HTTP server mux.
	mux := http.NewServeMux()

	sseHandler := &FlightSSEHandler{KV: kv}
	mux.Handle("/sse/flights", testVisitorIDMiddleware(testUserID)(sseHandler))

	// Load flights for search
	flights := generateSampleFlights()
	searchHandler := &SearchSSEHandler{Flights: flights}
	mux.Handle("/search-flights", testVisitorIDMiddleware(testUserID)(http.HandlerFunc(searchHandler.Search)))

	// The Home page handler.
	homeHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pages.HomePage().RenderStream(w)
	})
	mux.Handle("/", testVisitorIDMiddleware(testUserID)(homeHandler))

	// Handler for static files.
	staticDir := http.Dir("../../../webui/static")
	log.Info("Static directory path", slog.String("path", string(staticDir)))
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
