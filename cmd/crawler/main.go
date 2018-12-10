package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/ameliaikeda/crawler"
)

var (
	site           string
	workers, depth int
	timeout        time.Duration
)

func main() {
	flag.StringVar(&site, "site", "", "site to crawl. should be a valid URL.")
	flag.IntVar(&workers, "workers", 8, "number of worker processes to run at once.")
	flag.IntVar(&depth, "depth", 5, "how many levels of links to follow when crawling.")
	flag.DurationVar(&timeout, "timeout", time.Minute*5, "how long to wait before halting execution")
	flag.Parse()

	cfg := crawler.NewConfig(context.Background(), workers, depth, timeout)

	c, err := crawler.New(site, cfg)
	if err != nil {
		panic(err)
	}

	sitemap, err := c.Process()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s", sitemap)
}
