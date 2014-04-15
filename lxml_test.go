package lxml

import (
	"bytes"
	"testing"
)

// equalCount count the number of equal's bytes.
func equalCount(a, b []byte) int {
	i := 0
	for i = 0; i < len(a) && i < len(b) && a[i] == b[i]; i++ {
	}
	return i
}

// assertByte test two []byte agaist each other
// and call t.Error if there is a difference.
func assertByte(gen, res []byte, t *testing.T) {
	n := equalCount(gen, res)
	if n != len(res) {
		subGen := string(gen[max(0, n-5):min(len(gen), n+5)])
		subRes := string(res[max(0, n-5):min(len(res), n+5)])
		t.Errorf("Expecting:\n'%s'\nbut got:\n'%s'"+
			"\n\nComparaison failed on char %d with: '%s'\n"+
			"Sould have got: '%s'",
			string(res), string(gen), n, subGen, subRes)
	}

}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

var (
	ReadSvgTest  = XML1 + "\n" + SvgDoc + "\n" + "<svg>\n</svg>"
	ReadSvgTest2 = XML1 + "\n" + SvgDoc + "\n" +
		"<svg>\n" +
		"    <rect width=\"100%\" height='100%'>\n" +
		"        <rect width=\"100%\" height='100%' />\n" +
		"    </rect>\n" +
		"    <rect width=\"100%\" height='100%' />\n" +
		"</svg>"
)

func TestRead(t *testing.T) {
	var g Gen
	// Test empty Gen
	p := make([]byte, 512)
	n, _ := g.Read(p)
	if string(p[:n]) != "" {
		t.Errorf("Expecting '%s' but got '%s'", "", string(p[:n]))
	}
	// Test simple Gen
	gg := NewGen(XML1, SvgDoc, "svg")
	n, _ = gg.Read(p)
	if !bytes.Equal(p[:n], []byte(ReadSvgTest)) {
		t.Errorf("Expecting '%s' but got '%s'",
			ReadSvgTest, string(p[:n]))
	}
	// Test complex Gen
	ggg := NewGen(XML1, SvgDoc, "svg")
	ggg.OpenNode("rect", "width=\"100%\" height='100%'")
	ggg.AddNode("rect", "width=\"100%\" height='100%'")
	ggg.CloseNode()
	ggg.AddNode("rect", "width=\"100%\" height='100%'")
	n, _ = ggg.Read(p)
	assertByte(p, []byte(ReadSvgTest2), t)
}

var NewGenTest = XML1 + "\n" + SvgDoc + "\n" + "<svg"

func TestNewGen(t *testing.T) {
	g := NewGen(XML1, SvgDoc, "svg")
	if g == nil {
		t.Fatal("lxml.NewGen() has return 'nil'.")
	}
	if !bytes.Equal(g.b.Bytes(), []byte(NewGenTest)) {
		t.Errorf("Expecting '%s' but got '%s'",
			NewGenTest, string(g.b.Bytes()))
	}
	if g.depth != 1 {
		t.Fatal("Expecting a depth of 1 but got %d", g.depth)
	}
	if g.opened[g.depth-1] != "svg" {
		t.Fatal("Expecting 'svg' in g.opened bug got %s", g.opened[g.depth-1])

	}
}

var AddNodeTest = XML1 + "\n" + SvgDoc + "\n" +
	"<svg>\n" +
	"    <rect width=\"100%\" height='100%' />\n" +
	"    <rect width=\"100%\" height='100%'"

func TestAddNode(t *testing.T) {
	g := NewGen(XML1, SvgDoc, "svg")
	g.AddNode("rect", "width=\"100%\" height='100%'")
	g.AddNode("rect", "width=\"100%\" height='100%'")
	assertByte(g.b.Bytes(), []byte(AddNodeTest), t)
}

var AddAttrTest = XML1 + "\n" + SvgDoc + "\n" +
	"<svg>\n" +
	"    <rect width=\"100%\" height='100%' style='fill: blue' id='myRect'"

func TestAddAttr(t *testing.T) {
	g := NewGen(XML1, SvgDoc, "svg")
	g.AddNode("rect", "width=\"100%\" height='100%'")
	g.AddAttr("style='fill: blue'")
	g.AddAttr("id='myRect'")
	assertByte(g.b.Bytes(), []byte(AddAttrTest), t)
}

var OpenNodeTest = XML1 + "\n" + SvgDoc + "\n" +
	"<svg>\n" +
	"    <rect width=\"100%\" height='100%'>\n" +
	"        <rect width=\"100%\" height='100%'"

func TestOpenNode(t *testing.T) {
	g := NewGen(XML1, SvgDoc, "svg")
	g.OpenNode("rect", "width=\"100%\" height='100%'")
	g.OpenNode("rect", "width=\"100%\" height='100%'")
	assertByte(g.b.Bytes(), []byte(OpenNodeTest), t)
}

var CloseNodeTest = XML1 + "\n" + SvgDoc + "\n" +
	"<svg>\n" +
	"    <rect width=\"100%\" height='100%'>\n" +
	"        <rect width=\"100%\" height='100%' />\n" +
	"    </rect>\n" +
	"    <rect width=\"100%\" height='100%'"

func TestCloseNode(t *testing.T) {
	g := NewGen(XML1, SvgDoc, "svg")
	g.OpenNode("rect", "width=\"100%\" height='100%'")
	g.AddNode("rect", "width=\"100%\" height='100%'")
	g.CloseNode()
	g.OpenNode("rect", "width=\"100%\" height='100%'")
	assertByte(g.b.Bytes(), []byte(CloseNodeTest), t)
}

func TestCloseNamedNode(t *testing.T) {

}

func TestWriteTo(t *testing.T) {

}
