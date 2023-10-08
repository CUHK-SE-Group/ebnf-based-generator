package schemas

import (
	"fmt"
	"github.com/goccy/go-graphviz/cgraph"
	"github.com/golang/glog"
	"strings"
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

type Grammar struct {
	id      string
	gtype   GrammarType
	symbols *[]*Grammar

	content string
	root    *Grammar
	config  *Config
}

func NewGrammar(t GrammarType, id string, content string, conf *Config) *Grammar {
	newG := &Grammar{
		id:      id,
		gtype:   t,
		symbols: &[]*Grammar{},
		content: content,
		config:  conf,
	}
	return newG.SetRoot(newG)
}

func (g *Grammar) GetID() string {
	return g.id
}
func (g *Grammar) GetType() GrammarType {
	return g.gtype
}

func (g *Grammar) SetRoot(root *Grammar) *Grammar {
	g.root = root
	return g
}

func (g *Grammar) GetRoot() *Grammar {
	return g.root
}

func (g *Grammar) AddSymbol(new *Grammar) int {
	if new.gtype != GrammarProduction {
		new.root = g.root
	}
	*g.symbols = append(*g.symbols, new)
	return len(*g.symbols) - 1
}

func (g *Grammar) GetSymbols() *[]*Grammar {
	return g.symbols
}

func (g *Grammar) GetSymbol(idx int) *Grammar {
	if idx < len(*g.symbols) {
		return (*g.symbols)[idx]
	}
	return nil
}

func (g *Grammar) SetSymbol(idx int, sym *Grammar) *Grammar {
	if idx < len(*g.symbols) {
		(*g.symbols)[idx] = sym
	} else {
		glog.Fatal("idx out of range")
	}
	return g
}

func (g *Grammar) GetContent() string {
	return g.content
}

//func (g *Grammar) Visualize(filename string, expandSub bool) {
//	gh := graphviz.New()
//	graph, _ := gh.Graph()
//	visNode = make(map[string]*cgraph.Node)
//	g.addNodeToGraph(graph, nil, g.GetContent(), expandSub)
//
//	err := gh.RenderFilename(graph, graphviz.PNG, filename)
//	if err != nil {
//		panic(err)
//	}
//}

func (g *Grammar) addNodeToGraph(graph *cgraph.Graph, parent *cgraph.Node, label string, expandSub bool) {
	var n *cgraph.Node
	visited := false
	fmt.Println(g.id)
	if strings.Contains(g.id, "EOF") {
		return
	}
	if node, ok := visNode[g.id]; ok {
		n = node
		visited = true
		return
	} else {
		node, err := graph.CreateNode(g.GetID())
		if err != nil {
			panic(err)
		}
		n = node
		visNode[g.id] = n
	}

	if parent != nil {
		edge, err := graph.CreateEdge(label, parent, n)
		if err != nil {
			panic(err)
		}
		edge.SetLabel(label)
	}
	if (parent != nil && g.gtype == GrammarProduction && !expandSub) || visited {
		return
	}
	for _, child := range *g.symbols {
		child.addNodeToGraph(graph, n, g.GetID(), expandSub)
	}
}
