package schemas

import (
	"errors"
	"github.com/lucasjones/reggen"
	"log/slog"
	"math/rand"
	"regexp"
)

const (
	CatHandlerName      = "cat_handler"
	OrHandlerName       = "or_handler"
	IDHandlerName       = "id_handler"
	TerminalHandlerName = "terminal_handler"
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
	if len(*cur.GetSymbols()) == 0 {
		chain.Next(ctx, cb)
		return
	}
	// todo 此处的遍历树的手段会导致先生成后加入的节点，需要修改
	ctx.SymbolStack.Pop()
	ctx.SymbolStack.Push((*cur.GetSymbols())[0])
	chain.Next(ctx, cb)
	for i := 1; i < len(*cur.GetSymbols()); i++ {
		ctx.SymbolStack.Push((*cur.GetSymbols())[i])
	}
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
	if len(*cur.GetSymbols()) == 0 {
		chain.Next(ctx, cb)
		return
	}
	ctx.SymbolStack.Pop()
	idx := rand.Int() % len(*cur.GetSymbols())
	ctx.SymbolStack.Push((*cur.GetSymbols())[idx])
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

	if len(*cur.GetSymbols()) != 0 {
		slog.Error("Pattern mismatched[Identifier]")
		return
	}
	if _, ok := ctx.grammarMap[cur.content]; !ok {
		slog.Error("The identifier does not Existed")
		panic("fuck")
	}
	ctx.SymbolStack.Push(ctx.grammarMap[cur.content])
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

func (h *TermHandler) Handle(chain *Chain, ctx *Context, cb ResponseCallBack) {
	cur := ctx.SymbolStack.Top()
	ctx.SymbolStack.Pop()

	if len(*cur.GetSymbols()) != 0 {
		slog.Error("Pattern mismatched[Terminal]")
		return
	}
	result, err := reggen.Generate(cur.content, 1)
	if err != nil {
		slog.Error(err.Error())
		return
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
