package lxml

import (
	"bytes"
	"fmt"
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
	ReadSvgTest  = XML1 + "\n" + SVGDoc + "\n" + "<svg>\n</svg>"
	ReadSvgTest2 = XML1 + "\n" + SVGDoc + "\n" +
		"<svg>\n" +
		"    <g id='ImUnique'>\n" +
		"        <rect x='0' y='0' width=\"100%\" height='100%' />\n" +
		"    </g>\n" +
		"    <rect x='0' y='0' width=\"100%\" height='100%' />\n" +
		"    <g class='bench'>\n" +
		"        <g class='bench'>\n" +
		"            <rect x='0' y='0' width=\"100%\" height='100%' />\n" +
		"        </g>\n" +
		"    </g>\n" +
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
	gg := NewGen(XML1, SVGDoc, "svg")
	n, _ = gg.Read(p)
	if !bytes.Equal(p[:n], []byte(ReadSvgTest)) {
		t.Errorf("Expecting '%s' but got '%s'",
			ReadSvgTest, string(p[:n]))
	}
	// Test complex Gen
	ggg := NewGen(XML1, SVGDoc, "svg")
	ggg.OpenNode("g", "id='ImUnique'")
	ggg.AddNode("rect", "x='0' y='0' width=\"100%\" height='100%'")
	ggg.CloseNode()
	ggg.AddNode("rect", "x='0' y='0' width=\"100%\" height='100%'")
	ggg.OpenNode("g", "class='bench'")
	ggg.OpenNode("g", "class='bench'")
	ggg.AddNode("rect", "x='0' y='0' width=\"100%\" height='100%'")

	n, _ = ggg.Read(p)
	assertByte(p, []byte(ReadSvgTest2), t)
}

var NewGenTest = XML1 + "\n" + SVGDoc + "\n" + "<svg"

