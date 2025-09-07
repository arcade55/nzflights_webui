/*
ACTION: Update this file.
PATH:   webui/components/footer.go
PURPOSE: To fix the "flicker" bug by preventing default browser

	navigation and to use a declarative, signal-based approach
	for managing the active UI state.
*/
package components

import "github.com/arcade55/htma"

func FooterComponent() htma.Element {
	// --- SVG definitions remain the same ---
	homeIconSVG := `<svg viewBox="0 0 24 24"><path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"></path><polyline points="9 22 9 12 15 12 15 22"></polyline></svg>`
	sharedIconSVG := `<svg viewBox="0 0 24 24"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path><circle cx="9" cy="7" r="4"></circle><path d="M23 21v-2a4 4 0 0 0-3-3.87"></path><path d="M16 3.13a4 4 0 0 1 0 7.75"></path></svg>`
	addIconSVG := `<svg viewBox="0 0 24 24"><path d="M12 19V5M5 12h14" stroke-linecap="round"/></svg>`
	weatherIconSVG := `<svg viewBox="0 0 24 24"><path d="M18 10h-1.26A8 8 0 1 0 9 20h9a5 5 0 0 0 0-10z"></path></svg>`
	profileIconSVG := `<svg viewBox="0 0 24 24"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path><circle cx="12" cy="7" r="4"></circle></svg>`

	return htma.Footer().AddChild(
		// --- Home Link ---
		htma.A().ClassAttr("nav-button").
			// The button is active ONLY when the $page signal is 'home'
			DataClassAttr("{ 'active': $page === 'home' }").
			// CHANGE THIS: Href is now non-navigational
			HrefAttr("#!").
			// The button's only job is to change the state
			DataOnClickAttr("$page = 'home'").
			Attr("data-target-selector", "#main-content").
			AddChild(
				htma.RawContent(homeIconSVG),
				htma.Span().Text("Home"),
			),

		// --- Shared Link (unchanged) ---
		htma.A().ClassAttr("nav-button").HrefAttr("/shared").AddChild(
			htma.RawContent(sharedIconSVG),
			htma.Span().Text("Shared"),
		),

		// --- Add Flight Link ---
		htma.A().ClassAttr("nav-button").
			// The button is active ONLY when the $page signal is 'add-flight'
			DataClassAttr("{ 'active': $page === 'add-flight' }").
			// CHANGE THIS: Href is now non-navigational
			HrefAttr("#!").
			// The button's only job is to change the state
			DataOnClickAttr("$page = 'add-flight'").
			Attr("data-target-selector", "#main-content").
			AddChild(
				htma.RawContent(addIconSVG),
				htma.Span().Text("Add Flight"),
			),

		// --- Weather & Profile Links (unchanged) ---
		htma.A().ClassAttr("nav-button").HrefAttr("/weather").AddChild(
			htma.RawContent(weatherIconSVG),
			htma.Span().Text("Weather"),
		),
		htma.A().ClassAttr("nav-button").HrefAttr("/profile").AddChild(
			htma.RawContent(profileIconSVG),
			htma.Span().Text("Profile"),
		),
	)
}
