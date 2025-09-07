/*
ACTION: Update this file.
PATH:   webui/pages/add_flight.go
PURPOSE: To revert this component to its original state, where it only

	builds the full HTML page.
*/
package pages

import (
	"github.com/arcade55/htma"
	"github.com/arcade55/nzflights_webui/webui/components"
)

// This function now returns the full page layout directly.
func AddFlightPage() htma.Element {
	mainContent := htma.Div().ClassAttr("search-container").AddChild(
		htma.Div().ClassAttr("ticket-card search-ticket").AddChild(
			components.InputField("airplane_ticket", "", "", "Optional Flight Number"),
			components.Separator("Or"),
			components.InputField("flight_takeoff", "From", "Auckland, NZ", ""),
			components.InputField("flight_land", "", "", "Select destination"),
			components.InputField("calendar_month", "Departure", "Sep 01, 2025", ""),
			components.ActionButton("Search Flights"),
		),
	)

	return LayoutComponent("Add Flight", mainContent)
}
