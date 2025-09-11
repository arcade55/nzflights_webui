package natsclient

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/arcade55/logging"
	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// Client holds all the high-level NATS resources the application will interact with.
type Client struct {
	// Flights provides a dedicated interface for all flight data operations.
	Flights FlightStore
	// InMemoryKV provides direct access to the in-memory mirror for monitoring.
	InMemoryKV jetstream.KeyValue
	// Publish a message to trigger an API fetch for a flight.
	TriggerAPIFetch func(flightID string) error

	// A function to gracefully clean up connections and the server.
	Shutdown func()
}

const SynadiaCloudURL = "tls://connect.ngs.global:4222"

// New creates and configures the entire NATS stack for the application.
func New(ctx context.Context, logger *logging.Logger, cloudCreds []byte) (*Client, error) {
	log := logger.WithContext(ctx)

	// --- 1. Set up and run the Embedded Leaf Server ---
	embeddedNC, embeddedServer, err := runEmbeddedServer(true, true)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrEmbeddedServerFailed, err)
	}
	log.Info("✅ Embedded NATS leaf server started.")

	embeddedJS, err := jetstream.New(embeddedNC)
	if err != nil {
		embeddedServer.Shutdown()
		return nil, fmt.Errorf("%w on embedded server: %v", ErrJetStreamContextFailed, err)
	}

	// --- 2. Create the In-Memory Mirrored Key-Value Store ---
	mirrorConfig := jetstream.KeyValueConfig{
		Bucket:  "inMemoryFlights",
		Storage: jetstream.MemoryStorage,
		Mirror: &jetstream.StreamSource{
			Name:   "flights",
			Domain: "ngs", // Updated domain
		},
	}
	inMemoryKV, err := embeddedJS.CreateOrUpdateKeyValue(ctx, mirrorConfig)
	if err != nil {
		embeddedServer.Shutdown()
		return nil, fmt.Errorf("%w: %v", ErrKVStoreMirrorFailed, err)
	}
	log.Info("✅ In-memory KV store configured to mirror 'flights'.")

	// --- 3. Connect to the Cloud NATS Server ---
	cloudNC, err := nats.Connect(SynadiaCloudURL, nats.Name("nzflights_webui"), nats.UserCredentialBytes(cloudCreds))
	if err != nil {
		embeddedServer.Shutdown()
		return nil, fmt.Errorf("%w: %v", ErrCloudConnectionFailed, err)
	}
	log.Info("✅ Connected to cloud NATS server.")

	cloudJS, err := jetstream.New(cloudNC)
	if err != nil {
		cloudNC.Close()
		embeddedServer.Shutdown()
		return nil, fmt.Errorf("%w on cloud server: %v", ErrJetStreamContextFailed, err)
	}

	// --- 4. Get a handle to the actual Cloud Key-Value Store ---
	cloudKV, err := cloudJS.KeyValue(ctx, "flights")
	if err != nil {
		cloudNC.Close()
		embeddedServer.Shutdown()
		return nil, fmt.Errorf("%w: %v", ErrKVStoreBindFailed, err)
	}
	log.Info("✅ Bound to cloud 'flights' KV store.")

	// --- 5. Construct the final Client object ---
	client := &Client{
		// HERE is where the 'unused' code is now being used.
		Flights:    newFlightStore(inMemoryKV, cloudKV),
		InMemoryKV: inMemoryKV,

		TriggerAPIFetch: func(flightID string) error {
			subject := fmt.Sprintf("api.flightaware.fetch.%s", flightID)
			return cloudNC.Publish(subject, nil)
		},
		Shutdown: func() {
			log.Info("Shutting down NATS client and server...")
			cloudNC.Close()
			embeddedServer.Shutdown()
			log.Info("Shutdown complete.")
		},
	}

	return client, nil
}

// runEmbeddedServer starts an embedded NATS server configured as a leaf node.
func runEmbeddedServer(inProcess bool, enableLogging bool) (*nats.Conn, *server.Server, error) {
	leafURL, err := url.Parse("nats-leaf://connect.ngs.global")
	if err != nil {
		return nil, nil, err
	}
	opts := &server.Options{
		ServerName: "FlightApp_LeafNode",
		StoreDir:   "",
		DontListen: inProcess,
		JetStream:  true,
		LeafNode: server.LeafNodeOpts{
			Remotes: []*server.RemoteLeafOpts{{
				URLs:        []*url.URL{leafURL},
				Credentials: "./natsclient/leafnode.cred",
			}},
		},
	}
	ns, err := server.NewServer(opts)
	if err != nil {
		return nil, nil, err
	}
	if enableLogging {
		ns.ConfigureLogger()
	}
	go ns.Start()
	if !ns.ReadyForConnections(5 * time.Second) {
		return nil, nil, fmt.Errorf("embedded server not ready for connections")
	}
	clientOpts := []nats.Option{}
	if inProcess {
		clientOpts = append(clientOpts, nats.InProcessServer(ns))
	}
	nc, err := nats.Connect(nats.DefaultURL, clientOpts...)
	if err != nil {
		return nil, nil, err
	}
	return nc, ns, nil
}
