/*
ACTION: Update this file.
PATH:   main.go
PURPOSE: To correctly implement the SSE pattern for partial updates,

	with separate handlers for full pages and SSE streams.
*/
package main

import (
	"log"
	"net/http"

	// 1. ADD THIS IMPORT BACK IN
	"github.com/starfederation/datastar-go/datastar"

	// These imports are needed for the SSE handlers
	"github.com/arcade55/htma"
	"github.com/arcade55/nzflights-models"
	"github.com/arcade55/nzflights_webui/webui/components"
	"github.com/arcade55/nzflights_webui/webui/pages"
)

func main() {
	mux := http.NewServeMux()

	staticFS := http.FileServer(http.Dir("./webui/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", staticFS))

	// --- Handlers for FULL page loads (these remain) ---
	mux.HandleFunc("GET /home", handleHome)
	mux.HandleFunc("GET /add-flight", handleAddFlight)
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/home", http.StatusMovedPermanently)
	})

	// --- 2. NEW: Handlers for Datastar SSE updates ---
	mux.HandleFunc("GET /home-sse", handleHomeSSE)
	mux.HandleFunc("GET /add-flight-sse", handleAddFlightSSE)

	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}

// handleHome serves the initial, full HTML page.
func handleHome(w http.ResponseWriter, r *http.Request) {
	flights := getFlightsFromStore()
	page := pages.HomePage(flights)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	page.RenderStream(w)
}

// handleAddFlight serves the initial, full HTML page.
func handleAddFlight(w http.ResponseWriter, r *http.Request) {
	page := pages.AddFlightPage()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	page.RenderStream(w)
}

// --- 3. NEW SSE HANDLERS ---

// handleHomeSSE handles Datastar requests for the home page content.
func handleHomeSSE(w http.ResponseWriter, r *http.Request) {
	sse := datastar.NewSSE(w, r)
	flights := []nzflights.Flight{
		{
			Ident:        "SSE",
			IdentIATA:    "NZ",
			Origin:       "AKL",
			Destination:  "LAX",
			ScheduledOut: "2025-12-01T07:40:00Z",
			ScheduledIn:  "2025-12-01T19:45:00Z",
			Status:       "Delayed",
			GateOrigin:   "B7",
		},
		{
			Ident:        "QF140",
			IdentIATA:    "QF",
			Origin:       "AKL",
			Destination:  "SYD",
			ScheduledOut: "2025-12-01T18:00:00Z",
			ScheduledIn:  "2025-12-01T21:45:00Z",
			Status:       "On Time",
			GateOrigin:   "A12",
		},
		{
			Ident:        "SSE",
			IdentIATA:    "NZ",
			Origin:       "AKL",
			Destination:  "LAX",
			ScheduledOut: "2025-12-01T07:40:00Z",
			ScheduledIn:  "2025-12-01T19:45:00Z",
			Status:       "Delayed",
			GateOrigin:   "B7",
		},
		{
			Ident:        "QF140",
			IdentIATA:    "QF",
			Origin:       "AKL",
			Destination:  "SYD",
			ScheduledOut: "2025-12-01T18:00:00Z",
			ScheduledIn:  "2025-12-01T21:45:00Z",
			Status:       "On Time",
			GateOrigin:   "A12",
		},
		{
			Ident:        "SSE",
			IdentIATA:    "NZ",
			Origin:       "AKL",
			Destination:  "LAX",
			ScheduledOut: "2025-12-01T07:40:00Z",
			ScheduledIn:  "2025-12-01T19:45:00Z",
			Status:       "Delayed",
			GateOrigin:   "B7",
		},
		{
			Ident:        "QF140",
			IdentIATA:    "QF",
			Origin:       "AKL",
			Destination:  "SYD",
			ScheduledOut: "2025-12-01T18:00:00Z",
			ScheduledIn:  "2025-12-01T21:45:00Z",
			Status:       "On Time",
			GateOrigin:   "A12",
		},
		{
			Ident:        "SSE",
			IdentIATA:    "NZ",
			Origin:       "AKL",
			Destination:  "LAX",
			ScheduledOut: "2025-12-01T07:40:00Z",
			ScheduledIn:  "2025-12-01T19:45:00Z",
			Status:       "Delayed",
			GateOrigin:   "B7",
		},
		{
			Ident:        "QF140",
			IdentIATA:    "QF",
			Origin:       "AKL",
			Destination:  "SYD",
			ScheduledOut: "2025-12-01T18:00:00Z",
			ScheduledIn:  "2025-12-01T21:45:00Z",
			Status:       "On Time",
			GateOrigin:   "A12",
		},
		{
			Ident:        "SSE",
			IdentIATA:    "NZ",
			Origin:       "AKL",
			Destination:  "LAX",
			ScheduledOut: "2025-12-01T07:40:00Z",
			ScheduledIn:  "2025-12-01T19:45:00Z",
			Status:       "Delayed",
			GateOrigin:   "B7",
		},
		{
			Ident:        "QF140",
			IdentIATA:    "QF",
			Origin:       "AKL",
			Destination:  "SYD",
			ScheduledOut: "2025-12-01T18:00:00Z",
			ScheduledIn:  "2025-12-01T21:45:00Z",
			Status:       "On Time",
			GateOrigin:   "A12",
		},
		{
			Ident:        "SSE",
			IdentIATA:    "NZ",
			Origin:       "AKL",
			Destination:  "LAX",
			ScheduledOut: "2025-12-01T07:40:00Z",
			ScheduledIn:  "2025-12-01T19:45:00Z",
			Status:       "Delayed",
			GateOrigin:   "B7",
		},
		{
			Ident:        "QF140",
			IdentIATA:    "QF",
			Origin:       "AKL",
			Destination:  "SYD",
			ScheduledOut: "2025-12-01T18:00:00Z",
			ScheduledIn:  "2025-12-01T21:45:00Z",
			Status:       "On Time",
			GateOrigin:   "A12",
		},
		{
			Ident:        "SSE",
			IdentIATA:    "NZ",
			Origin:       "AKL",
			Destination:  "LAX",
			ScheduledOut: "2025-12-01T07:40:00Z",
			ScheduledIn:  "2025-12-01T19:45:00Z",
			Status:       "Delayed",
			GateOrigin:   "B7",
		},
		{
			Ident:        "QF140",
			IdentIATA:    "QF",
			Origin:       "AKL",
			Destination:  "SYD",
			ScheduledOut: "2025-12-01T18:00:00Z",
			ScheduledIn:  "2025-12-01T21:45:00Z",
			Status:       "On Time",
			GateOrigin:   "A12",
		},
	}

	var flightCards []htma.Renderable
	for _, flight := range flights {
		flightCards = append(flightCards, components.FlightCardComponent(flight))
	}
	content := htma.Div().ClassAttr("flight-card-container").AddChild(flightCards...)

	// Patch the rendered HTML into the #main-content target's innerHTML
	sse.PatchElements(content.Render(), datastar.WithSelector("#main-content"), datastar.WithModeInner())
}

