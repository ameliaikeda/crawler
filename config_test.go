package crawler

import (
	"context"
	"testing"
	"time"
)

func TestNewCrawlerConcurrency(t *testing.T) {
	cfg := NewConfig(context.Background(), 1234, 5, time.Second*1)

	if cfg.WorkerCount != 1234 {
		t.Error("when setting workers to 1234, WorkerCount should be 1234")
	}
}

func TestNewCrawlerDepth(t *testing.T) {
	cfg := NewConfig(context.Background(), 100, 5959, time.Second*1)

	if cfg.Depth != 5959 {
		t.Error("when setting depth to 5959, Depth should be 5959")
	}
}

func TestNewCrawlerContext(t *testing.T) {
	ctx := context.Background()
	cfg := NewConfig(ctx, 100, 5959, time.Second*5)

	if ctx == cfg.Context {
		t.Error("NewConfig should derive a context from whatever it is given")
	}

	cfg.CancelFunc()
	err := cfg.Context.Err()

	if err != context.Canceled {
		t.Error("cancelling config's context should return Canceled")
	}
}
