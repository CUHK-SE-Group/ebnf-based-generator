package schemas

import (
	"errors"
	"fmt"
	"github.com/CUHK-SE-Group/generic-generator/schemas/query"
	"log/slog"
	"math"
	"math/rand"
	"regexp"
	"strings"
)

const (
	CatHandlerName      = "cat_handler"
	OrHandlerName       = "or_handler"
	IDHandlerName       = "id_handler"
	BracketHandlerName  = "Bracket_handler"
	ParenHandlerName    = "paren_handler"
	TerminalHandlerName = "terminal_handler"
	SubHandlerName      = "sub_handler"
	BraceHandlerName    = "brace_handler"
	RepHandlerName      = "rep_handler"
	TraceHandlerName    = "trace_handler"
	OptionHandlerName   = "option_handler"
)

var errViolateBuildIn = errors.New("can not replace build-in handler func")
var buildIn = []string{}
var ErrDuplicatedHandler = errors.New("duplicated handler registration")
var funcMap = make(map[GrammarType][]func() Handler)

type Handler interface {
	Handle(*Chain, *Context, ResponseCallBack)
	HookRoute() []regexp.Regexp
	Name() string
	Type() GrammarType
}

type CatHandler struct {
}

func (h *CatHandler) Handle(chain *Chain, ctx *Context, cb ResponseCallBack) {
	cur := ctx.SymbolStack.Top()
	if len(cur.GetSymbols()) == 0 {
		chain.Next(ctx, cb)
		return
	}
	ctx.SymbolStack.Pop()
	for i := len(cur.GetSymbols()) - 1; i >= 0; i-- {
		ctx.SymbolStack.Push((cur.GetSymbols())[i])
	}
	chain.Next(ctx, cb)
}

func (h *CatHandler) HookRoute() []regexp.Regexp {
	return make([]regexp.Regexp, 0)
}

func (h *CatHandler) Name() string {
	return CatHandlerName
}

func (h *CatHandler) Type() GrammarType {
	return GrammarProduction | GrammarCatenate
}

type OrHandler struct {
}

func (h *OrHandler) Handle(chain *Chain, ctx *Context, cb ResponseCallBack) {
	cur := ctx.SymbolStack.Top()
	if len(cur.GetSymbols()) == 0 {
		chain.Next(ctx, cb)
		return
	}
	ctx.SymbolStack.Pop()
	idx := rand.Int() % len(cur.GetSymbols())
	ctx.SymbolStack.Push((cur.GetSymbols())[idx])
	ctx.VisitedEdge[GetEdgeID(cur.GetID(), (cur.GetSymbols())[idx].GetID())]++
	chain.Next(ctx, cb)
}

func (h *OrHandler) HookRoute() []regexp.Regexp {
	return make([]regexp.Regexp, 0)
}

func (h *OrHandler) Name() string {
	return OrHandlerName
}

func (h *OrHandler) Type() GrammarType {
	return GrammarOR
}

type IDHandler struct {
}

func (h *IDHandler) Handle(chain *Chain, ctx *Context, cb ResponseCallBack) {
	cur := ctx.SymbolStack.Top()
	ctx.SymbolStack.Pop()

	//if len(cur.GetSymbols()) != 0 {
	//	slog.Error("Pattern mismatched[Identifier]")
	//	//return
	//}
	node := ctx.Grammar.GetNode(cur.GetContent())
	if node.internal == nil {
		slog.Error("The identifier does not Existed", "id", cur.GetContent())
		ctx.Result += cur.GetContent()
		return // omit error
	}
	ctx.SymbolStack.Push(node)
	chain.Next(ctx, cb)
}

func (h *IDHandler) HookRoute() []regexp.Regexp {
	return make([]regexp.Regexp, 0)
}

func (h *IDHandler) Name() string {
	return IDHandlerName
}

func (h *IDHandler) Type() GrammarType {
	return GrammarID
}

type RepHandler struct {
}

func (r *RepHandler) Handle(chain *Chain, ctx *Context, cb ResponseCallBack) {
	cur := ctx.SymbolStack.Top()
	ctx.SymbolStack.Pop()
	for i := 0; i < rand.Intn(1); i++ {
		ctx.SymbolStack.Push(cur.GetSymbols()...)
	}
	chain.Next(ctx, cb)
}

func (r *RepHandler) HookRoute() []regexp.Regexp {
	return make([]regexp.Regexp, 0)
}

func (r *RepHandler) Name() string {
	return RepHandlerName
}

func (r *RepHandler) Type() GrammarType {
	return GrammarREP
}

type TermHandler struct {
}

func (h *TermHandler) isTermPreserve(g *Node) bool {
	content := g.GetContent()
	return (content[0] == content[len(content)-1]) && (content[0] == '\'')
}

func (h *TermHandler) stripQuote(content string) string {
	if content[0] == content[len(content)-1] {
		if (content[0] == '\'') || (content[0] == '"') {
			return content[1 : len(content)-1]
		}
	}
	return content
}

