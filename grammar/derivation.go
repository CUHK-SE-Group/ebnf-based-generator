package grammar

import (
	"context"
	"fmt"
	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
)

type Handler func(ctx *DerivationContext, tree *DerivationTree)

type DerivationContext struct {
	context.Context
	coverage      map[string]int
	handlers      []Handler
	handlerIndex  int
	preExpansions []ExpansionPair
	finish        bool
}

func (ctx *DerivationContext) Next() Handler {
	if ctx.handlerIndex <= len(ctx.handlers) {
		ctx.handlerIndex++
	}
	return ctx.handlers[ctx.handlerIndex-1]
}

type DerivationTree struct {
	gram        Grammar
	startSymbol string
	traversal   TraversalAlgo
	expandAlgo  Handler
	root        *Node

	handlers       map[string][]Handler
	systemHandlers []Handler
}
type Node struct {
	ExpansionTuple
	Children []*Node
}
type ExpansionPair struct {
	From ExpansionTuple
	To   ExpansionTuple
}

func (dt *DerivationTree) Construct() {
	for key, value := range dt.gram.G {
		if key == dt.startSymbol {
			dt.root = &Node{ExpansionTuple: ExpansionTuple{name: key}}
			dt.root.Children = make([]*Node, 0)
			for _, v := range value {
				dt.root.Children = append(dt.root.Children, &Node{ExpansionTuple: v})
			}
		}
	}
}

// GetNonTerminals 获取叶子节点上的所有非终端符号
func (dt *DerivationTree) GetNonTerminals() []*Node {
	nonTermi := make([]*Node, 0)
	dt.traversal(dt.root, func(node *Node) {
		if node == nil {
			return
		}
		// 如果当前节点有孩子节点，则不是叶子结点
		if node.Children != nil {
			return
		}
		if IsNonTerminals(node.GetName()) {
			nonTermi = append(nonTermi, node)
		}
	})
	return nonTermi
}
func (dt *DerivationTree) GetNode(symbol string) *Node {
	var n *Node
	dt.traversal(dt.root, func(node *Node) {
		if node.GetName() == symbol {
			n = node
		}
	})
	return n
}
func (dt *DerivationTree) GetRoot() *Node {
	return dt.root
}
func (dt *DerivationTree) RegisterHandler(name string, f Handler) {
	if _, ok := dt.handlers[name]; !ok {
		dt.handlers[name] = make([]Handler, 0)
	}
	dt.handlers[name] = append(dt.handlers[name], f)
}

func (dt *DerivationTree) registerSystemHandler(f Handler) {
	dt.systemHandlers = append(dt.systemHandlers, f)
}

func (dt *DerivationTree) ExpandNode(ctx *DerivationContext) (*DerivationContext, bool) {
	ctx.Next()(ctx, dt)
	return ctx, ctx.finish
}

func (dt *DerivationTree) GetLeafNodes() ([]*Node, string) {
	arr := make([]*Node, 0)
	// force to use DFS
	res := ""
	DFS(dt.root, func(node *Node) {
		if node == nil {
			return
		}
		if node.Children == nil {
			arr = append(arr, node)
			res += node.GetName()
		}
	})
	return arr, res
}
func (dt *DerivationTree) addNode(graph *cgraph.Graph, parent *cgraph.Node, node *Node) {
	if node == nil {
		return
	}
	n, err := graph.CreateNode(node.name)
	if err != nil {
		panic(err)
	}
	if parent != nil {
		_, err := graph.CreateEdge("", parent, n)
		if err != nil {
			panic(err)
		}
	}

	for _, child := range node.Children {
		dt.addNode(graph, n, child)
	}
}

func (dt *DerivationTree) Visualize(filename string) {
	g := graphviz.New()
	graph, _ := g.Graph()

	dt.addNode(graph, nil, dt.root)

	// Render graph to file
	err := g.RenderFilename(graph, graphviz.PNG, filename)
	if err != nil {
		panic(err)
	}
}

func (dt *DerivationTree) ResetCoverage(ctx *DerivationContext) {
	ctx.coverage = make(map[string]int)
	for symbol, expansions := range dt.gram.G {
		for _, expansion := range expansions {
			ctx.coverage[fmt.Sprintf("%s->%s", symbol, expansion.GetName())] = 0
		}
	}
}

func (dt *DerivationTree) Configure(ctx *DerivationContext) *DerivationContext {
	dt.ResetCoverage(ctx)
	if hds, ok := dt.handlers["default"]; ok {
		for _, hd := range hds {
			ctx.handlers = append(ctx.handlers, hd)
		}
	}
	ctx.handlers = append(ctx.handlers, dt.systemHandlers...)
	return ctx
}
func (ctx *DerivationContext) GetCoverage() (map[string]int, float64) {
	covered := 0.0
	for _, cnt := range ctx.coverage {
		if cnt != 0 {
			covered++
		}
	}
	return ctx.coverage, covered / float64(len(ctx.coverage))
}

func NewDerivationTree(gram Grammar, startSymbol string, traversal TraversalAlgo, expandAlgo ExpansionAlgo) *DerivationTree {
	dt := &DerivationTree{
		gram:        gram,
		startSymbol: startSymbol,
		traversal:   traversal,
		expandAlgo:  Handler(expandAlgo),
		handlers:    make(map[string][]Handler),
	}
	dt.registerSystemHandler(countCoverage)
	dt.registerSystemHandler(checkFinish)
	dt.registerSystemHandler(dt.expandAlgo)
	dt.Construct()

	return dt
}

func NewDerivationContext(ctx context.Context) *DerivationContext {
	return &DerivationContext{
		Context:  ctx,
		coverage: make(map[string]int),
	}
}
