// Package lxml contain tools to generate xml files on the fly.
package lxml

import (
	"bytes"
	"errors"
	"io"
)

// Predifined constantes.
const (
	XML1        = "<?xml version=\"1.0\" encoding=\"utf-8\"?>"
	SVGDoc      = "<!DOCTYPE svg PUBLIC \"-//W3C//DTD SVG 1.1//EN\"\n\"http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd\">"
	MaxDepth    = 127
	MaxNodeName = 255
)

// Predifined errors
var (
	ErrDepthOverflow  = errors.New("lxml.Gen: too deep")
	ErrNothingToClose = errors.New("lxml.Gen: nothing to close")
	ErrEmpty          = errors.New("lxml.Gen: empty struct")
)

// A Gen represente an xml generator.
type Gen struct {
	xml     bytes.Buffer     // Contains the opened xml
	closing bytes.Buffer     // Contains the closing elements.
	opened  [MaxDepth][]byte // Contains the name of the opened elements
	depth   int8             // Contains the number of opened elements
	added   bool             // Indicate if the last element has been added or opened
	changed bool             // Indicate if there has been a change after the generation.
}

// NewGen create a new xml.Gen with the specified root element.
func NewGen(root string) *Gen {
	g := &Gen{added: false, changed: true}
	g.xml.WriteString("<" + root)
	g.opened[0] = []byte(root)
	g.depth++
	return g
}

// Close the previous element and prepare
// the new one with correct indentation and
// opening sign '<'.
func (g *Gen) closePreviousNode() {
	// Close the previous node if added
	buf := make([]byte, MaxDepth*4+5)

	buf[0], buf[1] = '>', '\n'
	n := 2
	if g.added {
		buf[0], buf[1], buf[2], buf[3] = ' ', '/', '>', '\n'
		n = 4
		g.added = false
	}

	// Fill the buffer with depth-1 tabs (no tab for the root element)
	for i, max := 0, int(g.depth); i < max; i++ {
		//n += copy(buf[n:], "    ")
		buf[n], buf[n+1], buf[n+2], buf[n+3] = ' ', ' ', ' ', ' '
		n += 4
	}

	// Write the opening sign
	buf[n] = '<'
	n++

	// Write the node in the buffer
	g.xml.Write(buf[:n])
	g.changed = true
}

// OpenNode open a new xml node.
func (g *Gen) OpenNode(name, attr string) (err error) {
	if g.depth >= MaxDepth {
		return ErrDepthOverflow
	}

	g.closePreviousNode()

	// Save the new element
	//fmt.Printf("OpenNode: depth=%v\n", g.depth)
	g.opened[g.depth] = []byte(name)
	g.depth++

	// Write the new element
	g.xml.WriteString(name + " " + attr)
	g.changed = true

	return err
}

// AddNode add a closed node to the generator.
// Convenience function avoiding to call CloseNode for
// each element without childs.
func (g *Gen) AddNode(name, attr string) {
	g.closePreviousNode()
	//fmt.Printf("AddNode: depth=%v\n", g.depth)

	// Write the new element.
	g.xml.WriteString(name + " " + attr)

	// Keep in mind how the node has been open.
	g.added = true
	g.changed = true
}

// AddAttr add an attribute to the last opened node.
func (g *Gen) AddAttr(a string) {
	g.xml.WriteString(" " + a)
	g.changed = true
}

// CloseNode close the last opened node in the generator.
// This function cannot close the root element, and
// shoul not be use to close elements added with Gen.AddNode().
func (g *Gen) CloseNode() error {
	if g.depth <= 1 {
		return ErrNothingToClose
	}

	// We go back from one level.
	g.depth--

	g.closePreviousNode()

	// We close the node.
	g.xml.WriteByte('/')
	g.xml.Write(g.opened[g.depth])
	g.changed = true

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

func (g *Gen) generate() {
	g.closing.Reset()

	// We close the previous element if added.
	if g.added {
		//p[n], p[n+1] = ' ', '/'
		g.closing.Write([]byte(" /"))
	}

	// We create a string fill with space
	buf := make([]byte, MaxDepth*4+5)
	buf[0], buf[1] = '>', '\n'
	b := 2
	for j, max := 0, int(g.depth-1); j < max; j++ {
		buf[b] = ' '
		buf[b+1] = ' '
		buf[b+2] = ' '
		buf[b+3] = ' '
		b += 4
	}
	buf[b], buf[b+1] = '<', '/'
	b += 2
	// We close all the other element
	for i := int(g.depth - 1); i > 0; i-- {
		g.closing.Write(buf[:b])
		g.closing.Write(g.opened[i])
		// We prepare the next iteration
		b -= 4
		buf = buf[:b]
		buf[b-2], buf[b-1] = '<', '/'
	}
	g.closing.Write(buf[:b])
	g.closing.Write(g.opened[0])
	g.closing.WriteByte('>')
	g.changed = false
}

// Read read len(p) byte from the generator.
func (g *Gen) Read(p []byte) (n int, err error) {
	if g.depth == 0 || g.xml.Len() == 0 {
		return n, ErrEmpty
	}

	// We copy the opened XML
	n = copy(p, g.xml.Bytes())
	if len(p) <= g.xml.Len() {
		return
	}

	// If the xml has changed we rewrite the closing elements.
	if g.changed {
		g.generate()
	}

	n += copy(p[n:], g.closing.Bytes())
	return
}

// WriteTo write the content of the generator until the end of the document.
func (g *Gen) WriteTo(w io.Writer) (int64, error) {
	if g.depth == 0 || g.xml.Len() == 0 {
		return 0, ErrEmpty
	}

	// We copy the opened XML
	n, err := w.Write(g.xml.Bytes())
	if err != nil {
		return int64(n), err
	}

	// We copy the closing tags.
	var tmp int
	if g.changed {
		g.generate()
	}
	tmp, err = w.Write(g.closing.Bytes())
	n += tmp

	return int64(n), err
}
