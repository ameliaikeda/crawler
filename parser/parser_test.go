package parser

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestBasicLinks(t *testing.T) {
	r := strings.NewReader(`
<a href="https://foo.example/bar"></a>
<a href="http://foo.example/bar"></a>
<a href="/bar"></a>
<a href="tel:55555555"></a>
<img src="/img/foo.png">
<link rel="stylesheet" href="http://foo.example/foo">
`)

	links, err := ParseLinks(r)
	if err != nil {
		t.Error(err)
	}

	if len(links) != 4 {
		t.Errorf("FindLinks should find 4 links, %d found", len(links))
	}

}

func TestInvalidHTML(t *testing.T) {
	r := strings.NewReader(`<html`)

	links, err := ParseLinks(r)
	if err == io.EOF {
		t.Error("FindLinks should not return io.EOF on invalid HTML.")
	}

	// links should be empty on a failure
	if len(links) != 0 {
		t.Error("On an error, FindLinks should return a zero-length slice")
	}
}

func TestValidFixture(t *testing.T) {
	f, err := os.Open("fixtures/pass.html")
	if err != nil {
		t.Error(err)
	}

	defer f.Close()

	links, err := ParseLinks(f)
	if err != nil {
		t.Error(err)
	}

	if len(links) != 46 {
		t.Errorf("46 links expected from fixture, got %d", len(links))
	}
}

func TestFragmentLinks(t *testing.T) {
	r := strings.NewReader(`
<a href="https://foo.example/bar"></a>
<a href="http://foo.example/bar"></a>
<a href="/bar"></a>
<a href="tel:55555555"></a>
<img src="/img/foo.png">
<link rel="stylesheet" href="http://foo.example/foo">
<a href="#skip-to-content"></a>
<a href="#/foo-bar-baz"></a>
`)

	links, err := ParseLinks(r)
	if err != nil {
		t.Error(err)
	}

	if len(links) != 4 {
		t.Errorf("fragments should be ignored; expected 4 links, got %d", len(links))
	}
}

func TestDuplicateLinks(t *testing.T) {
	r := strings.NewReader(`
<a href="http://foo.example/bar"></a>
<a href="http://foo.example/bar"></a>
<a href="http://foo.example/bar"></a>
<a href="http://foo.example/bar"></a>
<a href="http://foo.example/bar"></a>
<a href="http://foo.example/bar"></a>
<a href="http://foo.example/bar"></a>
<a href="http://foo.example/bar"></a>
<a href="http://foo.example/bar"></a>
`)

	links, err := ParseLinks(r)
	if err != nil {
		t.Error(err)
	}

	if len(links) != 1 {
		t.Error("duplicate links should be removed by the parser.")
	}
}
