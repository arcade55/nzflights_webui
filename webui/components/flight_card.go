/*
ACTION: Update this file.
PATH:   webui/components/flight_card.go
PURPOSE: To add the new, shared "ticket-card" class to the component's

	main Div element. This applies the parent styles.
*/
package components

import (
	"fmt"
	"strings"
	"time"

	"github.com/arcade55/htma"
	"github.com/arcade55/nzflights-models"
)

func FlightCardComponent(flight nzflights.Flight) htma.Element {
	// CHANGE THIS LINE: Add "ticket-card" to the class list.
	return htma.Div().ClassAttr("ticket-card flight-card").AddChild(
		// Card Header
		htma.Div().ClassAttr("card-header").AddChild(
			htma.Div().ClassAttr("airline-info").AddChild(
				htma.Div().
					ClassAttr(fmt.Sprintf("airline-logo airline-%s", strings.ToLower(flight.IdentIATA))).
					Text(flight.IdentIATA),
				htma.Div().AddChild(
					htma.Div().ClassAttr("flight-number").Text(flight.Ident),
					htma.Div().ClassAttr("airline-name").Text(getAirlineName(flight.IdentIATA)),
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
				htma.P().Text(getAirportCity(flight.Origin)),
			),
			htma.Div().ClassAttr("path-icon").AddChild(
				htma.Span().ClassAttr("material-symbols-outlined").Text("east"),
			),
			htma.Div().ClassAttr("location").AddChild(
				htma.H2().Text(flight.Destination),
				htma.P().Text(getAirportCity(flight.Destination)),
			),
		),

		// Flight Details
		htma.Div().ClassAttr("flight-details").AddChild(
			htma.Div().ClassAttr("detail-item").AddChild(
				htma.H3().Text(flight.GateOrigin),
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

// --- Helper Functions ---
func createStatusSpan(status string) htma.Element {
	span := htma.Span().Text(status)
	switch strings.ToLower(status) {
	case "delayed":
		return span.ClassAttr("status-delayed")
	case "on time":
		return span.ClassAttr("status-ontime")
	default:
		return span
	}
}

func formatTime(isoString string) string {
	t, err := time.Parse(time.RFC3339, isoString)
	if err != nil {
		return "??:??"
	}
	return t.Format("03:04 PM")
}

func getAirlineName(code string) string {
	names := map[string]string{"NZ": "Air New Zealand", "MU": "China Eastern", "QF": "Qantas"}
	if name, ok := names[code]; ok {
		return name
	}
	return "Airline"
}

func getAirportCity(code string) string {
	cities := map[string]string{"AKL": "Auckland", "LAX": "Los Angeles", "PVG": "Shanghai", "SYD": "Sydney"}
	if city, ok := cities[code]; ok {
		return city
	}
	return "Location"
}
