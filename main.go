package main

import (
	"context"
	"embed"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	// 1. ADD THIS IMPORT BACK IN

	"github.com/starfederation/datastar-go/datastar"

	// These imports are needed for the SSE handlers
	"github.com/arcade55/htma"
	"github.com/arcade55/logging"
	correlation "github.com/arcade55/nzflights-correlation"
	"github.com/arcade55/nzflights_webui/natsclient"
	"github.com/arcade55/nzflights_webui/server/handlers/middleware"
	"github.com/arcade55/nzflights_webui/server/handlers/standard"
	"github.com/arcade55/nzflights_webui/webui/components"
	"github.com/arcade55/nzflights_webui/webui/pages"
)

//go:embed nats.cred
var credsFile embed.FS



func main() {
	ctx := correlation.EnsureCorrelationID(context.Background())

	cfg := logging.Config{
		ServiceName: "my-app-service", // Required for base attributes.
		Output:      os.Stdout,        // Or a file, etc.
		Level:       logging.LevelInfo,
		Format:      logging.FormatPretty, // Or JSON/Text.
		AddSource:   true,                 // Includes file/line in logs.
	}

	logger, cleanup, err := logging.Init(ctx, cfg)

	creds, err := credsFile.ReadFile("nats.cred")
	if err != nil {
		log.Error(err, slog.String("action", "embedded_file_error"), slog.String("message", "failed to read embedded credentials file"))
		os.Exit(1)
	}

	// --- 2. Initialize the NATS Client ---
	// This single call sets up everything: embedded server, cloud connection, mirrors, etc.
	client, err := natsclient.New(ctx, logger, creds)
	if err != nil {
		if errors.Is(err, natsclient.ErrCloudConnectionFailed) {
			log.Error(err)
		} else {
			log.Error(err)
		}
		os.Exit(1)
	}
	defer client.Shutdown()
	log.Info("🚀 Application started successfully. NATS client is ready.")

	// --- Setup Graceful Shutdown ---
	// Create a channel to listen for OS signals.
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGINT, syscall.SIGTERM)

	mux := http.NewServeMux()

	staticFS := http.FileServer(http.Dir("./webui/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", staticFS))

	home := http.HandlerFunc(standard.HomeHandler)

	mux.Handle("GET /home", middleware.VisitorID(home))
	mux.HandleFunc("GET /home-sse", handleHomeSSE)
	mux.HandleFunc("GET /add-flight", handleAddFlight)
	mux.HandleFunc("GET /add-flight-sse", handleAddFlightSSE)
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/home", http.StatusMovedPermanently)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Info("Starting server on port " + port)
	err = http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Error(err)
	}
}

// handleHomeSSE handles Datastar requests for the home page content.
func handleHomeSSE(w http.ResponseWriter, r *http.Request) {
	sse := datastar.NewSSE(w, r)
	// The ID for this container must match the selector used in the SSE handler
	mainContent := htma.Div().ClassAttr("flight-card-container").IDAttr("flights")
	sse.PatchElements(mainContent.Render(), datastar.WithSelector("#main-content"), datastar.WithModeInner())
}

// handleAddFlight serves the initial, full HTML page.
func handleAddFlight(w http.ResponseWriter, r *http.Request) {
	page := pages.AddFlightPage()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	page.RenderStream(w)
}

// handleAddFlightSSE handles Datastar requests for the add flight page content.
func handleAddFlightSSE(w http.ResponseWriter, r *http.Request) {
	sse := datastar.NewSSE(w, r)
	content := htma.Div().ClassAttr("search-container").AddChild(
		htma.Div().ClassAttr("ticket-card search-ticket").AddChild(
			components.InputField("airplane_ticket", "", "", "Optional Flight Number"),
			components.Separator("Or"),
			components.InputField("flight_takeoff", "From", "Auckland, NZ", ""),
			components.InputField("flight_land", "", "", "Select destination"),
			components.InputField("calendar_month", "Departure", "Sep 01, 2025", ""),
			components.ActionButton("SSE Search Flights"),
		),
	)

	sse.PatchElements(content.Render(), datastar.WithSelector("#main-content"), datastar.WithModeInner())
}
