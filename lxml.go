// Package lxml contain tool to generate xml files on the fly.
package lxml

import (
	"bytes"
	"errors"
	"io"
	"strings"
)

// Predifined constantes.
const (
	XML1     = "<?xml version=\"1.0\" encoding=\"utf-8\"?>"
	SvgDoc   = "<!DOCTYPE svg PUBLIC \"-//W3C//DTD SVG 20010904//EN\" \"http://www.w3.org/TR/2001/REC-SVG-20010904/DTD/svg10.dtd\">"
	MaxDepth = 64
)

// Predifined errors
var (
	ErrDepthOverflow  = errors.New("lxml.Gen: too deep.")
	ErrNothingToClose = errors.New("lxml.Gen: nothing to close.")
)

// A Gen represente an xml generator.
type Gen struct {
	b      bytes.Buffer     // Contain the xml
	opened [MaxDepth]string // Contain the name of the opened elements
	depth  int8             // Contain the number of opened elements
	added  bool             // Indicate if the last element has been added or opened
}

// Create a new Gen with the specified Document type doc.
func NewGen(version, doctype, root string) *Gen {
	//g := &Gen{b: bytes.NewBufferString(version + "\n" + doctype + "\n<" +
	//	root)}
	g := new(Gen)
	g.b.WriteString(version + "\n" + doctype + "\n<" + root)
	g.opened[0] = root
	g.depth++
	return g
}

func (g *Gen) closeLastNode() (err error) {
	if g.depth <= 0 {
		return ErrNothingToClose
	}
	if g.added {
		_, err = g.b.WriteString(" /")
		g.added = false
	} else {
		_, err = g.b.WriteString("></" + g.opened[g.depth-1])
	}
	g.depth--
	return err
}

// AddNode add a closed node to the generator.
// Convenience function avoiding to call CloseNode for
// each element without childs.
func (g *Gen) AddNode(name, attr string) (err error) {
	err = g.OpenNode(name, attr)
	if err != nil {
		return err
	}
	g.added = true
	return err
}

// AddAttr add an attribute to the last opened node.
func (g *Gen) AddAttr(a string) {
	g.b.WriteString(" " + a)
}

// CloseNode close the last opened node in the generator.
func (g *Gen) CloseNode() error {
	if g.depth <= 0 {
		return ErrNothingToClose
	}

	delim := ">\n"
	if g.added {
		if g.depth <= 1 {
			return ErrNothingToClose
		}
		delim = " />\n"
		g.depth--
		g.added = false
	}

	g.b.WriteString(delim + strings.Repeat("    ", int(g.depth-1)) +
		"</" + g.opened[g.depth-1])
	g.depth--
	return nil
}

// CloseNamedNode close all the nodes starting from
// the last opened node until it reach the node matching
// name.
func (g *Gen) CloseNamedNode(name string) {

}

// OpenNode open a new node named name with the attributes
// contains in attr.
func (g *Gen) OpenNode(name, attr string) (err error) {
	if g.depth >= MaxDepth {
		return ErrDepthOverflow
	}
	// Close the previous node if added
	delim := ">\n"
	if g.added {
		if g.depth <= 0 {
			return ErrNothingToClose
		}
		delim = " />\n"
		g.added = false
		g.depth--
	}
	// Add to opened nodes
	g.opened[g.depth] = name
	g.depth++

	// Write the node in the buffer
	_, err = g.b.WriteString(delim + strings.Repeat("    ", int(g.depth-1)) +
		"<" + name + " " + attr)
	return err
}

// Read read len(p) byte from the generator.
func (g *Gen) Read(p []byte) (n int, err error) {
	// If the buffer is smaller than the xml we copy what we can
	if l := len(p); l < g.b.Len() {
		n = copy(p, g.b.Bytes()[:l])
		return
	}
	n = copy(p, g.b.Bytes())
	var b bytes.Buffer
	// We handle added element
	i := g.depth - 1
	if g.added {
		_, err = b.WriteString(" /")
		i--
	}
	// We close all the other element
	for ; err == nil && i >= 0; i-- {
		_, err = b.WriteString(">\n" + strings.Repeat("    ", int(i)) + "</" + g.opened[i])
	}
	if g.depth > 0 {
		b.WriteString(">")
		tmp, err := b.Read(p[n:])
		return n + tmp, err
	}
	return
}

// Write the content of the generator until the end of the document.
func (g *Gen) WriteTo(w io.Writer) (n int64, err error) {
	return n, err
}
