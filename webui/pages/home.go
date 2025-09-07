package pages

import (
	"github.com/arcade55/htma"
	"github.com/arcade55/nzflights-models"
	"github.com/arcade55/nzflights_webui/webui/components"
)

func HomePage(flights []nzflights.Flight) htma.Element {
	var flightCards []htma.Renderable
	for _, flight := range flights {
		flightCards = append(flightCards, components.FlightCardComponent(flight))
	}

	// Add the new class to this wrapper Div
	mainContent := htma.Div().ClassAttr("flight-card-container").AddChild(flightCards...)

	return LayoutComponent("My Flights", mainContent)
}
