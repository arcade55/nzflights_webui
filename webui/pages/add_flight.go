package pages

import "github.com/arcade55/htma"

// AddFlightPage composes the view for the "Add Flight" page.
func AddFlightPage() htma.Element {
	// 1. Build this page's unique content (e.g., a form).
	mainContent := htma.Div().ClassAttr("form-container").AddChild(
		htma.H1().Text("Add a New Flight"),
		// ... your form elements would go here ...
		htma.P().Text("Flight addition form will be here."),
	)

	// 2. Wrap it with the same layout.
	return LayoutComponent("Add Flight", mainContent)
}
