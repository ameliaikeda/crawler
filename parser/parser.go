package parser

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

// ParseLinks takes a HTML body and returns a slice of all URLs found within <a> tags.
// It will remove duplicates and links that are simply fragments on the same page.
func ParseLinks(body io.Reader) ([]string, error) {
	links := make([]string, 0)
	duplicates := make(map[string]bool)
	s := html.NewTokenizer(body)

	for {
		token := s.Next()

		switch {
		case token == html.ErrorToken:
			err := s.Err()
			if err == io.EOF {
				err = nil
			}

			return links, err
		case token == html.StartTagToken:
			t := s.Token()

			if t.Data != "a" {
				continue
			}

			for _, tag := range t.Attr {
				if tag.Key == "href" && strings.Index(tag.Val, "#") != 0 {
					// check for duplicates; if there's any, we don't include them.
					if _, ok := duplicates[tag.Val]; !ok {
						duplicates[tag.Val] = true
						links = append(links, tag.Val)
					}
				}
			}
		}
	}
}
