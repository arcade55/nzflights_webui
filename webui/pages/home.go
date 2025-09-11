package pages

import (
	"github.com/arcade55/htma"
	"github.com/arcade55/nzflights-models"
	"github.com/arcade55/nzflights_webui/webui/components"
)

func HomePage(flights []nzflights.Flight) htma.Element {
	var flightCards []htma.Renderable
	for _, flight := range flights {
		// FIX: Wrap the 'flight' in a 'FlightValue' struct to match
		// the component's expected input type.
		flightValue := nzflights.FlightValue{
			ElementId: flight.Ident, // Set the ElementId for Datastar
			Flight:    flight,
		}
		flightCards = append(flightCards, components.FlightCardComponent(flightValue))
	}

	// The ID for this container must match the selector used in the SSE handler
	mainContent := htma.Div().IDAttr("flight-card-container").AddChild(flightCards...)

	return LayoutComponent("My Flights", mainContent)
}
