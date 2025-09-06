package pages

import (
	"github.com/arcade55/htma"
	"github.com/arcade55/nzflights-models"
	"github.com/arcade55/nzflights_webui/webui/components"
)

// HomePage composes the entire view for the home page.
// It takes the necessary data (a slice of flights) to render its content.
func HomePage(flights []nzflights.Flight) htma.Element {
	// 1. Build the page-specific main content.
	var flightCards []htma.Renderable
	for _, flight := range flights {
		flightCards = append(flightCards, components.FlightCardComponent(flight))
	}
	mainContent := htma.Div().AddChild(flightCards...)

	// 2. Wrap the main content with the common layout.
	return LayoutComponent("My Flights", mainContent)
}
