/*
ACTION: Update this file.
PATH:   main.go
PURPOSE: To fix both the client-side debounce syntax and the server-side

	"body already closed" error, making the active search functional.
*/
package search

/*
import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/arcade55/htma"
	"github.com/starfederation/datastar-go/datastar"
)

// The Flight struct and helper functions remain the same...
type Flight struct {
	ID          string `json:"id"`
	Origin      string `json:"origin"`
	Destination string `json:"destination"`
}

func generateFlights() []Flight {
	rand.Seed(time.Now().UnixNano())
	flights := make([]Flight, 100)
	airports := []string{"AKL", "WLG", "CHC", "ZQN", "LAX", "SYD", "MEL"}
	for i := 0; i < 100; i++ {
		origin := airports[rand.Intn(len(airports))]
		dest := airports[rand.Intn(len(airports))]
		if origin == dest {
			dest = airports[(rand.Intn(len(airports)-1)+1)%len(airports)]
		}
		flights[i] = Flight{
			ID:          fmt.Sprintf("NZ%d", 500+rand.Intn(499)),
			Origin:      origin,
			Destination: dest,
		}
	}
	return flights
}

func miniCardComponent(flight Flight) htma.Element {
	onClickScript := fmt.Sprintf(`
            if (!$myFlights.find(f => f.id === '%s')) {
                $myFlights.push({ id: '%s', origin: '%s', destination: '%s' })
            };
            $searchTerm = '';
        `, flight.ID, flight.ID, flight.Origin, flight.Destination)

	return htma.Div().ClassAttr("mini-card").
		DataOnClickAttr(onClickScript).
		AddChild(
			htma.Div().ClassAttr("flight-info").AddChild(
				htma.Div().ClassAttr("flight-id").Text(flight.ID),
				htma.Div().ClassAttr("flight-route").Text(fmt.Sprintf("%s → %s", flight.Origin, flight.Destination)),
			),
			htma.Span().Text("＋ Add"),
		)
}

func main() {
	// The main page handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		doc := htma.HTML().LangAttr("en").AddChild(
			htma.Head().AddChild(
				htma.Meta().CharsetAttr("UTF-8"),
				htma.Title("NZ Flight Search"),
				htma.Script().TypeAttr("module").SrcAttr("https://cdn.jsdelivr.net/gh/starfederation/datastar@main/bundles/datastar.js"),
				htma.RawContent(`<style>
                    body { font-family: system-ui, sans-serif; max-width: 600px; margin: 2rem auto; padding: 0 1rem; background: #f8f9fa; color: #212529; }
                    h1, h2 { font-weight: 700; color: #000; }
                    h1 { font-size: 2rem; margin-bottom: 1.5rem; }
                    h2 { font-size: 1.5rem; margin-top: 2.5rem; margin-bottom: 1rem; }
                    input { width: 100%; padding: 0.75rem; font-size: 1rem; margin-bottom: 0.5rem; border: 1px solid #ced4da; border-radius: 8px; box-shadow: 0 1px 3px rgba(0,0,0,0.05); }
                    input:focus { outline: none; border-color: #4A90E2; box-shadow: 0 0 0 3px rgba(74, 144, 226, 0.25); }
                    .search-results-container { min-height: 1px; background: #fff; border-radius: 8px; box-shadow: 0 4px 12px rgba(0,0,0,0.1); overflow: hidden; }
                    .mini-card { display: flex; justify-content: space-between; align-items: center; padding: 0.75rem 1rem; cursor: pointer; border-bottom: 1px solid #e9ecef; }
                    .mini-card:last-child { border-bottom: none; }
                    .mini-card:hover { background-color: #f1f3f5; }
                    .mini-card .flight-id { font-weight: 600; }
                    .mini-card .flight-route { font-size: 0.9em; color: #6c757d; }
                    .tracked-flight-list { list-style: none; padding: 0; margin: 0; }
                    .tracked-flight-item { display: flex; justify-content: space-between; align-items: center; padding: 1rem; margin-bottom: 0.5rem; background: #fff; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.05); }
                </style>`),
			),
			htma.Body().AddChild(
				htma.Div().
					DataSignalsAttr(`{ "searchTerm": "", "myFlights": [] }`).
					AddChild(
						htma.H1().Text("✈️ Add a Flight"),
						htma.Input().
							TypeAttr("text").
							PlaceholderAttr("Search flight number (e.g. NZ622)...").
							Attr("data-bind", "searchTerm").
							// FIX: Use the correct double-underscore syntax for the modifier
							Attr("data-on-input__debounce-300ms", "@post('/search-flights')"),

						htma.Div().IDAttr("search-results").ClassAttr("search-results-container"),

						htma.Div().DataShowAttr("$myFlights.length > 0").AddChild(
							htma.H2().Text("Your Tracked Flights"),
							htma.Ul().ClassAttr("tracked-flight-list").AddChild(
								htma.Template().Attr("data-for-flight", "$myFlights").AddChild(
									htma.Li().ClassAttr("tracked-flight-item").AddChild(
										htma.Span().DataTextAttr("`${$flight.id}: ${$flight.origin} → ${$flight.destination}`"),
									),
								),
							),
						),
					),
			),
		)

		w.Header().Set("Content-Type", "text/html")
		doc.RenderStream(w)
	})

	// SSE handler that performs the search on the server.
	http.HandleFunc("/search-flights", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		log.Println("-> Received request on /search-flights")

		// FIX: Read signals from the body BEFORE creating the SSE stream
		signals := struct {
			SearchTerm string `json:"searchTerm"`
		}{}
		if err := datastar.ReadSignals(r, &signals); err != nil {
			log.Printf("ERROR: Could not read signals: %v", err)
			http.Error(w, "Could not read signals", http.StatusBadRequest)
			return
		}
		log.Printf("   - Received SearchTerm: '%s'", signals.SearchTerm)

		// NOW it's safe to create the SSE stream
		sse := datastar.NewSSE(w, r)

		// The rest of the logic remains the same...
		if signals.SearchTerm == "" {
			log.Println("   - Search term is empty, clearing results.")
			sse.PatchElements("", datastar.WithSelector("#search-results"), datastar.WithModeInner())
			return
		}

		allFlights := generateFlights()
		matches := []Flight{}
		for _, flight := range allFlights {
			if strings.Contains(strings.ToLower(flight.ID), strings.ToLower(signals.SearchTerm)) {
				matches = append(matches, flight)
				if len(matches) >= 5 {
					break
				}
			}
		}
		log.Printf("   - Found %d matches.", len(matches))

		var sb strings.Builder
		for _, flight := range matches {
			miniCardComponent(flight).RenderStream(&sb)
		}
		htmlFragment := sb.String()
		log.Printf("   - Generated HTML fragment of length: %d", len(htmlFragment))

		log.Println("   - Sending fragment to client...")
		sse.PatchElements(htmlFragment, datastar.WithSelector("#search-results"), datastar.WithModeInner())
		log.Println("<- Request finished.")
	})

	log.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
*/
