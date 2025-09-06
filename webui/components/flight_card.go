package components

import (
	"fmt"
	"strings"
	"time"

	"github.com/arcade55/htma"
	"github.com/arcade55/nzflights-models"
)

// FlightCardComponent generates an HTML flight card from a Flight struct.
// This version assumes the flight struct is fully populated with no nil pointers.
func FlightCardComponent(flight nzflights.Flight) htma.Element {
	// With guaranteed data, we no longer need helper variables for nil checks.
	// We can now access all fields directly.
	return htma.Div().ClassAttr("flight-card").AddChild(
		// Card Header
		htma.Div().ClassAttr("card-header").AddChild(
			htma.Div().ClassAttr("airline-info").AddChild(
				htma.Div().
					ClassAttr(fmt.Sprintf("airline-logo airline-%s", strings.ToLower(flight.IdentIATA))).
					Text(flight.IdentIATA),
				htma.Div().AddChild(
					htma.Div().ClassAttr("flight-number").Text(flight.Ident),
					htma.Div().ClassAttr("airline-name").Text(getAirlineName(flight.IdentIATA)), // Placeholder lookup
				),
			),
			htma.Button().ClassAttr("card-share-button").AddChild(
				htma.Span().ClassAttr("material-symbols-outlined").Text("share"),
			),
		),

		htma.Hr().ClassAttr("card-separator"),

		// Flight Path
		htma.Div().ClassAttr("flight-path").AddChild(
			htma.Div().ClassAttr("location").AddChild(
				htma.H2().Text(flight.Origin),
				htma.P().Text(getAirportCity(flight.Origin)), // Placeholder lookup
			),
			htma.Div().ClassAttr("path-icon").AddChild(
				htma.Span().ClassAttr("material-symbols-outlined").Text("east"),
			),
			htma.Div().ClassAttr("location").AddChild(
				htma.H2().Text(flight.Destination),
				htma.P().Text(getAirportCity(flight.Destination)), // Placeholder lookup
			),
		),

		// Flight Details
		htma.Div().ClassAttr("flight-details").AddChild(
			htma.Div().ClassAttr("detail-item").AddChild(
				htma.H3().Text(flight.GateOrigin), // Direct access, no nil check needed
				htma.P().Text("Gate"),
			),
			htma.Div().ClassAttr("detail-item").AddChild(
				htma.H3().Text(formatTime(flight.ScheduledOut)),
				htma.P().Text("Boarding"),
			),
		),

		htma.Hr().ClassAttr("card-separator"),

		// Card Footer
		htma.Div().ClassAttr("card-footer").AddChild(
			htma.Span().Text(formatTime(flight.ScheduledOut)),
			createStatusSpan(flight.Status),
			htma.Span().Text(formatTime(flight.ScheduledIn)),
		),
	)
}

// --- Helper Functions (remain the same) ---

// createStatusSpan conditionally creates the status span with the correct class.
func createStatusSpan(status string) htma.Element {
	span := htma.Span().Text(status)
	switch strings.ToLower(status) {
	case "delayed":
		return span.ClassAttr("status-delayed")
	case "on time":
		return span.ClassAttr("status-ontime")
	default:
		return span // No special class
	}
}

// formatTime parses an ISO 8601 string and returns it in "03:04 PM" format.
func formatTime(isoString string) string {
	t, err := time.Parse(time.RFC3339, isoString)
	if err != nil {
		return "??:??"
	}
	return t.Format("03:04 PM")
}

// getAirlineName is a placeholder for a data lookup.
func getAirlineName(code string) string {
	names := map[string]string{"NZ": "Air New Zealand", "MU": "China Eastern", "QF": "Qantas"}
	if name, ok := names[code]; ok {
		return name
	}
	return "Airline"
}

// getAirportCity is a placeholder for a data lookup.
func getAirportCity(code string) string {
	cities := map[string]string{"AKL": "Auckland", "LAX": "Los Angeles", "PVG": "Shanghai", "SYD": "Sydney"}
	if city, ok := cities[code]; ok {
		return city
	}
	return "Location"
}
