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
	SVGDoc   = "<!DOCTYPE svg PUBLIC \"-//W3C//DTD SVG 1.1//EN\"\n\"http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd\">"
	MaxDepth = 127
)

// Predifined errors
var (
	ErrDepthOverflow  = errors.New("lxml.Gen: too deep.")
	ErrNothingToClose = errors.New("lxml.Gen: nothing to close.")
)

// A Gen represente an xml generator.
type Gen struct {
	b      bytes.Buffer     // Contain the opened xml
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
	g.depth--
	if g.added {
		_, err = g.b.WriteString(" /")
		g.added = false
	} else {
		_, err = g.b.WriteString("></" + g.opened[g.depth])
	}
	return err
}

// OpenNode open a new node named name with the attributes
// contains in attr.
func (g *Gen) OpenNode(name, attr string) (err error) {
	if g.depth >= MaxDepth {
		return ErrDepthOverflow
	}
	// Close the previous node if added
	buf := make([]byte, MaxDepth*4+5)

	n := copy(buf, ">\n")
	if g.added {
		if g.depth <= 0 {
			return ErrNothingToClose
		}
		n = copy(buf, " />\n")
		g.added = false
		g.depth--
	}
	// Add to opened nodes
	g.opened[g.depth] = name
	g.depth++

	// Fill the buffer with depth-1 tabs (no tab for the root element)
	for i, max := 0, int(g.depth-1); i < max; i++ {
		n += copy(buf[n:], "    ")
	}

	// Write the node in the buffer
	g.b.Write(buf[:n])
	g.b.WriteString("<" + name + " " + attr)

	return err
}

// AddNode add a closed node to the generator.
// Convenience function avoiding to call CloseNode for
// each element without childs.
func (g *Gen) AddNode(name, attr string) (err error) {
	err = g.OpenNode(name, attr)
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

	// We allocate a buffer able to contain 4 space for
	// each opened element + 6 char.
	buf := make([]byte, MaxDepth*4+6)
	buf[0], buf[1] = '>', '\n'
	n := 2

	// We close the previous element if added
	if g.added {
		if g.depth <= 1 {
			return ErrNothingToClose
		}
		buf[0], buf[1], buf[2], buf[3] = ' ', '/', '>', '\n'
		n = 4
		g.depth--
		g.added = false
	}

	// Fill the buffer with depth-1 tabs (no tab for the root element)
	for i, max := 0, int(g.depth-1); i < max; i++ {
		n += copy(buf[n:], "    ")
	}
	buf[n], buf[n+1] = '<', '/'
	n += 2

	g.b.Write(buf[:n])
	g.b.WriteString(g.opened[g.depth-1])
	g.depth--
	return nil
}

// CloseNNode close n nodes starting from
// the last opened node.
func (g *Gen) CloseNNode(n int) error {
	if g.depth-int8(n) <= 0 {
		return ErrNothingToClose
	}
	for i := 0; i < n; i++ {
		g.CloseNode()
	}
	return nil
}

// Read read len(p) byte from the generator.
func (g *Gen) Read(p []byte) (n int, err error) {
	if g.b.Len() == 0 {
		return
	}
	// We copy the opened XML
	n = copy(p, g.b.Bytes())
	// The slice is already full ? We escape.
	if l := len(p); l < g.b.Len() {
		return
	}

	// We create the closing elements
	// We handle added element
	i := int(g.depth - 1)
	if g.added {
		p[n], p[n+1] = ' ', '/'
		n += 2
		i--
	}

	// We create a string fill with space
	buf := make([]byte, MaxDepth*4+5)
	buf[0], buf[1] = '>', '\n'
	b := 2
	for j := 0; j < i; j++ {
		b += copy(buf[b:], "    ")
	}
	buf[b], buf[b+1] = '<', '/'
	b += 2
	// We close all the other element
	for ; err == nil && i > 0; i-- {
		n += copy(p[n:], buf[:b])
		n += copy(p[n:], g.opened[i])
		// We prepare the next iteration
		b -= 4
		buf = buf[:b]
		buf[b-2], buf[b-1] = '<', '/'
	}
	n += copy(p[n:], buf[:b])
	n += copy(p[n:], g.opened[i])

	if g.depth > 0 {
		p[n] = '>'
		n++
		return n, err
	}
	return
}

// Write the content of the generator until the end of the document.
func (g *Gen) WriteTo(w io.Writer) (n int64, err error) {
	return n, err
}
