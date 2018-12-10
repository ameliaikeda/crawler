package node

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// Node is a simple struct that can hold a list of URLs per-page, sitemap-style.
type Node struct {
	Children     []*Node
	URL          *url.URL
	ResponseCode int
	Visited      bool
	Depth        int
	Err          error
}

// String is a simple recursive function to pretty-print nodes.
func (n *Node) String() string {
	u := n.URL.String()

	if n.Err != nil || n.ResponseCode == http.StatusInternalServerError {
		return fmt.Sprintf("❌ %s %s", u, n.Err)
	}

	if !n.Visited {
		return fmt.Sprintf("⭕️ %s", u)
	}

	str := fmt.Sprintf("%s", u)

	if len(n.Children) > 0 {
		for _, c := range n.Children {
			padding := strings.Repeat("    ", n.Depth+1)

			str = str + "\n" + padding + "- " + c.String()
		}
	}

	return str
}
