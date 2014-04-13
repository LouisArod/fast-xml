// Package lxml contain tool to generate xml files on the fly.
package lxml

import (
	"bytes"
	"errors"
	"io"
)

// Predifined constantes.
const (
	XML1     = "<?xml version=\"1.0\" encoding=\"utf-8\"?>"
	SvgDoc   = "<!DOCTYPE svg PUBLIC \"-//W3C//DTD SVG 20010904//EN\" \"http://www.w3.org/TR/2001/REC-SVG-20010904/DTD/svg10.dtd\">"
	MaxDepth = 64
)

// Predifined errors
var (
	ErrDepthOverflow = errors.New("lxml.Gen: too deep.")
)

// A Gen represente an xml generator.
type Gen struct {
	b      *bytes.Buffer    // Contain the xml
	opened [MaxDepth]string // Contain the name of the opened elements
	depth  int8             // Contain the number of opened elements
	added  bool             // Indicate if the last element has been added or opened
}

// Create a new Gen with the specified Document type doc.
func NewGen(version, doctype string) *Gen {
	return &Gen{b: bytes.NewBufferString(version + "\n" + doctype + "\n")}
}

// AddNode add a closed node to the generator.
func (g *Gen) AddNode(name, attr string) (e error) {
	if g.depth >= MaxDepth {
		return ErrDepthOverflow
	}
}

// AddAttr add an attribute to the last opened node.
func (g *Gen) AddAttr(name, value string) {

}

// CloseNode close the last opened node in the generator.
func (g *Gen) CloseNode() {

}

// CloseNamedNode close all the nodes starting from
// the last opened node until it reach the node matching
// name.
func (g *Gen) CloseNamedNode(name string) {

}

// OpenNode open a new node named name with the attributes
// contains in attr.
func (g *Gen) OpenNode(name, attr string) {

}

// Read read len(p) byte from the generator.
func (g *Gen) Read(p []byte) (n int, err error) {
	return n, err
}

// Write the content of the generator until the end of the document.
func (g *Gen) WriteTo(w io.Writer) (n int64, err error) {
	return n, err
}
