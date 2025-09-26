package sse

import (
	"fmt"
	"net/http"
	"strings"

	"log/slog"

	"github.com/arcade55/htma"
	"github.com/arcade55/nzflights-models"
	"github.com/starfederation/datastar-go/datastar"
)

type SearchSSEHandler struct {
	Flights map[string]nzflights.FlightValue
}

func (h *SearchSSEHandler) Search(w http.ResponseWriter, r *http.Request) {
	// DEBUG: Log when the endpoint is hit.
	log.Info("-> Received request for /search-flights")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	signals := struct {
		SearchTerm string `json:"searchTerm"`
	}{}
	if err := datastar.ReadSignals(r, &signals); err != nil {
		log.Info("   - ERROR: Could not read signals", slog.Any("error", err))
		http.Error(w, "Could not read signals", http.StatusBadRequest)
		return
	}

	sse := datastar.NewSSE(w, r)

	if signals.SearchTerm == "" {
		sse.PatchElements("", datastar.WithSelector("#search-results"), datastar.WithModeInner())
		return
	}

	var matches []nzflights.FlightValue
	for _, flight := range h.Flights {

		log.Info(flight.Flight.Ident)
		if strings.Contains(strings.ToLower(flight.Flight.Ident), strings.ToLower(signals.SearchTerm)) {
			matches = append(matches, flight)
			if len(matches) >= 5 {
				break
			}
		}
	}
	log.Info("   - Found matches.", slog.Int("count", len(matches)))

	var sb strings.Builder
	for _, flight := range matches {
		miniCardComponent(flight).RenderStream(&sb)
	}
	htmlFragment := sb.String()

	log.Info("sending this to server")
	log.Info(htmlFragment)

	err := sse.PatchElements(htmlFragment, datastar.WithSelector("#search-results"), datastar.WithModeInner())
	if err != nil {
		log.Error(err)
	}
}

func miniCardComponent(flight nzflights.FlightValue) htma.Element {
	onClickScript := fmt.Sprintf(`
        if (!$myFlights.find(f => f.id === '%s')) {
            $myFlights.push({ id: '%s', origin: '%s', destination: '%s' })
        };
        $searchTerm = '';
    `, flight.Flight.Ident, flight.Flight.Ident, flight.Flight.OriginCity, flight.Flight.DestinationCity)

	return htma.Div().ClassAttr("mini-card").
		DataOnClickAttr(onClickScript).
		AddChild(
			htma.Div().ClassAttr("flight-info").AddChild(
				htma.Div().ClassAttr("flight-id").Text(flight.Flight.Ident),
				htma.Div().ClassAttr("flight-route").Text(fmt.Sprintf("%s → %s", flight.Flight.OriginCity, flight.Flight.DestinationCity)),
			),
			htma.Span().Text("＋ Add"),
		)
}
