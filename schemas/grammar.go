package schemas

import (
	"github.com/dominikbraun/graph"
	"github.com/dominikbraun/graph/draw"
	"github.com/golang/glog"
	"os"
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

func (g *Grammar) Visualize(filename string, expandSub bool) {
	ghash := func(gram *Grammar) string { return gram.content }
	gr := graph.New(ghash, graph.Directed(), graph.Rooted())
	_ = gr.AddVertex(g)

	queue := make([]*Grammar, 0)
	queue = append(queue, g)
	visited := make(map[string]bool, 0)
	for len(queue) != 0 {
		cur := queue[0]
		queue = queue[1:]
		visited[cur.id] = true
		for _, v := range *cur.GetSymbols() {
			_ = gr.AddVertex(v)
			_ = gr.AddEdge(cur.content, v.content)
			if visited[v.id] {
				continue
			}
			queue = append(queue, v)
		}
	}

	file, _ := os.Create(filename)
	_ = draw.DOT(gr, file)
}
