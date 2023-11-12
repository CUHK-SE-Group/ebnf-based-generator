package schemas

import (
	"github.com/CUHK-SE-Group/generic-generator/graph"
	A "github.com/IBM/fp-go/array"
	"github.com/google/uuid"
)

type GrammarType int

// 带yes标记的symbol要指定生成策略
const (
	GrammarProduction GrammarType = 1 << iota
	GrammarExpr                   // yes
	GrammarTerm
	GrammarPAREN
	GrammarBRACKET // yes
	GrammarBRACE   // yes
	GrammarREP     // yes
	GrammarPLUS    // yes
	GrammarEXT     // yes
	GrammarSUB     // yes
	GrammarID
	GrammarTerminal
)

type property struct {
	Type    GrammarType
	Root    *Node
	Gram    *Grammar
	Content string
}

const (
	prop = "property"
)

type Grammar struct {
	internal graph.Graph[property]
}

func NewGrammar() *Grammar {
	newG := &Grammar{
		internal: graph.NewGraph[property](),
	}
	return newG
}

type Node struct {
	internal graph.Vertex[property]
}

func newEdge(id string, from, to *Node) graph.Edge[property] {
	res := graph.NewEdge[property]()
	res.SetID(id)
	res.SetFrom(from.internal)
	res.SetTo(to.internal)
	return res
}

func NewNode(g *Grammar, tp GrammarType) *Node {
	n := graph.NewVertex[property]()
	n.SetProperty(prop, property{
		Type:    tp,
		Root:    nil,
		Gram:    g,
		Content: "",
	})
	return &Node{internal: n}
}

func (g *Node) GetType() GrammarType {
	return g.internal.GetProperty(prop).Type
}

func (g *Node) SetRoot(r *Node) {
	ori := g.internal.GetProperty(prop)
	ori.Root = r
	g.internal.SetProperty(prop, ori)
}

func (g *Node) GetRoot() *Node {
	return g.internal.GetProperty(prop).Root
}

func (g *Node) GetGrammar() *Grammar {
	return g.internal.GetProperty(prop).Gram
}

func (g *Node) AddSymbol(new *Node) int {
	e := newEdge(uuid.NewString(), g, new)
	g.GetGrammar().internal.AddEdge(e)
	return len(g.GetGrammar().internal.GetOutEdges(g.internal)) - 1
}

func (g *Node) GetSymbols() []*Node {
	f := func(edge graph.Edge[property]) *Node {
		return &Node{internal: edge.GetTo()}
	}
	return A.Map(f)(g.GetGrammar().internal.GetOutEdges(g.internal))
}

func (g *Node) GetSymbol(idx int) *Node {
	syms := g.GetSymbols()
	if idx < len(syms) {
		return (syms)[idx]
	}
	return nil
}

func (g *Node) GetContent() string {
	return g.internal.GetProperty(prop).Content
}

func (g *Grammar) Visualize(filename string, expandSub bool) {
	//ghash := func(gram *Grammar) string { return gram.content }
	//gr := graph.New(ghash, graph.Directed(), graph.Rooted())
	//_ = gr.AddVertex(g)
	//
	//queue := make([]*Grammar, 0)
	//queue = append(queue, g)
	//visited := make(map[string]bool, 0)
	//for len(queue) != 0 {
	//	cur := queue[0]
	//	queue = queue[1:]
	//	visited[cur.id] = true
	//	for _, v := range *cur.GetSymbols() {
	//		_ = gr.AddVertex(v)
	//		_ = gr.AddEdge(cur.content, v.content)
	//		if visited[v.id] {
	//			continue
	//		}
	//		queue = append(queue, v)
	//	}
	//}
	//
	//file, _ := os.Create(filename)
	//_ = draw.DOT(gr, file)
}
