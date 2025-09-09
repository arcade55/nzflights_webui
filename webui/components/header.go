package components

import "github.com/arcade55/htma"

// HeaderComponent creates the main application header.
func HeaderComponent() htma.Element {
	// The SVG content is defined as a string
	searchIconSVG := `<svg viewBox="0 0 24 24"><circle cx="11" cy="11" r="8"></circle><line x1="21" y1="21" x2="16.65" y2="16.65"></line></svg>`

	return htma.Header().ClassAttr("main-header").AddChild(
		htma.H1().Text("My Flights"),
		htma.Button().ClassAttr("header-button").AddChild(
			// Use RawContent to inject the SVG string without escaping
			htma.RawContent(searchIconSVG),
		),
	)
}
