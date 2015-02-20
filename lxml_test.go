package lxml

import (
	"bytes"
	"fmt"
	"io/ioutil"
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

var NewGenTest = "<svg"

func TestNewGen(t *testing.T) {
	g := NewGen("svg")
	if g == nil {
		t.Fatal("lxml.NewGen() has return 'nil'.")
	}
	if !bytes.Equal(g.xml.Bytes(), []byte(NewGenTest)) {
		t.Errorf("Expecting '%s' but got '%s'",
			NewGenTest, string(g.xml.Bytes()))
	}
	if g.depth != 1 {
		t.Fatal("Expecting a depth of 1 but got ", g.depth)
	}
	if string(g.opened[g.depth-1]) != "svg" {
		t.Fatal("Expecting 'svg' in g.opened bug got ", g.opened[g.depth-1])

	}
}

var AddNodeTest = `<svg>
    <rect width="100%" height='100%' />
    <rect width="100%" height='100%'`

func TestAddNode(t *testing.T) {
	g := NewGen("svg")
	g.AddNode("rect", "width=\"100%\" height='100%'")
	g.AddNode("rect", "width=\"100%\" height='100%'")
	assertByte(g.xml.Bytes(), []byte(AddNodeTest), t)
}

var AddAttrTest = `<svg>
    <rect width="100%" height='100%' style='fill: blue' id='myRect'`

func TestAddAttr(t *testing.T) {
	g := NewGen("svg")
	g.AddNode("rect", "width=\"100%\" height='100%'")
	g.AddAttr("style='fill: blue'")
	g.AddAttr("id='myRect'")
	assertByte(g.xml.Bytes(), []byte(AddAttrTest), t)
}

var OpenNodeTest = `<svg>
    <rect width="100%" height='100%'>
        <rect width="100%" height='100%'`

func TestOpenNode(t *testing.T) {
	g := NewGen("svg")
	g.OpenNode("rect", "width=\"100%\" height='100%'")
	g.OpenNode("rect", "width=\"100%\" height='100%'")
	assertByte(g.xml.Bytes(), []byte(OpenNodeTest), t)
}

var CloseNodeTest = `<svg>
    <rect width="100%" height='100%'>
        <rect width="100%" height='100%' />
    </rect>
    <rect width="100%" height='100%'`

func TestCloseNode(t *testing.T) {
	g := NewGen("svg")
	g.OpenNode("rect", "width=\"100%\" height='100%'")
	g.AddNode("rect", "width=\"100%\" height='100%'")
	g.CloseNode()
	g.OpenNode("rect", "width=\"100%\" height='100%'")
	assertByte(g.xml.Bytes(), []byte(CloseNodeTest), t)
}

var CloseNNodeTest = `<svg>
    <g class='bench'>
        <g class='bench'>
            <g class='bench'>
                <rect x='0' y='0' width="100%" height='100%' />
            </g>
        </g>
    </g>
    <rect x='0' y='0' width="100%" height='100%'`

func TestCloseNNode(t *testing.T) {
	g := NewGen("svg")
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
	assertByte(g.xml.Bytes(), []byte(CloseNNodeTest), t)
}

func TestGenerate(t *testing.T) {
	g := NewGen("svg")
	g.generate()
	t.Log("Gen.closing:", string(g.closing.Bytes()))
}

var (
	ReadSvgTest  = "<svg>\n</svg>"
	ReadSvgTest2 = `<svg>
    <g id='ImUnique'>
        <rect x='0' y='0' width="100%" height='100%' />
    </g>
    <rect x='0' y='0' width="100%" height='100%' />
    <g class='bench'>
        <g class='bench'>
            <rect x='0' y='0' width="100%" height='100%' />
        </g>
    </g>
</svg>`
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
	gg := NewGen("svg")
	n, _ = gg.Read(p)
	if !bytes.Equal(p[:n], []byte(ReadSvgTest)) {
		t.Errorf("Expecting '%s' but got '%s'",
			ReadSvgTest, string(p[:n]))
	}
	// Test complex Gen
	ggg := NewGen("svg")
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

func TestWriteTo(t *testing.T) {
	buf := &bytes.Buffer{}
	buf.Grow(2048)

	assertErr := func(err error) {
		if err != nil {
			t.Error(err)
		}
	}
	assertLen := func(n, dn int64) {
		if n != dn {
			t.Errorf("Expecting length of '%d' but got '%d'", dn, n)
		}
	}

	// Test empty Gen
	var g Gen
	n, err := g.WriteTo(buf)
	if err != ErrEmpty {
		t.Errorf("Expecting error '%v' but got '%v'", ErrEmpty, err)
	}
	buf.Reset()

	// Test simple Gen
	gg := NewGen("svg")
	n, err = gg.WriteTo(buf)
	assertErr(err)
	assertLen(n, int64(len(ReadSvgTest)))
	if !bytes.Equal(buf.Bytes(), []byte(ReadSvgTest)) {
		t.Errorf("Expecting '%s' but got '%s'",
			ReadSvgTest, string(buf.Bytes()))
	}
	buf.Reset()

	// Test complex Gen
	ggg := NewGen("svg")
	ggg.OpenNode("g", "id='ImUnique'")
	ggg.AddNode("rect", "x='0' y='0' width=\"100%\" height='100%'")
	ggg.CloseNode()
	ggg.AddNode("rect", "x='0' y='0' width=\"100%\" height='100%'")
	ggg.OpenNode("g", "class='bench'")
	ggg.OpenNode("g", "class='bench'")
	ggg.AddNode("rect", "x='0' y='0' width=\"100%\" height='100%'")

	n, err = ggg.WriteTo(buf)
	assertErr(err)
	assertLen(n, int64(len(ReadSvgTest2)))
	assertByte(buf.Bytes(), []byte(ReadSvgTest2), t)

}

/*func BenchmarkEq(b *testing.B) {
	buf := make([]byte, 2048)
	n, m := 0, 0
	for i := 0; i < b.N; i++ {
		buf[m], buf[m+1], buf[m+2] = ':', ')', ' '
		n += 3
		m = n % (2048 - 3)
	}
}

func BenchmarkCopy(b *testing.B) {
	buf := make([]byte, 2048)
	n, m := 0, 0
	for i := 0; i < b.N; i++ {
		n += copy(buf[m:], ":) ")
		m = n % (2048 - 3)
	}
}*/

func BenchmarkAddNode(b *testing.B) {
	g := NewGen("svg")
	g.AddAttr("width='100%' height='100%' preserveAspectRatio='xMidYMid meet' viewBox='0 0 400 400' xmlns='http://www.w3.org/2000/svg'")
	g.AddNode("rect", "x='0' y='0' width=\"400\" height='400' style='fill: white; stroke: black; stroke-width: 1px;'")

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		g.AddNode("rect", "x='10' y='10' width=\"50\" height='50'")
	}
}

func BenchmarkAddAttr(b *testing.B) {
	g := NewGen("svg")
	g.AddAttr("width='100%' height='100%' preserveAspectRatio='xMidYMid meet' viewBox='0 0 400 400' xmlns='http://www.w3.org/2000/svg'")
	g.AddNode("rect", "x='0' y='0' width=\"400\" height='400' style='fill: white; stroke: black; stroke-width: 1px;'")

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		g.AddAttr("style='fill: rgb(2,3,4);'")
	}
}

