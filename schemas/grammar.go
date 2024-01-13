package schemas

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/CUHK-SE-Group/generic-generator/graph"
	A "github.com/IBM/fp-go/array"
	"github.com/lucasjones/reggen"
)

type GrammarType int

// 带yes标记的symbol要指定生成策略
const (
	GrammarProduction GrammarType = 1 << iota
	GrammarOR                     // yes
	GrammarCatenate
	GrammarOptional // yes
	GrammarREP      // yes
	GrammarPLUS     // yes
	GrammarEXT      // yes
	GrammarSUB      // yes
	GrammarID
	GrammarTerminal
	GrammarChoice
)
const (
	Prop     = "Property"
	Index    = "index"
	Visited  = "visited_"
	Distance = "distance"
)

var typeStrRep = map[GrammarType]string{
	GrammarProduction: "GrammarProduction",
	GrammarOR:         "GrammarOR",
	GrammarCatenate:   "GrammarCatenate",
	GrammarOptional:   "GrammarOptional",
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
	Type               GrammarType
	Root               *Node
	Gram               *Grammar
	Content            string
	DistanceToTerminal int
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

type Derivation struct {
	*Grammar
	EdgeHistory []string
	SymbolCnt   map[string]int // 为了给语法图上的Node做标记
}

func (d *Derivation) getNodeID(id string) string {
	return id
	// return fmt.Sprintf("%s#%d", id, d.SymbolCnt[id])
}

func (d *Derivation) AddEdge(from, to *Node) {
	newfrom := from.Clone(d.Grammar)
	newto := to.Clone(d.Grammar)

	newfrom.SetID(d.getNodeID(from.GetID())) // update the `from` node to the current existing node
	newto.SetID(d.getNodeID(to.GetID()))

	newfrom.AddSymbol(newto)
	d.EdgeHistory = append(d.EdgeHistory, GetEdgeID(newfrom.GetID(), newto.GetID()))
	d.SymbolCnt[to.GetID()]++ // denote the `to` node to a new node

}
func isTermPreserve(content string) bool {
	return (content[0] == content[len(content)-1]) && ((content[0] == '\'') || content[0] == '"')
}

func (d *Derivation) GetResult(custom func(content string) string) string {
	d.SymbolCnt = make(map[string]int)
	res := ""
	for _, e := range d.EdgeHistory {
		from, to := ExtractEdgeID(e)
		if from == to {
			cur := d.GetNode(to)
			content := cur.GetContent()
			if isTermPreserve(content) {
				tmp := strings.Trim(content, "'\"")
				switch content[0] {
				case '"':
					var err error
					content, err = reggen.Generate(tmp, 10)
					if err != nil {
						panic(err)
					}
				case '\'':
					content = tmp
				default:
					panic("error in generating terminal")
				}
			}
			if custom != nil {
				content = custom(content)
			}

			res += content
			d.SymbolCnt[to]++
			continue
		}
	}
	return res
}

type Grammar struct {
	internal  graph.Graph[string, Property]
	startSym  string
	terminals []string // cache
}

func NewGrammar(start ...string) *Grammar {
	var st string
	if len(start) != 0 {
		st = start[0]
	}
	newG := &Grammar{
		internal: graph.NewGraph[string, Property](),
		startSym: st,
	}
	return newG
}
func (g *Grammar) GetInternal() graph.Graph[string, Property] {
	return g.internal
}

func (g *Grammar) GetNode(id string) *Node {
	if inter := g.internal.GetVertexById(id); inter != nil {
		return &Node{internal: inter}
	}
	return nil
}

func (g *Grammar) GetEdge(id string) (*Node, *Node) {
	from, to := ExtractEdgeID(id)
	return g.GetNode(from), g.GetNode(to)
}

func (p *Grammar) MergeProduction() {
	start := p.startSym
	queue := []*Node{p.GetNode(start)}
	visited := make(map[string]any)
	productions := []*Node{p.GetNode(start)}
	for len(queue) != 0 {
		for _, n := range queue[0].GetSymbols() {
			if n.GetType() == GrammarID {
				productions = append(productions, n)
				v := p.GetNode(fmt.Sprintf("%s", n.GetContent()))
				if v != nil {
					n.AddSymbol(v)
					queue = append(queue, v)
				}
			}
			if _, ok := visited[n.GetID()]; !ok {
				queue = append(queue, n)
				visited[n.GetID()] = ""
			}
		}
		queue = queue[1:]
	}
}

// QueryDistance query for the largest and minimum distance to the terminal
func (g *Grammar) QueryDistance(id string) (float64, float64) {
	if g.internal.GetMetadata(Distance) == nil {
		//tmpG := graph.Clone[string, Property](g.internal, graph.NewGraph[string, Property], graph.NewEdge[string, Property], graph.NewVertex[Property])
		//gg := &Grammar{internal: tmpG, startSym: g.startSym}
		g.MergeProduction()
		distances := graph.FloydAlgorithm(g.internal)
		g.internal.SetMetadata(Distance, distances)
	}
	dis := g.internal.GetMetadata(Distance).(map[string]map[string]float64)
	if g.terminals == nil || len(g.terminals) == 0 {
		g.terminals = make([]string, 0)
		for _, v := range g.GetInternal().GetAllVertices() {
			if v.GetProperty(Prop).Type == GrammarTerminal {
				g.terminals = append(g.terminals, v.GetID())
			}
		}
	}
	mi, ma := math.Inf(1), math.Inf(-1)
	for _, v := range g.terminals {
		mi = min(dis[id][v], mi)
		ma = max(dis[id][v], ma)
	}
	return mi, ma
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

func (g *Grammar) BuildShortestNotation() {
	vertices := g.internal.GetAllVertices()
	numVertices := len(vertices)
	sort.Slice(vertices, func(i, j int) bool {
		return vertices[i].GetID() > vertices[j].GetID()
	})
	vertexMap := make(map[string]int)
	for i, vertex := range vertices {
		vertexMap[vertex.GetID()] = i
	}

	// some sufficiently large value
	// indicating not-reachable
	inf := int(1e8)
	distance := make([]int, numVertices)
	for i, vertex := range vertices {
		if vertex.GetProperty(Prop).Type == GrammarTerminal {
			distance[i] = 0
		} else {
			distance[i] = inf
		}
	}

	// Bellman-ford-like process;
	// this should terminate within O(numVertices) iterations,
	// i.e. no more relaxations after that
	round := 0
	for {
		round++
		stop := true
		for index, current := range vertices {
			adjacent := g.internal.GetOutEdges(current)
			pre := distance[index]
			if current.GetProperty(Prop).Type == GrammarTerminal {
				// do nothing
			} else if current.GetProperty(Prop).Type == GrammarOR {
				//if strings.Contains(current.GetProperty(Prop).Content, "1") {
				//	fmt.Println()
				//}
				// 1 + min of {distances of outgoing neighbors}
				best := inf
				for _, e := range adjacent {
					next_index := vertexMap[e.GetTo().GetID()]
					best = min(best, distance[next_index])
				}
				best += 1
				distance[index] = min(distance[index], best)
			} else {
				// 1 + sum of {distances of outgoing neighbors}
				sum := 0
				for _, e := range adjacent {
					next_index := vertexMap[e.GetTo().GetID()]
					sum += distance[next_index]
					// watch out for overflows
					if distance[next_index] >= inf {
						break
					}
				}
				sum += 1
				distance[index] = min(distance[index], sum)
				if distance[index] < 1 {
					fmt.Println("fuck")
				}
			}
			if pre != distance[index] {
				stop = false
			}
		}
		if stop {
			fmt.Println(round)
			break
		}
	}

	for index, v := range vertices {
		if v.GetProperty(Prop).Type == GrammarTerminal {
			continue
		}
		prop := v.GetProperty(Prop)
		prop.DistanceToTerminal = distance[index]
		v.SetProperty(Prop, prop)
	}
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

func (g *Node) Clone(belongto *Grammar) *Node {
	newInternal := graph.CloneVertex(g.internal, graph.NewVertex[Property])
	if belongto != nil {
		p := newInternal.GetProperty(Prop)
		p.Gram = belongto
		newInternal.SetProperty(Prop, p)
	}
	return &Node{internal: newInternal}
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
func (g *Node) SetID(id string) {
	g.internal.SetID(id)
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
	e := newEdge(GetEdgeID(g.GetID(), new.GetID()), g, new)
	g.GetGrammar().internal.AddEdge(e)
	return len(g.GetGrammar().internal.GetOutEdges(g.internal)) - 1
}
func getNumber(id string) int {
	ids := strings.Split(id, "#")
	if len(ids) != 2 {
		panic("the id format should be xxx#yyy")
	}
	num1, err := strconv.Atoi(ids[1])
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
func (g *Node) SetContent(content string) {
	p := g.internal.GetProperty(Prop)
	p.Content = content
	g.internal.SetProperty(Prop, p)
}
func (g *Node) GetDistance() int {
	return g.internal.GetProperty(Prop).DistanceToTerminal
}
func GetEdgeID(father string, child string) string {
	return fmt.Sprintf("%s,%s", father, child)
}

func ExtractEdgeID(id string) (string, string) {
	res := strings.Split(id, ",")
	if len(res) != 2 {
		panic("error in id")
	}
	return res[0], res[1]
}
