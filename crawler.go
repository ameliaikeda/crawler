package crawler

import (
	"errors"
	"fmt"
	"net/url"
	"sync"

	"github.com/ameliaikeda/crawler/node"
	"github.com/ameliaikeda/crawler/queue"
	"github.com/hashicorp/errwrap"
)

const urlParsingFailure = "crawler: failed to parse url '%s', {{err}}"

var errNoHost = errors.New("no domain specified when crawling")

// Crawler is a generic interface for something that can crawl a sitemap.
type Crawler interface {
	// Process will run the crawler with its given concurrency settings,
	// running until all pages have been crawled or a timeout is hit.
	Process() (*node.Node, error)
}

type crawler struct {
	sync.RWMutex
	pages map[string]bool

	queue  queue.Queue
	config *Config
	base   *node.Node
}

// Process will run the crawler with its given concurrency settings, running until all pages have been crawled.
func (c *crawler) Process() (*node.Node, error) {
	// push the base node onto the queue to kick everything off.
	c.queue.Add(c.base)

	done := make(chan bool)

	for i := 0; i < c.config.WorkerCount; i++ {
		go c.worker(c.config.Context, done, i)
	}

	// wait on workers to complete
	for w := 0; w < c.config.WorkerCount; {
		<-done
		w++
	}

	return c.base, nil
}

// New creates a new Crawler, set up to process the given URL.
func New(site string, config *Config) (Crawler, error) {
	u, err := url.Parse(site)
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf(urlParsingFailure, u), err)
	}

	if u.Host == "" {
		return nil, errNoHost
	}

	return &crawler{
		pages:  make(map[string]bool),
		queue:  queue.New(config.CancelFunc),
		base:   &node.Node{URL: u, Children: make([]*node.Node, 0)},
		config: config,
	}, nil
}

// normalize takes a node and sets up domain+scheme if they are not set (relative URLs)
func (c *crawler) normalize(n *node.Node) {
	if !n.URL.IsAbs() {
		n.URL.Host = c.base.URL.Host
		n.URL.Scheme = c.base.URL.Scheme
	}
}

func (c *crawler) pushChildren(n *node.Node, links []string) {
	n.Children = createNodes(links, n.Depth+1)

	// now loop through children and queue them all.
	for _, child := range n.Children {
		c.normalize(child)
		u := child.URL.String()

		// avoid re-queuing visited nodes; our buffer is only so big.
		if !c.visited(u) {
			// visit the node to lock other goroutines from queuing it.
			c.visit(u)
			c.queue.Add(child)
		}
	}
}

// ignore returns true if we should skip this node.
func (c *crawler) ignore(n *node.Node) bool {
	return n.URL.Host != c.base.URL.Host || n.Depth > c.config.Depth
}

// visited will check if we've already got a request pending for visit.
func (c *crawler) visited(u string) bool {
	c.RLock()
	defer c.RUnlock()

	_, ok := c.pages[u]
	return ok
}

// visit marks a URL as visited.
func (c *crawler) visit(u string) {
	c.Lock()
	defer c.Unlock()

	c.pages[u] = true
}

func createNodes(urls []string, depth int) []*node.Node {
	nodes := make([]*node.Node, 0)

	for _, raw := range urls {
		u, err := url.Parse(raw)
		if err != nil {
			continue
		}

		nodes = append(nodes, &node.Node{URL: u, Depth: depth})
	}

	return nodes
}