// handleAddFlightSSE handles Datastar requests for the add flight page content.
func handleAddFlightSSE(w http.ResponseWriter, r *http.Request) {
	sse := datastar.NewSSE(w, r)
	content := htma.Div().ClassAttr("search-container").AddChild(
		htma.Div().ClassAttr("ticket-card search-ticket").AddChild(
			components.InputField("airplane_ticket", "", "", "Optional Flight Number"),
			components.Separator("Or"),
			components.InputField("flight_takeoff", "From", "Auckland, NZ", ""),
			components.InputField("flight_land", "", "", "Select destination"),
			components.InputField("calendar_month", "Departure", "Sep 01, 2025", ""),
			components.ActionButton("SSE Search Flights"),
		),
	)

	sse.PatchElements(content.Render(), datastar.WithSelector("#main-content"), datastar.WithModeInner())
}

// getFlightsFromStore remains the same
func getFlightsFromStore() []nzflights.Flight {
	return []nzflights.Flight{
		{
			Ident:        "NZ4",
			IdentIATA:    "NZ",
			Origin:       "AKL",
			Destination:  "LAX",
			ScheduledOut: "2025-12-01T07:40:00Z",
			ScheduledIn:  "2025-12-01T19:45:00Z",
			Status:       "Delayed",
			GateOrigin:   "B7",
		},
		{
			Ident:        "QF140",
			IdentIATA:    "QF",
			Origin:       "AKL",
			Destination:  "SYD",
			ScheduledOut: "2025-12-01T18:00:00Z",
			ScheduledIn:  "2025-12-01T21:45:00Z",
			Status:       "On Time",
			GateOrigin:   "A12",
		},
	}
}
