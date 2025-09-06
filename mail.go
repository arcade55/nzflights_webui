package main

import (
	"log"
	"net/http"

	// Import your internal packages
	"github.com/arcade55/nzflights-models"
	"github.com/arcade55/nzflights_webui/webui/pages"
)

func main() {
	mux := http.NewServeMux()

	// Handle static files
	staticFS := http.FileServer(http.Dir("./webui/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", staticFS))

	// Register page handlers
	mux.HandleFunc("/", handleHome)
	mux.HandleFunc("/add-flight", handleAddFlight)

	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}

// handleHome is the HTTP handler for the home page.
func handleHome(w http.ResponseWriter, r *http.Request) {
	// 1. Get data.
	flights := getFlightsFromStore() // Placeholder for your data access logic

	// 2. Call the page composition function from the 'pages' package.
	page := pages.HomePage(flights)

	// 3. Render.
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	page.RenderStream(w)
}

// handleAddFlight is the HTTP handler for the "add flight" page.
func handleAddFlight(w http.ResponseWriter, r *http.Request) {
	page := pages.AddFlightPage()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	page.RenderStream(w)
}

// getFlightsFromStore is a placeholder function. In a real application,
// this would fetch data from your canonical store (e.g., database, NATS JetStream).
func getFlightsFromStore() []nzflights.Flight {
	// Returning mock data for demonstration purposes
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
