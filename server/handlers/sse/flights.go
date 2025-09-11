package sse

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"

	"github.com/arcade55/htma"
	"github.com/arcade55/nzflights-models"
	"github.com/arcade55/nzflights_webui/server/handlers/middleware"
	"github.com/arcade55/nzflights_webui/webui/components"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/starfederation/datastar-go/datastar"
)

type FlightSSEHandler struct {
	KV jetstream.KeyValue
}

func (h *FlightSSEHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sse := datastar.NewSSE(w, r)
	ctx := r.Context()

	// Get the validated User ID from the context set by the AuthMiddleware.
	userID, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		// This should not happen if the middleware is applied, but it's a good safeguard.
		log.Println("Error: User ID not found in context.")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// --- END MODIFICATION ---

	ownedFlightsPattern := fmt.Sprintf("users.%s.flights.owned.>", userID)

	renderFlights := func() {
		var flights []nzflights.FlightValue
		// IMPORTANT: Use the request context for NATS operations
		keyLister, err := h.KV.ListKeysFiltered(ctx, ownedFlightsPattern)
		if err != nil {
			// This will now correctly log an error if the client disconnects mid-operation
			log.Printf("Error listing keys for %s: %v", ownedFlightsPattern, err)
			return
		}

		for key := range keyLister.Keys() {
			kvEntry, err := h.KV.Get(ctx, key)
			if err != nil {
				log.Printf("Error fetching flight %s: %v", key, err)
				continue
			}
			var fv nzflights.FlightValue
			if err := json.Unmarshal(kvEntry.Value(), &fv); err != nil {
				log.Printf("Error unmarshaling flight %s: %v", key, err)
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

		content := htma.Div().IDAttr("flight-card-container").AddChild(flightCards...)
		if err := sse.PatchElements(content.Render(),
			datastar.WithSelector("#flight-card-container"),
			datastar.WithMode("replace"),
		); err != nil {
			log.Printf("Error sending patch: %v", err)
		}
	}

	// **THE FIX**
	// 1. Explicitly render the initial state first.
	//renderFlights()

	// 2. Then, create the watcher to listen for *future* changes.
	watcher, err := h.KV.Watch(ctx, ownedFlightsPattern)
	if err != nil {
		log.Printf("Error creating watcher for user %s: %v", userID, err)
		return
	}
	defer watcher.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Printf("Client for user %s disconnected.", userID)
			return
		case entry := <-watcher.Updates():
			if entry == nil {
				continue
			}
			renderFlights()
		}
	}
}
