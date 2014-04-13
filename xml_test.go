package lxml_test

import (
	"testing"

	"github.com/LouisArod/lxml"
)

func TestNewGen(t *testing.T) {
	g := lxml.NewGen(lxml.XML1, lxml.SvgDoc)
	if g == nil {
		t.Fatal("lxml.NewGen() has return 'nil'.")
	}
	if string(g.Version) != lxml.XML1 {
		t.Errorf("Gen.Version is eq to '%s' and should be eq to '%s'",
			string(g.Version), lxml.XML1)
	}
	if string(g.Doctype) != lxml.SvgDoc {
		t.Errorf("Gen.Doctype is eq to '%s' and should be eq to '%s'",
			string(g.Doctype), lxml.SvgDoc)
	}
}

func TestAddNode(t *testing.T) {
	g = lxml.NewGen(lxml.XML1, lxml.SvgDoc)
	g.AddNode("svg", "width=\"100%\" height='100%'")
}

func TestAddAttr(t *testing.T) {

}

func TestCloseNode(t *testing.T) {

}

func TestCloseNamedNode(t *testing.T) {

}

func TestOpenNode(t *testing.T) {

}

func TestRead(t *testing.T) {

}

func TestWriteTo(t *testing.T) {

}
