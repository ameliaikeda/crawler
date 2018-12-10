package crawler

import (
	"context"
	"net/http"

	"github.com/ameliaikeda/crawler/node"
	"github.com/ameliaikeda/crawler/parser"
)

func (c *crawler) worker(ctx context.Context, done chan bool, id int) {
	// when done, exit.
	defer func() {
		done <- true
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case node, ok := <-c.queue.Get():
			if !ok {
				return
			}

			c.process(ctx, node)
		}
	}
}

func (c *crawler) process(ctx context.Context, n *node.Node) {
	// defer marking this request as complete.
	defer func() {
		err := c.queue.Del()
		if err != nil {
			panic(err)
		}
	}()

	c.normalize(n)

	// HTTP request only if Host matches that of base.URL.Host, otherwise return.
	// this means we can queue all nodes and make our code path a little neater.
	// we can also bail out if we've crawled too far down a path.
	if c.ignore(n) {
		return
	}

	u := n.URL.String()

	resp, err := c.config.HTTPClient.Get(u)
	if err != nil {
		n.Err = err

		return
	}

	defer resp.Body.Close()
	n.ResponseCode = resp.StatusCode
	n.Visited = true

	// bail out if we're not a 200
	if resp.StatusCode != http.StatusOK {
		return
	}

	links, err := parser.ParseLinks(resp.Body)
	if err != nil {
		n.Err = err

		return
	}

	c.pushChildren(n, links)
}
