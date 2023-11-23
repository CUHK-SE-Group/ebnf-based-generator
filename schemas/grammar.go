package schemas

import (
	"encoding/json"
	"fmt"
	"github.com/CUHK-SE-Group/generic-generator/graph"
	A "github.com/IBM/fp-go/array"
	"log/slog"
	"sort"
	"strconv"
	"strings"
)

type GrammarType int

// 带yes标记的symbol要指定生成策略
const (
	GrammarProduction GrammarType = 1 << iota
	GrammarOR                     // yes
	GrammarCatenate
	GrammarOptional // yes
	GrammarPAREN
	GrammarBRACE // yes
	GrammarREP   // yes
	GrammarPLUS  // yes
	GrammarEXT   // yes
	GrammarSUB   // yes
	GrammarID
	GrammarTerminal
	GrammarChoice
)
const (
	Prop    = "Property"
	Index   = "index"
	Visited = "visited_"
)

var typeStrRep = map[GrammarType]string{
	GrammarProduction: "GrammarProduction",
	GrammarOR:         "GrammarOR",
	GrammarCatenate:   "GrammarCatenate",
	GrammarOptional:   "GrammarOptional",
	GrammarPAREN:      "GrammarPAREN",
	GrammarBRACE:      "GrammarBRACE",
	GrammarREP:        "GrammarREP",
	GrammarPLUS:       "GrammarPLUS",
	GrammarEXT:        "GrammarEXT",
	GrammarSUB:        "GrammarSUB",
	GrammarID:         "GrammarID",
	GrammarTerminal:   "GrammarTerminal",
	GrammarChoice:     "GrammarChoice",
}

func GetGrammarTypeStr(t GrammarType) string {
	return typeStrRep[t]
}

type Property struct {
	Type    GrammarType
	Root    *Node
	Gram    *Grammar
	Content string
}
type TrieTree struct {
	Root  *TrieNode
	Index map[string]*TrieNode
}
type TrieNode struct {
	Children map[string]*TrieNode `json:"children"`
	IsEnd    bool                 `json:"-"`
	Node     *Node
}

func (t *TrieNode) ToJSON() (string, error) {
	type Alias TrieNode

	data := make(map[string]interface{})
	for key, child := range t.Children {
		childJSON, err := child.ToJSON()
		if err != nil {
			return "", err
		}
		data[key] = json.RawMessage(childJSON)
	}

	jsonData, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

func NewTrieNode(node *Node) *TrieNode {
	return &TrieNode{Children: make(map[string]*TrieNode), Node: node}
}

func (t *TrieNode) Insert(n *Node, path []string, index *map[string]*TrieNode) {
	node := t
	for _, element := range path {
		if _, ok := node.Children[element]; !ok {
			node.Children[element] = NewTrieNode(n)
			(*index)[element] = node.Children[element]
		}
		node = node.Children[element]
	}
	node.IsEnd = true
}

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
	if g.internal.GetVertexById(id) == nil {
		return nil
	}
	return &Node{internal: g.internal.GetVertexById(id)}
}
func (g *Grammar) GetIndex(id string) *TrieNode {
	if g.GetInternal().GetMetadata(Index) == nil {
		g.BuildPath(id)
	}
	return g.GetInternal().GetMetadata(Index).(*TrieTree).Index[id]
}
func (g *Grammar) BuildPath(id string) {
	root := g.GetNode(id)
	if root == nil {
		return
	}
	visited := make(map[string]bool)
	trie := NewTrieNode(root)
	index := &(map[string]*TrieNode{})
	var dfs func(node *Node, path []string)
	dfs = func(node *Node, path []string) {
		if node == nil {
			return
		}
		newPath := append(path, node.GetID())

		if node.GetType() == GrammarTerminal {
			trie.Insert(node, newPath, index)
		} else if node.GetType() == GrammarID {
			// do not visit back to the identifier, but we will add a notice
			childNode := g.GetNode(node.GetContent())
			if childNode == nil {
				fmt.Println("empty node")
				return
			}
			if visited[childNode.GetID()] {
				trie.Insert(childNode, append(path, Visited+node.GetContent()), index)
			}
			if !visited[childNode.GetID()] {
				visited[childNode.GetID()] = true
				dfs(childNode, newPath)
			}
		} else {
			for _, child := range node.GetSymbols() {
				dfs(child, newPath)
			}
		}
	}

	dfs(root, []string{})
	updateVisitedNodes(trie)
	tree := &TrieTree{
		Root:  trie,
		Index: *index,
	}

	g.GetInternal().SetMetadata(Index, tree)
}

func updateVisitedNodes(node *TrieNode) bool {
	if node == nil || len(node.Children) == 0 {
		return false
	}

	allChildrenVisited := true
	anyChildVisited := false
	for key, child := range node.Children {
		childVisited := updateVisitedNodes(child)
		if childVisited && strings.HasPrefix(key, Visited) {
			anyChildVisited = true
		}
		allChildrenVisited = allChildrenVisited && childVisited
	}

	// 根据节点的状态决定是否标记为 Visited
	if (node.Node.GetType() == GrammarOR && allChildrenVisited) || (node.Node.GetType() == GrammarCatenate && anyChildVisited) {
		// 更新当前节点为 Visited
		// 这可能需要修改 TrieNode 的结构或添加一个标记
		return true
	}

	return false
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
	if g.internal == nil {
		return 0
	}
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

func (g *Node) SetType(t GrammarType) {
	ori := g.internal.GetProperty(Prop)
	ori.Type = t
	g.internal.SetProperty(Prop, ori)
}

func (g *Node) GetGrammar() *Grammar {
	return g.internal.GetProperty(Prop).Gram
}

func (g *Node) AddSymbol(new *Node) int {
	e := newEdge(fmt.Sprintf("%s,%s", g.GetID(), new.GetID()), g, new)
	g.GetGrammar().internal.AddEdge(e)
	return len(g.GetGrammar().internal.GetOutEdges(g.internal)) - 1
}
func getNumber(id string) int {
	num1, err := strconv.Atoi(strings.Split(id, "#")[1])
	if err != nil {
		slog.Error("strconv atoi", "error", err)
	}
	return num1
}
func (g *Node) GetSymbols() []*Node {
	f := func(edge graph.Edge[string, Property]) *Node {
		return &Node{internal: edge.GetTo()}
	}
	ori := A.Map(f)(g.GetGrammar().internal.GetOutEdges(g.internal))
	sort.Slice(ori, func(i, j int) bool {
		num1 := getNumber(ori[i].GetID())
		num2 := getNumber(ori[j].GetID())
		return num1 < num2
	})
	return ori
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
