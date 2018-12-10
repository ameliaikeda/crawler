// Package crawler will run a pseudo "search crawler" on a
// given domain, when given a URL, following any internal links on that domain.
// The crawler will only crawl internal pages on that domain,
// and will discard any other links it sees.
//
//  package main
//
//  import (
//      "github.com/ameliaikeda/crawler"
//  )
//
//  func main() {
//      site := "https://site.example/"
//      cfg := crawler.NewConfig(context.Background(), 100, 10, 5*time.Second)
//      c := crawler.New(site, cfg)
//      sitemap, err := c.Process()
//      fmt.Printf("%s", sitemap)
//  }
package crawler