func BenchmarkRead(b *testing.B) {
	g := NewGen("svg")
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

func BenchmarkWriteTo(b *testing.B) {
	g := NewGen("svg")
	g.AddAttr("width='100%' height='100%' preserveAspectRatio='xMidYMid meet' viewBox='0 0 400 400' xmlns='http://www.w3.org/2000/svg'")
	g.AddNode("rect", "x='0' y='0' width=\"400\" height='400' style='fill: white; stroke: black; stroke-width: 1px;'")

	for i := 0; i < 15; i++ {
		g.OpenNode("g", "")
		g.AddNode("rect", fmt.Sprintf("x='%d' y='%d' width=\"%d\" height='%d'", i*10, i*10, 400-i*20, 400-i*20))
		g.AddAttr(fmt.Sprintf("style='fill: rgb(%d,%d,%d);'", i*10, 0, i*20))
	}

	var n int64
	var err error
	buf := new(bytes.Buffer)
	buf.Grow(32768)
	//buf := make([]byte, 32768)
	//b.Logf("g.depth: %d, g.b.Len(): %d", g.depth, g.b.Len())
	b.ReportAllocs()
	for i := 0; err == nil && i < b.N; i++ {
		//b.Logf("g.depth: %d, g.b.Len(): %d", g.depth, g.b.Len())
		buf.Reset()
		n, err = g.WriteTo(buf)
		b.SetBytes(n)
	}
	if err != nil {
		b.Fatal("Gen.WriteTo() has return the error:", err)
	}
	if n <= 0 {
		b.Fatal("Gen.WriteTo() expected to write more than 0 bytes and wrote :", n)
	}
	ioutil.WriteFile("benchmarkRead.svg", buf.Bytes(), 0755)
}

func BenchmarkGenerate(b *testing.B) {
	g := NewGen("svg")
	g.AddAttr("width='100%' height='100%' preserveAspectRatio='xMidYMid meet' viewBox='0 0 400 400' xmlns='http://www.w3.org/2000/svg'")
	g.AddNode("rect", "x='0' y='0' width=\"400\" height='400' style='fill: white; stroke: black; stroke-width: 1px;'")

	for i := 0; i < 15; i++ {
		g.OpenNode("g", "")
		g.AddNode("rect", fmt.Sprintf("x='%d' y='%d' width=\"%d\" height='%d'", i*10, i*10, 400-i*20, 400-i*20))
		g.AddAttr(fmt.Sprintf("style='fill: rgb(%d,%d,%d);'", i*10, 0, i*20))
	}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		g.generate()
	}
}
