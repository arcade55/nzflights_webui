package pages

import (
	"github.com/arcade55/htma"
	"github.com/arcade55/nzflights_webui/webui/components"
)

// LayoutComponent creates the full HTML document structure.
// This is the corrected version with all necessary head tags.
func LayoutComponent(title string, content htma.Renderable) htma.Element {
	return htma.HTML().LangAttr("en").AddChild(
		// This Head() section now includes all the missing tags
		htma.Head().AddChild(
			htma.Meta().CharsetAttr("UTF-8"),
			// ADD THIS: The viewport meta tag for responsive design
			htma.Meta().NameAttr("viewport").Attr("content", "width=device-width, initial-scale=1.0"),
			htma.Title(title),

			// ADD THESE: The links for Google Fonts and Material Symbols
			htma.Link().RelAttr("preconnect").HrefAttr("https://fonts.googleapis.com"),
			htma.Link().RelAttr("preconnect").HrefAttr("https://fonts.gstatic.com").CrossOriginAttr(""),
			htma.Link().HrefAttr("https://fonts.googleapis.com/css2?family=Roboto:wght@400;500;700&display=swap").RelAttr("stylesheet"),
			htma.Link().HrefAttr("https://fonts.googleapis.com/css2?family=Material+Symbols+Outlined:opsz,wght,FILL,GRAD@24,400,0,0").RelAttr("stylesheet"),

			// This link to your local stylesheet is correct
			htma.Link().RelAttr("stylesheet").HrefAttr("/static/style.css"),
		),
		htma.Body().AddChild(
			components.HeaderComponent(),
			htma.Main().AddChild(content),
			components.FooterComponent(),
		),
	)
}
