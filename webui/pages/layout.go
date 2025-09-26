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
			htma.Script().TypeAttr("module").SrcAttr("/static/datastar.js"),
			htma.Script().TypeAttr("module").SrcAttr("/static/flightcard.js"),
			htma.Script().TypeAttr("module").SrcAttr("/static/searchcard.js"),
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
				htma.Main().
					AddChild(htma.SearchCard().
						Attr("data-on-searchtermchange__debounce300ms", "$searchTerm = evt.detail.value; @post('/search-flights')").
						Attr("data-fetch-url", "/search-flights").
						Attr("data-fetch-method", "post").
						Attr("data-fetch-body", "signals"),
						htma.Div().ClassAttr("container").DataSignalsAttr(`{ "searchTerm": "", "myFlights": [] }`),

						htma.Div().IDAttr("search-results").ClassAttr("search-results-container"),
						htma.Div().DataShowAttr("$myFlights.length > 0").AddChild(
							htma.H2().Text("Matching Flights"),
							htma.Ul().ClassAttr("tracked-flight-list").AddChild(
								htma.Template().Attr("data-for-flight", "$myFlights").AddChild(
									htma.Li().ClassAttr("tracked-flight-item").AddChild(
										htma.Span().DataTextAttr("`${$flight.id}: ${$flight.origin} → ${$flight.destination}`"),
									),
								),
							),
						),

						htma.Div().ClassAttr("flight-card-container").IDAttr("flights")).DataOnLoadAttr("@get('/sse/flights')"),
				components.FooterComponent(),
			),
	)
}
