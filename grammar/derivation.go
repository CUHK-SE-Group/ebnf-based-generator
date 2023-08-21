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
	// The keys of this map are strings in the format "symbol->expansion", representing a specific grammar rule, and the values are integers representing the number of times that rule has been used.
	coverage      map[string]int // track how many times each grammar rule (from symbol to expansion) has been used during the derivation process.
	handlers      []Handler      // handlers to be executed.
	handlerIndex  int
	preExpansions []ExpansionPair // previous expansions records. May be used for coverage analysis
	finish        bool
}

// Returns the next handler to be executed in the given context
func (ctx *DerivationContext) Next() Handler {
	if ctx.handlerIndex <= len(ctx.handlers) {
		ctx.handlerIndex++
	}
	return ctx.handlers[ctx.handlerIndex-1]
}

type DerivationTree struct {
	gram        Grammar
	startSymbol string
	traversal   TraversalAlgo // The algorithm used for tree traversal.
	expandAlgo  Handler       // The algorithm used for expanding nodes.
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

// Initializes the tree based on the grammar:
// Root with start symbol as its name, and its children are the expansions of the start symbol.
func (dt *DerivationTree) Construct() {
	for key, value := range dt.gram.G {
		if key == dt.startSymbol {
			dt.root = &Node{ExpansionTuple: ExpansionTuple{name: key}}
			dt.root.Children = make([]*Node, 0)
			for _, v := range value {
				dt.root.Children = append(dt.root.Children, &Node{ExpansionTuple: v})
			}
			break
		}
	}
}

// Returns all non-terminal nodes from leaves.
func (dt *DerivationTree) GetNonTerminals() []*Node {
	nonTermi := make([]*Node, 0)
	dt.traversal(dt.root, func(node *Node) {
		if node == nil {
			return
		}
		// has children -> not leaf nodes
		if node.Children != nil {
			return
		}
		if IsNonTerminals(node.GetName()) {
			nonTermi = append(nonTermi, node)
		}
	})
	return nonTermi
}

// Returns a node with a specific symbol.
func (dt *DerivationTree) GetNode(symbol string) *Node {
	var n *Node
	dt.traversal(dt.root, func(node *Node) {
		if node.GetName() == symbol {
			n = node
		}
	})
	return n
}

// Returns the root node of the DerivationTree
func (dt *DerivationTree) GetRoot() *Node {
	return dt.root
}

// Registers a handler for a specific symbol.
// Simply append the handler to the list of handlers.
func (dt *DerivationTree) RegisterHandler(name string, f Handler) {
	if _, ok := dt.handlers[name]; !ok {
		dt.handlers[name] = make([]Handler, 0)
	}
	dt.handlers[name] = append(dt.handlers[name], f)
}

// Registers a system-level handler
func (dt *DerivationTree) registerSystemHandler(f Handler) {
	dt.systemHandlers = append(dt.systemHandlers, f)
}

// Expands a node using the next registered handlers.
func (dt *DerivationTree) ExpandNode(ctx *DerivationContext) (*DerivationContext, bool) {
	ctx.Next()(ctx, dt)
	return ctx, ctx.finish
}

// Returns all leaf nodes and their names concatenated.
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

// Helper method for visualization. Adds a node to the graph.
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

// Visualizes the tree and saves it as an image, given the filename.
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

// Resets the coverage map in the context.
func (dt *DerivationTree) ResetCoverage(ctx *DerivationContext) {
	ctx.coverage = make(map[string]int)
	for symbol, expansions := range dt.gram.G {
		for _, expansion := range expansions {
			ctx.coverage[fmt.Sprintf("%s->%s", symbol, expansion.GetName())] = 0
		}
	}
}

// Configures the context, using default & system handlers.
func (dt *DerivationTree) Configure(ctx *DerivationContext) *DerivationContext {
	dt.ResetCoverage(ctx)
	// Add default handlers
	if hds, ok := dt.handlers["default"]; ok {
		ctx.handlers = append(ctx.handlers, hds...)
	}
	// Add system handlers
	ctx.handlers = append(ctx.handlers, dt.systemHandlers...)
	return ctx
}

// for a context, returns the coverage map and the coverage value calculated (0.0-1.0)
func (ctx *DerivationContext) GetCoverage() (map[string]int, float64) {
	covered := 0.0
	for _, cnt := range ctx.coverage {
		if cnt != 0 {
			covered++
		}
	}
	return ctx.coverage, covered / float64(len(ctx.coverage))
}

// Creates a new DerivationTree based on the provided grammar, start symbol, traversal algorithm, and expansion algorithm.
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

// Creates a new DerivationContext with empty coverage based on the provided context.
func NewDerivationContext(ctx context.Context) *DerivationContext {
	return &DerivationContext{
		Context:  ctx,
		coverage: make(map[string]int),
	}
}