func (h *TermHandler) Handle(chain *Chain, ctx *Context, cb ResponseCallBack) {
	cur := ctx.SymbolStack.Top()
	for key, fn := range FirstPlaceNode {
		if !strings.Contains(cur.GetID(), ctx.SymbolStack.ProductionTrace[len(ctx.SymbolStack.ProductionTrace)-1]) {
			continue
		}
		if query.MatchPattern(ctx.SymbolStack.ProductionTrace, key) {
			ctx = fn(ctx)
			chain.Next(ctx, cb)
			return
		}
	}
	for key, fn := range SecondPlaceNode {
		if !strings.Contains(cur.GetID(), ctx.SymbolStack.ProductionTrace[len(ctx.SymbolStack.ProductionTrace)-1]) {
			continue
		}
		if query.MatchPattern(ctx.SymbolStack.ProductionTrace, key) {
			var err error
			ctx, err = fn.Func(ctx)
			if err != nil {
				if errors.Is(err, ErrSymbolNotFound) {
					continue
				}
				panic(err)
			}
			chain.Next(ctx, cb)
			return
		}
	}
	ctx.SymbolStack.Pop()

	if len(cur.GetSymbols()) != 0 {
		slog.Error("Pattern mismatched[Terminal]")
		return
	}
	result := ""
	regex := h.stripQuote(cur.GetContent())
	if h.isTermPreserve(cur) {
		result = regex
	} else {
		result = regex
	}
	ctx.Result += result
	chain.Next(ctx, cb)
}

func (h *TermHandler) HookRoute() []regexp.Regexp {
	return make([]regexp.Regexp, 0)
}

func (h *TermHandler) Name() string {
	return TerminalHandlerName
}

func (h *TermHandler) Type() GrammarType {
	return GrammarTerminal
}

type BracketHandler struct {
}

func (h *BracketHandler) Handle(chain *Chain, ctx *Context, cb ResponseCallBack) {
	//cur := ctx.SymbolStack.Top()
	ctx.SymbolStack.Pop()

	//if len(cur.GetSymbols()) == 0 {
	//	slog.Error("Pattern mismatched[Identifier]")
	//	return
	//}
	//for i := len(cur.GetSymbols()) - 1; i >= 0; i-- {
	//	ctx.SymbolStack.Push(cur.GetSymbol(i))
	//}
	chain.Next(ctx, cb)
}

func (h *BracketHandler) HookRoute() []regexp.Regexp {
	return make([]regexp.Regexp, 0)
}

func (h *BracketHandler) Name() string {
	return BracketHandlerName
}

func (h *BracketHandler) Type() GrammarType {
	return GrammarOptional
}

type ParenHandler struct {
}

func (h *ParenHandler) Handle(chain *Chain, ctx *Context, cb ResponseCallBack) {
	cur := ctx.SymbolStack.Top()
	ctx.SymbolStack.Pop()
	for i := len(cur.GetSymbols()) - 1; i >= 0; i-- {
		ctx.SymbolStack.Push(cur.GetSymbol(i))
	}
	chain.Next(ctx, cb)
}

func (h *ParenHandler) HookRoute() []regexp.Regexp {
	return make([]regexp.Regexp, 0)
}

func (h *ParenHandler) Name() string {
	return ParenHandlerName
}

func (h *ParenHandler) Type() GrammarType {
	return GrammarPAREN
}

type BraceHandler struct {
}

func (h *BraceHandler) Handle(chain *Chain, ctx *Context, cb ResponseCallBack) {
	cur := ctx.SymbolStack.Top()
	ctx.SymbolStack.Pop()

	if len(cur.GetSymbols()) == 0 {
		slog.Error("Pattern mismatched[Identifier]")
		return
	}

	for i := len(cur.GetSymbols()) - 1; i >= 0; i-- {
		ctx.SymbolStack.Push(cur.GetSymbol(i))
	}
	chain.Next(ctx, cb)
}

func (h *BraceHandler) HookRoute() []regexp.Regexp {
	return make([]regexp.Regexp, 0)
}

func (h *BraceHandler) Name() string {
	return BraceHandlerName
}

func (h *BraceHandler) Type() GrammarType {
	return GrammarBRACE
}

type SubHandler struct {
}

func (h *SubHandler) Handle(chain *Chain, ctx *Context, cb ResponseCallBack) {
	ctx.SymbolStack.Pop()
	chain.Next(ctx, cb)
}

func (h *SubHandler) HookRoute() []regexp.Regexp {
	return make([]regexp.Regexp, 0)
}

func (h *SubHandler) Name() string {
	return SubHandlerName
}

func (h *SubHandler) Type() GrammarType {
	return GrammarSUB
}

type TraceHandler struct {
}

func (h *TraceHandler) Handle(chain *Chain, ctx *Context, cb ResponseCallBack) {
	chain.Next(ctx, cb)
}

func (h *TraceHandler) HookRoute() []regexp.Regexp {
	return make([]regexp.Regexp, 0)
}

func (h *TraceHandler) Name() string {
	return TraceHandlerName
}

func (h *TraceHandler) Type() GrammarType {
	return math.MaxInt
}

type OptionHandler struct {
}

func (h *OptionHandler) Handle(chain *Chain, ctx *Context, cb ResponseCallBack) {
	cur := ctx.SymbolStack.Top()
	fmt.Println("dealing with", cur.GetContent())
	ctx.SymbolStack.Pop()
	chain.Next(ctx, cb)
}

func (h *OptionHandler) HookRoute() []regexp.Regexp {
	return make([]regexp.Regexp, 0)
}

func (h *OptionHandler) Name() string {
	return OptionHandlerName
}

func (h *OptionHandler) Type() GrammarType {
	return GrammarChoice
}
