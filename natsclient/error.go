package natsclient

import "errors"

var (
	// --- General Setup Errors ---
	ErrEmbeddedServerFailed   = errors.New("embedded NATS server setup failed")
	ErrCloudConnectionFailed  = errors.New("cloud NATS server connection failed")
	ErrJetStreamContextFailed = errors.New("failed to create JetStream context")

	// --- Key-Value Store Errors ---
	ErrKVStoreMirrorFailed = errors.New("failed to create mirrored Key-Value store")
	ErrKVStoreBindFailed   = errors.New("failed to bind to cloud Key-Value store")

	// --- Watcher Errors ---
	ErrWatcherCreationFailed = errors.New("failed to create Key-Value watcher")
)
