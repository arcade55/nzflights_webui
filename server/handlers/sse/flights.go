package sse

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	"github.com/arcade55/htma"
	"github.com/arcade55/logging"
	"github.com/arcade55/nzflights-models"
	"github.com/arcade55/nzflights_webui/server/handlers/middleware"
	"github.com/arcade55/nzflights_webui/webui/components"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/starfederation/datastar-go/datastar"
)

type FlightSSEHandler struct {
	KV jetstream.KeyValue
}

// Initialize the logger
var logger, _, _ = logging.Init(context.Background(), logging.Config{
	Format:    logging.FormatPretty,
	Level:     logging.LevelDebug,
	AddSource: true,
})
var log = logger.WithContext(context.Background())

func (h *FlightSSEHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	log.Info("Starting interactive UI test server...")

	visitorCookie, err := r.Cookie(middleware.VisitorCookieName)
	if err != nil {
		// This should technically never happen because the middleware adds it.
		http.Error(w, "User could not be identified ", http.StatusInternalServerError)
		return
	}

	visitorID := visitorCookie.Value

	sse := datastar.NewSSE(w, r)
	ctx := r.Context()

	if visitorID == "" {
		// This should not happen if the middleware is applied, but it's a good safeguard.
		log.Error(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// --- END MODIFICATION ---

	ownedFlightsPattern := fmt.Sprintf("users.%s.flights.owned.>", visitorID)

	renderFlights := func() {
		var flights []nzflights.FlightValue
		// IMPORTANT: Use the request context for NATS operations
		keyLister, err := h.KV.ListKeysFiltered(ctx, ownedFlightsPattern)
		if err != nil {
			// This will now correctly log an error if the client disconnects mid-operation
			log.Error(err)
			return
		}

		for key := range keyLister.Keys() {
			kvEntry, err := h.KV.Get(ctx, key)
			if err != nil {
				log.Error(err)
				continue
			}
			var fv nzflights.FlightValue
			if err := json.Unmarshal(kvEntry.Value(), &fv); err != nil {
				log.Error(err)
				continue
			}
			flights = append(flights, fv)
		}

		sort.Slice(flights, func(i, j int) bool {
			return flights[i].Flight.Ident < flights[j].Flight.Ident
		})

		var flightCards []htma.Renderable
		for _, fv := range flights {
			flightCards = append(flightCards, components.FlightCardComponent(fv))
		}
		content := htma.Div().ClassAttr("flight-card-container").IDAttr("flights").AddChild(flightCards...)
		log.Info(content.Render())
		if err := sse.PatchElements(content.Render(),
			datastar.WithSelector("#flights"),
			datastar.WithMode("replace"),
		); err != nil {
			log.Error(err)
		}
	}

	watcher, err := h.KV.Watch(ctx, ownedFlightsPattern)
	if err != nil {
		log.Error(err)
		return
	}
	defer watcher.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info(fmt.Sprintf("Client for user %s disconnected.", visitorID))
			return
		case entry := <-watcher.Updates():
			if entry == nil {
				continue
			}
			log.Info(string(entry.Value()))
			renderFlights()
		}
	}
}
