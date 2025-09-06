package pages

import (
	"github.com/arcade55/htma"
	"github.com/arcade55/nzflights_webui/webui/components"
)

// LayoutComponent creates the full HTML document structure.
// It is the base layout that all pages will use.
func LayoutComponent(title string, content htma.Renderable) htma.Element {
	return htma.HTML().LangAttr("en").AddChild(
		htma.Head().AddChild(
			// ... all your meta and link tags ...
			htma.Meta().CharsetAttr("UTF-8"),
			htma.Title(title),
			htma.Link().RelAttr("stylesheet").HrefAttr("/static/style.css"),
		),
		htma.Body().AddChild(
			components.HeaderComponent(),
			htma.Main().AddChild(content), // Page-specific content is injected here
			components.FooterComponent(),
		),
	)
}
