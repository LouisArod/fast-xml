// Package xml contain tool to create xml files.
package lxml

import "io"

// Predifined Document constantes.
const (
	Xml1       = "<?xml version=\"1.0\" encoding=\"utf-8\"?>"
	SvgDoctype = "<!DOCTYPE svg PUBLIC \"-//W3C//DTD SVG 20010904//EN\" \"http://www.w3.org/TR/2001/REC-SVG-20010904/DTD/svg10.dtd\">"
)

// A Document represente the type and encoding of an XML document
type Document struct {
	Version string
	Doctype string
}

// A Gen represente an xml generator.
type Gen struct {
	Doc Document
}

// Create a new Gen with the specified Document type doc.
func NewGen(doc Document) *Gen {
	return &Gen{doc}
}

// AddNode add a closed node to the generator.
func (g *Gen) AddNode(name, attr string) {

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
