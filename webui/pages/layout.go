/*
ACTION: Update this file.
PATH:   webui/pages/layout.go
PURPOSE: To initialize a Datastar signal that will manage the

	currently active page state on the client.
*/
package pages

import (
	"github.com/arcade55/htma"
	"github.com/arcade55/nzflights_webui/webui/components"
)

func LayoutComponent(title string, content htma.Renderable) htma.Element {
	return htma.HTML().LangAttr("en").AddChild(
		htma.Head().AddChild(
			// ... head content is unchanged ...
			htma.Meta().CharsetAttr("UTF-8"),
			htma.Meta().NameAttr("viewport").Attr("content", "width=device-width, initial-scale=1.0"),
			htma.Title(title),
			htma.Link().RelAttr("preconnect").HrefAttr("https://fonts.googleapis.com"),
			htma.Link().RelAttr("preconnect").HrefAttr("https://fonts.gstatic.com").CrossOriginAttr(""),
			htma.Link().HrefAttr("https://fonts.googleapis.com/css2?family=Roboto:wght@400;500;700&display=swap").RelAttr("stylesheet"),
			htma.Link().HrefAttr("https://fonts.googleapis.com/css2?family=Material+Symbols+Outlined:opsz,wght,FILL,GRAD@24,400,0,0").RelAttr("stylesheet"),
			htma.Link().RelAttr("stylesheet").HrefAttr("/static/style.css"),
		),
		// ADD data-signals and data-effect attributes HERE 👇
		htma.Body().
			DataSignalsAttr(`{ "page": "home" }`).
			DataEffectAttr(`
                if ($page === 'home') @get('/home-sse');
                if ($page === 'add-flight') @get('/add-flight-sse');
            `).
			AddChild(
				components.HeaderComponent(),
				htma.Main().IDAttr("main-content").AddChild(content),
				components.FooterComponent(),
				htma.Script().TypeAttr("module").SrcAttr("/static/datastar.js"),
			),
	)
}
