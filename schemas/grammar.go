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

type Property struct {
	Type    GrammarType
	Root    *Node
	Gram    *Grammar
	Content string
}

const (
	Prop = "Property"
)

type Grammar struct {
	internal graph.Graph[string, Property]
}

func NewGrammar() *Grammar {
	newG := &Grammar{
		internal: graph.NewGraph[string, Property](),
	}
	return newG
}
func (g *Grammar) GetInternal() graph.Graph[string, Property] {
	return g.internal
}

func (g *Grammar) GetNode(id string) *Node {
	return &Node{internal: g.internal.GetVertexById(id)}
}

type Node struct {
	internal graph.Vertex[Property]
}

func newEdge(id string, from, to *Node) graph.Edge[string, Property] {
	res := graph.NewEdge[string, Property]()
	res.SetID(id)
	res.SetFrom(from.internal)
	res.SetTo(to.internal)
	return res
}

func NewNode(g *Grammar, tp GrammarType, id, content string) *Node {
	n := graph.NewVertex[Property]()
	n.SetProperty(Prop, Property{
		Type:    tp,
		Root:    nil,
		Gram:    g,
		Content: content,
	})
	n.SetID(id)
	return &Node{internal: n}
}

func (g *Node) GetType() GrammarType {
	return g.internal.GetProperty(Prop).Type
}

func (g *Node) GetID() string {
	return g.internal.GetID()
}

func (g *Node) SetRoot(r *Node) {
	ori := g.internal.GetProperty(Prop)
	ori.Root = r
	g.internal.SetProperty(Prop, ori)
}

func (g *Node) GetGrammar() *Grammar {
	return g.internal.GetProperty(Prop).Gram
}

func (g *Node) AddSymbol(new *Node) int {
	e := newEdge(uuid.NewString(), g, new)
	g.GetGrammar().internal.AddEdge(e)
	return len(g.GetGrammar().internal.GetOutEdges(g.internal)) - 1
}

func (g *Node) GetSymbols() []*Node {
	f := func(edge graph.Edge[string, Property]) *Node {
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
	return g.internal.GetProperty(Prop).Content
}
