package schemas

import (
	"errors"
	"github.com/golang/glog"
	"github.com/lucasjones/reggen"
	"math/rand"
	"regexp"
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
	return GrammarProduction | GrammarTerm
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
	chain.Next(ctx, cb)
}

func (h *OrHandler) HookRoute() []regexp.Regexp {
	return make([]regexp.Regexp, 0)
}

func (h *OrHandler) Name() string {
	return OrHandlerName
}

func (h *OrHandler) Type() GrammarType {
	return GrammarExpr
}

type IDHandler struct {
}

func (h *IDHandler) Handle(chain *Chain, ctx *Context, cb ResponseCallBack) {
	cur := ctx.SymbolStack.Top()
	ctx.SymbolStack.Pop()

	if len(cur.GetSymbols()) != 0 {
		glog.Error("Pattern mismatched[Identifier]")
		return
	}
	node := ctx.grammarMap.GetNode(cur.GetContent())
	if node == nil {
		glog.Errorf("The identifier [%v] does not Existed", cur.GetContent())
		panic("The identifier does not Existed")
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
	ctx.SymbolStack.Pop()

	if len(cur.GetSymbols()) != 0 {
		glog.Error("Pattern mismatched[Terminal]")
		return
	}
	result := ""
	regex := h.stripQuote(cur.GetContent())
	if h.isTermPreserve(cur) {
		result = regex
	} else {
		var err error
		result, err = reggen.Generate(regex, 1)
		if err != nil {
			glog.Error(err.Error())
			return
		}

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
	cur := ctx.SymbolStack.Top()
	ctx.SymbolStack.Pop()

	if len(cur.GetSymbols()) == 0 {
		glog.Error("Pattern mismatched[Identifier]")
		return
	}
	for i := len(cur.GetSymbols()) - 1; i >= 0; i-- {
		ctx.SymbolStack.Push(cur.GetSymbol(i))
	}
	chain.Next(ctx, cb)
}

func (h *BracketHandler) HookRoute() []regexp.Regexp {
	return make([]regexp.Regexp, 0)
}

func (h *BracketHandler) Name() string {
	return BracketHandlerName
}

func (h *BracketHandler) Type() GrammarType {
	return GrammarBRACKET
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
		glog.Error("Pattern mismatched[Identifier]")
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
	//_ := ctx.SymbolStack.Top()
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
