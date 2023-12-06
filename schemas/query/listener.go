package query

import (
	"github.com/CUHK-SE-Group/generic-generator/schemas/query/pathQuery"
	"github.com/antlr4-go/antlr/v4"
)

const (
	direct = 1
	skip   = 2
)

type pathListener struct {
	*pathQuery.BasePathQueryListener
	queries []string
}

type path struct {
	queries []pathNode
	idx     int
}
type pathNode struct {
	content string
	t       int // type
}

func (p *path) Fetch() pathNode {
	res := p.queries[p.idx]
	p.idx++
	return res
}
func (p *path) Finish() bool {
	return p.idx == len(p.queries)
}

// EnterNode is called when production Node is entered.
func (s *pathListener) EnterNode(ctx *pathQuery.NodeContext) {
	s.queries = append(s.queries, ctx.GetText())
}

// EnterAny is called when production Any is entered.
func (s *pathListener) EnterAny(ctx *pathQuery.AnyContext) {
	s.queries = append(s.queries, ctx.GetText())
}

// EnterRootNode is called when production rootNode is entered.
func (s *pathListener) EnterRootNode(ctx *pathQuery.RootNodeContext) {
	//s.queries = append(s.queries, "/")
}

// EnterChild is called when production Child is entered.
func (s *pathListener) EnterChild(ctx *pathQuery.ChildContext) {
	s.queries = append(s.queries, "/")
}

// EnterAll is called when production All is entered.
func (s *pathListener) EnterAll(ctx *pathQuery.AllContext) {
	s.queries = append(s.queries, "//")
}

func Parse(input string) []string {
	is := antlr.NewInputStream(input)
	lexer := pathQuery.NewPathQueryLexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	parser := pathQuery.NewPathQueryParser(stream)
	listener := &pathListener{}
	antlr.ParseTreeWalkerDefault.Walk(listener, parser.Query())
	return listener.queries
}
