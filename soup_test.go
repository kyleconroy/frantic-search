package main

import (
	"code.google.com/p/go.net/html"
	"strings"
	"testing"
)

func TestFlatten(t *testing.T) {
	s := `<p>Links:</p><ul><li><a href="foo">Foo</a><li><a href="/bar/baz">BarBaz</a></ul>`
	doc, _ := html.Parse(strings.NewReader(s))
	if Flatten(doc) != "Links:FooBarBaz" {
		t.Fatalf("%s was wrong", Flatten(doc))
	}
}

func TestFind(t *testing.T) {
	s := `<p>Links:</p><ul><li><a href="foo">Foo</a><li><a class="goo" href="/bar/baz">BarBaz</a></ul>`
	doc, _ := html.Parse(strings.NewReader(s))

	_, found := Find(doc, "#foo")

	if found {
		t.Errorf("There is no node with id 'foo'")
	}

	p, found := Find(doc, "p")

	if !found || p.Data != "p" {
		t.Errorf("Couldn't find p")
	}

	a, found := Find(doc, "ul a")

	if !found || a.Data != "a" || Flatten(a) != "Foo" {
		t.Errorf("Couldn't find a")
	}

	goo, found := Find(doc, "ul .goo")

	if !found || goo.Data != "a" || Flatten(goo) != "BarBaz" {
		t.Errorf("Couldn't find a with class goo")
	}

}

func TestFindAll(t *testing.T) {
	s := `<ul><li><a href="foo">Foo</a><li><a class="goo" href="/bar/baz">BarBaz</a></ul>`
	doc, _ := html.Parse(strings.NewReader(s))

	links := FindAll(doc, "li a")

	if len(links) != 2 {
		t.Errorf("There should be 2 link nodes")
	}
}
