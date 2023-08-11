package parser

import (
	"fmt"
	"github.com/goccy/go-graphviz"
	"strings"

	"github.com/goccy/go-graphviz/cgraph"
	"github.com/golang/glog"
)

var visNode map[string]*cgraph.Node

const (
	GrammarProduction = 1
	GrammarTerminal   = 2
	GrammarInner      = 3
)

type Result struct {
	path   []*Grammar
	output []string
}

func NewResult() *Result {
	return &Result{
		path: []*Grammar{},
	}
}

func (r *Result) AddNode(n *Grammar) *Result {
	r.path = append(r.path, n)
	return r
}

func (r *Result) AddOutput(s string) *Result {
	r.output = append(r.output, s)
	return r
}

func (r *Result) GetPath() []*Grammar {
	return r.path
}

func (r *Result) GetOutput() []string {
	return r.output
}

func (r *Result) Visualize(filename string) {
	gh := graphviz.New()
	graph, _ := gh.Graph()

	var prev *cgraph.Node = nil
	var outputIdx int = 0
	for idx, n := range r.GetPath() {
		current, err := graph.CreateNode(n.GetID())
		if err != nil {
			glog.Fatalf("something unexpected when noding %s: %v", n.GetID(), err)
		}
		if n.Type == GrammarTerminal {
			// current.Set(r.GetOutput()[outputIdx])
			outputNode, err := graph.CreateNode(r.GetOutput()[outputIdx])
			if err != nil {
				glog.Fatalf("something unexpected when outputing %s: %v", n.GetID(), err)
			}
			_, err = graph.CreateEdge("", current, outputNode)
			if err != nil {
				glog.Fatalf("something unexpected when connecting output %s: %v", n.GetID(), err)
			}
			outputIdx++
		}
		if prev != nil {
			e, err := graph.CreateEdge("", prev, current)
			if err != nil {
				glog.Fatalf("something unexpected when labeling %s: %v", n.GetID(), err)

			}
			e.SetLabel(r.GetPath()[idx-1].GetOperator().GetText())
		}
		prev = current
	}

	err := gh.RenderFilename(graph, graphviz.PNG, filename)
	if err != nil {
		glog.Fatalf("something unexpected happened when saving %s: %v", filename, err)
	}
}

type Context struct {
	// TODO: extend
	Operators map[string]Operator
}

func NewContext() *Context {
	return &Context{
		Operators: map[string]Operator{},
	}
}

func (ctx *Context) GetOperator(id string) Operator {
	if op, ok := ctx.Operators[id]; ok {
		return op
	}
	return nil
}

func (ctx *Context) UpdateOperator(id string, op Operator) {
	ctx.Operators[id] = op
}

func (ctx *Context) Copy() *Context {
	new := NewContext()
	for k, v := range ctx.Operators {
		new.UpdateOperator(k, v)
	}
	return new
}

type Grammar struct {
	ID      string
	Type    int
	Symbols *[]*Grammar
	Content string
	Root    *Grammar
	Ctx     *Context
}

func NewGrammar(ctx *Context, t int, op Operator, id string, content string) *Grammar {
	ctx.Operators[id] = op
	new := &Grammar{
		ID:      id,
		Type:    t,
		Symbols: &[]*Grammar{},
		Content: content,
		Ctx:     ctx,
	}
	return new.SetRoot(new)
}

func (g *Grammar) SetID(id string) *Grammar {
	if g.GetCtx() != nil {
		g.GetCtx().UpdateOperator(id, g.GetOperator())
	}
	g.ID = id
	return g
}

func (g *Grammar) GetID() string {
	return g.ID
}

func (g *Grammar) SetRoot(root *Grammar) *Grammar {
	g.Root = root
	return g
}

func (g *Grammar) GetRoot() *Grammar {
	return g.Root
}

func (g *Grammar) GetCtx() *Context {
	return g.GetRoot().Ctx
}

func (g *Grammar) SetCtx(ctx *Context) *Grammar {
	g.Ctx = ctx
	return g
}

func (g *Grammar) AddSymbol(new *Grammar) int {
	if new.Type != GrammarProduction {
		new.Root = g.Root
	}
	*g.Symbols = append(*g.Symbols, new)
	return len(*g.Symbols) - 1
}

func (g *Grammar) GetSymbol(idx int) *Grammar {
	if idx < len(*g.Symbols) {
		return (*g.Symbols)[idx]
	}
	return nil
}

func (g *Grammar) SetSymbol(idx int, sym *Grammar) *Grammar {
	if idx < len(*g.Symbols) {
		(*g.Symbols)[idx] = sym
	} else {
		glog.Fatal("idx out of range")
	}
	return g
}

func (g *Grammar) GetOperator() Operator {
	return g.GetRoot().GetCtx().GetOperator(g.ID)
}

func (g *Grammar) SetOperator(op Operator) {
	g.Ctx.UpdateOperator(g.ID, op)
}

func (g *Grammar) GetContent() string {
	return g.Content
}

func (g *Grammar) ShallowCopy() *Grammar {
	new := *g
	return &new
}

func (g *Grammar) ForkContext(newID string) (*Grammar, error) {
	if g.Type != GrammarProduction {
		return nil, fmt.Errorf("it is not necessary to fork a non-production grammar: %s", g.ID)
	}
	new := g.ShallowCopy()
	new = new.SetCtx(g.GetCtx().Copy())
	if newID != "" {
		new = new.SetID(newID)
		new.GetCtx().UpdateOperator(newID, g.GetOperator())
	}
	return new, nil
}

func (g *Grammar) Generate(r *Result) *Result {
	if r == nil {
		r = &Result{}
	}
	op := g.GetOperator()
	if op != nil {
		op.BeforeGen(g.GetCtx(), g, r)
		op.Gen(g.GetCtx(), g, r)
		op.AfterGen(g.GetCtx(), g, r)
	}
	return r
}

func (g *Grammar) Visualize(filename string, expandSub bool) {
	gh := graphviz.New()
	graph, _ := gh.Graph()
	visNode = make(map[string]*cgraph.Node)
	g.addNodeToGraph(graph, nil, g.GetOperator().GetText(), expandSub)

	err := gh.RenderFilename(graph, graphviz.PNG, filename)
	if err != nil {
		panic(err)
	}
}

func (g *Grammar) addNodeToGraph(graph *cgraph.Graph, parent *cgraph.Node, label string, expandSub bool) {
	var n *cgraph.Node
	visited := false
	fmt.Println(g.ID)
	if strings.Contains(g.ID, "EOF") {
		return
	}
	if node, ok := visNode[g.ID]; ok {
		n = node
		visited = true
		return
	} else {
		node, err := graph.CreateNode(g.ID)
		if err != nil {
			panic(err)
		}
		n = node
		visNode[g.ID] = n
	}

	if parent != nil {
		edge, err := graph.CreateEdge(label, parent, n)
		if err != nil {
			panic(err)
		}
		edge.SetLabel(label)
	}
	if (parent != nil && g.Type == GrammarProduction && !expandSub) || visited {
		return
	}
	for _, child := range *g.Symbols {
		child.addNodeToGraph(graph, n, g.GetOperator().GetText(), expandSub)
	}
}

func (g *Grammar) ConvertCypher() {

}
