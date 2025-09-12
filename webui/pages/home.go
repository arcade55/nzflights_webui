package pages

import (
	"github.com/arcade55/htma"
)

func HomePage() htma.Element {

	// The ID for this container must match the selector used in the SSE handler
	mainContent := htma.Div().ClassAttr("flight-card-container").IDAttr("flights")

	return LayoutComponent("My Flights", mainContent)
}
