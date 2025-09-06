package components

import "github.com/arcade55/htma"

// FooterComponent creates the main application footer with navigation buttons.
func FooterComponent() htma.Element {
	// Define SVG content for each icon
	homeIconSVG := `<svg viewBox="0 0 24 24"><path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"></path><polyline points="9 22 9 12 15 12 15 22"></polyline></svg>`
	sharedIconSVG := `<svg viewBox="0 0 24 24"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path><circle cx="9" cy="7" r="4"></circle><path d="M23 21v-2a4 4 0 0 0-3-3.87"></path><path d="M16 3.13a4 4 0 0 1 0 7.75"></path></svg>`
	addIconSVG := `<svg viewBox="0 0 24 24"><path d="M12 19V5M5 12h14" stroke-linecap="round"/></svg>`
	weatherIconSVG := `<svg viewBox="0 0 24 24"><path d="M18 10h-1.26A8 8 0 1 0 9 20h9a5 5 0 0 0 0-10z"></path></svg>`
	profileIconSVG := `<svg viewBox="0 0 24 24"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path><circle cx="12" cy="7" r="4"></circle></svg>`

	return htma.Footer().AddChild(
		// Home Button (Active)
		htma.Button().ClassAttr("nav-button active").AddChild(
			htma.RawContent(homeIconSVG),
			htma.Span().Text("Home"),
		),
		// Shared Button
		htma.Button().ClassAttr("nav-button").AddChild(
			htma.RawContent(sharedIconSVG),
			htma.Span().Text("Shared"),
		),
		// Add Flight Button
		htma.Button().ClassAttr("nav-button").AddChild(
			htma.RawContent(addIconSVG),
			htma.Span().Text("Add Flight"),
		),
		// Weather Button
		htma.Button().ClassAttr("nav-button").AddChild(
			htma.RawContent(weatherIconSVG),
			htma.Span().Text("Weather"),
		),
		// Profile Button
		htma.Button().ClassAttr("nav-button").AddChild(
			htma.RawContent(profileIconSVG),
			htma.Span().Text("Profile"),
		),
	)
}