func TestNewGen(t *testing.T) {
	g := NewGen(XML1, SVGDoc, "svg")
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

var AddNodeTest = XML1 + "\n" + SVGDoc + "\n" +
	"<svg>\n" +
	"    <rect width=\"100%\" height='100%' />\n" +
	"    <rect width=\"100%\" height='100%'"

func TestAddNode(t *testing.T) {
	g := NewGen(XML1, SVGDoc, "svg")
	g.AddNode("rect", "width=\"100%\" height='100%'")
	g.AddNode("rect", "width=\"100%\" height='100%'")
	assertByte(g.b.Bytes(), []byte(AddNodeTest), t)
}

var AddAttrTest = XML1 + "\n" + SVGDoc + "\n" +
	"<svg>\n" +
	"    <rect width=\"100%\" height='100%' style='fill: blue' id='myRect'"

func TestAddAttr(t *testing.T) {
	g := NewGen(XML1, SVGDoc, "svg")
	g.AddNode("rect", "width=\"100%\" height='100%'")
	g.AddAttr("style='fill: blue'")
	g.AddAttr("id='myRect'")
	assertByte(g.b.Bytes(), []byte(AddAttrTest), t)
}

var OpenNodeTest = XML1 + "\n" + SVGDoc + "\n" +
	"<svg>\n" +
	"    <rect width=\"100%\" height='100%'>\n" +
	"        <rect width=\"100%\" height='100%'"

func TestOpenNode(t *testing.T) {
	g := NewGen(XML1, SVGDoc, "svg")
	g.OpenNode("rect", "width=\"100%\" height='100%'")
	g.OpenNode("rect", "width=\"100%\" height='100%'")
	assertByte(g.b.Bytes(), []byte(OpenNodeTest), t)
}

var CloseNodeTest = XML1 + "\n" + SVGDoc + "\n" +
	"<svg>\n" +
	"    <rect width=\"100%\" height='100%'>\n" +
	"        <rect width=\"100%\" height='100%' />\n" +
	"    </rect>\n" +
	"    <rect width=\"100%\" height='100%'"

func TestCloseNode(t *testing.T) {
	g := NewGen(XML1, SVGDoc, "svg")
	g.OpenNode("rect", "width=\"100%\" height='100%'")
	g.AddNode("rect", "width=\"100%\" height='100%'")
	g.CloseNode()
	g.OpenNode("rect", "width=\"100%\" height='100%'")
	assertByte(g.b.Bytes(), []byte(CloseNodeTest), t)
}

var CloseNNodeTest = XML1 + "\n" + SVGDoc + "\n" +
	"<svg>\n" +
	"    <g class='bench'>\n" +
	"        <g class='bench'>\n" +
	"            <g class='bench'>\n" +
	"                <rect x='0' y='0' width=\"100%\" height='100%' />\n" +
	"            </g>\n" +
	"        </g>\n" +
	"    </g>\n" +
	"    <rect x='0' y='0' width=\"100%\" height='100%'"

func TestCloseNNode(t *testing.T) {
	g := NewGen(XML1, SVGDoc, "svg")
	n := 3
	for i := 0; i < n; i++ {
		g.OpenNode("g", "class='bench'")
	}
	g.AddNode("rect", "x='0' y='0' width=\"100%\" height='100%'")
	err := g.CloseNNode(n + 1)
	if err != ErrNothingToClose {
		t.Errorf("CloseNNode should have failed on error: %v\n"+
			"But we got: %v", ErrNothingToClose, err)
	}
	err = g.CloseNNode(n)
	if err != nil {
		t.Fatalf("Fail to close %d nodes with error: %v", n, err)
	}
	g.AddNode("rect", "x='0' y='0' width=\"100%\" height='100%'")
	assertByte(g.b.Bytes(), []byte(CloseNNodeTest), t)
}

func TestWriteTo(t *testing.T) {

}

func BenchmarkRead(b *testing.B) {
	g := NewGen(XML1, SVGDoc, "svg")
	g.AddAttr("width='100%' height='100%' preserveAspectRatio='xMidYMid meet' viewBox='0 0 400 400' xmlns='http://www.w3.org/2000/svg'")
	g.AddNode("rect", "x='0' y='0' width=\"400\" height='400' style='fill: white; stroke: black; stroke-width: 1px;'")

	for i := 0; i < 15; i++ {
		g.OpenNode("g", "")
		g.AddNode("rect", fmt.Sprintf("x='%d' y='%d' width=\"%d\" height='%d'", i*10, i*10, 400-i*20, 400-i*20))
		g.AddAttr(fmt.Sprintf("style='fill: rgb(%d,%d,%d);'", i*10, 0, i*20))
	}

	var n int
	var err error
	buf := make([]byte, 32768)
	b.ReportAllocs()
	for i := 0; err == nil && i < b.N; i++ {
		n, err = g.Read(buf)
		b.SetBytes(int64(n))
	}
	if err != nil {
		b.Fatal("Gen.Read() has return the error:", err)
	}
	if n == len(buf) {
		b.Error("Buffer overflow.")
	}
	//ioutil.WriteFile("benchmarkRead.svg", buf[:n], 0755)
}

func BenchmarkAddNode(b *testing.B) {
	g := NewGen(XML1, SVGDoc, "svg")
	g.AddAttr("width='100%' height='100%' preserveAspectRatio='xMidYMid meet' viewBox='0 0 400 400' xmlns='http://www.w3.org/2000/svg'")
	g.AddNode("rect", "x='0' y='0' width=\"400\" height='400' style='fill: white; stroke: black; stroke-width: 1px;'")

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		g.AddNode("rect", "x='10' y='10' width=\"50\" height='50'")
	}
}

func BenchmarkAddAttr(b *testing.B) {
	g := NewGen(XML1, SVGDoc, "svg")
	g.AddAttr("width='100%' height='100%' preserveAspectRatio='xMidYMid meet' viewBox='0 0 400 400' xmlns='http://www.w3.org/2000/svg'")
	g.AddNode("rect", "x='0' y='0' width=\"400\" height='400' style='fill: white; stroke: black; stroke-width: 1px;'")

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		g.AddAttr("style='fill: rgb(2,3,4);'")
	}
}

func BenchmarkOpenNode(b *testing.B) {
	g := NewGen(XML1, SVGDoc, "svg")
	g.AddAttr("width='100%' height='100%' preserveAspectRatio='xMidYMid meet' viewBox='0 0 400 400' xmlns='http://www.w3.org/2000/svg'")
	g.AddNode("rect", "x='0' y='0' width=\"400\" height='400' style='fill: white; stroke: black; stroke-width: 1px;'")

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		g.OpenNode("g", "class='benchmark'")
	}
}
