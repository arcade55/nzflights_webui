package components

import (
	"fmt"
	"strings"
	"time"

	"github.com/arcade55/htma"
	"github.com/arcade55/nzflights-models"
)

// FlightCardComponent now renders HTML that perfectly matches the static index.html prototype.
func FlightCardComponent(flightValue nzflights.FlightValue) htma.Element {
	flight := flightValue.Flight
	return htma.Div().
		IDAttr(flightValue.ElementId).
		ClassAttr(fmt.Sprintf("%s ticket-card flight-card", flightValue.ElementId)).
		AddChild(
			// --- Card Header ---
			htma.Div().ClassAttr("card-header").AddChild(
				htma.Div().ClassAttr("airline-info").AddChild(
					htma.Div().ClassAttr(fmt.Sprintf("airline-logo airline-%s", strings.ToLower(flight.IdentIATA))).Text(flight.IdentIATA),
					htma.Div().AddChild(
						htma.Div().ClassAttr("flight-number").Text(flight.Ident),
						htma.Div().ClassAttr("airline-name").Text(getAirlineName(flight.IdentIATA)),
					),
				),
				htma.Button().ClassAttr("card-share-button").AddChild(
					htma.Span().ClassAttr("material-symbols-outlined").Text("share"),
				),
			),

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

			// --- Card Footer ---
			htma.Div().ClassAttr("flight-card-container").AddChild(
				htma.Span().Text(formatTime(flight.ScheduledOut)),
				createStatusSpan(flight.Status),
				htma.Span().Text(formatTime(flight.ScheduledIn)),
			),
		)
}

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

func getAirlineName(iata string) string {
	// In a real app, this would be more robust
	if iata == "NZ" {
		return "Air New Zealand"
	}
	return "Airline"
}
