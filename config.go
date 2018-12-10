package crawler

import (
	"context"
	"net/http"
	"time"
)

var defaultHTTPClient = &http.Client{
	Timeout: time.Second * 5,
}

// Config holds basic configuration for a Crawler
type Config struct {
	// Depth is how many links we will follow consecutively.
	Depth int

	// WorkerCount is how many workers we run to process URLs.
	WorkerCount int

	// HTTPClient allows customising the HTTP client used to request.
	HTTPClient *http.Client

	// Context is used to send cancellation signals to worker nodes.
	Context context.Context

	// CancelFunc is called when signalling that we are done processing.
	CancelFunc context.CancelFunc
}

// NewConfig returns a new Config instance.
func NewConfig(ctx context.Context, workers, depth int, timeout time.Duration) *Config {
	c, cancel := context.WithTimeout(ctx, timeout)

	return &Config{
		Depth:       depth,
		WorkerCount: workers,
		HTTPClient:  defaultHTTPClient,
		Context:     c,
		CancelFunc:  cancel,
	}
}
