package grammar

import (
	"fmt"
	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
)

type DerivationTree struct {
	gram        Grammar
	startSymbol string
	traversal   TraversalAlgo
	expandAlgo  ExpansionAlgo
	root        *Node
}
type Node struct {
	ExpansionTuple
	Children []*Node
}

func (dt *DerivationTree) Construct() {
	for key, value := range dt.gram {
		if key == dt.startSymbol {
			dt.root = &Node{ExpansionTuple: ExpansionTuple{name: key}}
			dt.root.Children = make([]*Node, 0)
			for _, v := range value {
				dt.root.Children = append(dt.root.Children, &Node{ExpansionTuple: v})
			}
		}
	}
}
func (dt *DerivationTree) GetNonTerminals() []*Node {
	nonTermi := make([]*Node, 0)
	dt.traversal(dt.root, func(node *Node) {
		if node == nil {
			return
		}
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
func (dt *DerivationTree) ExpandNode() bool {
	dt.expandAlgo(dt)
	finish := true
	dt.traversal(dt.root, func(node *Node) {
		if node == nil {
			return
		}
		if node.Children == nil && IsNonTerminals(node.GetName()) {
			finish = false
		}
		fmt.Println(node.GetName())
	})
	return finish
}
func (dt *DerivationTree) addNode(graph *cgraph.Graph, parent *cgraph.Node, node *Node) {
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
