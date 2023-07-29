package grammar

import (
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
	})
	return finish
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
