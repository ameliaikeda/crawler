package crawler

import (
	"context"
	"testing"
	"time"
)

func TestNewCrawlerWithInvalidURL(t *testing.T) {
	site := "http://192.168.0.%31/"
	cfg := NewConfig(context.Background(), 100, 5, time.Second*1)

	_, err := New(site, cfg)
	if err == nil || err == errNoHost {
		t.Error("invalid URLs should generate an error")
	}
}

func TestNewCrawlerWithRelativeURL(t *testing.T) {
	site := "/blog/"
	cfg := NewConfig(context.Background(), 100, 5, time.Second*1)

	_, err := New(site, cfg)
	if err != errNoHost {
		t.Error("relative URLs should return errNoHost")
	}
}
