package schemas

import (
	"errors"
	"fmt"
	"regexp"
)

const (
	CatHandlerName = "cat_handler"
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
	children := ctx.SymbolStack.Top().GetSymbols()
	if len(*children) == 0 {
		ctx.SymbolStack.Pop()
		chain.Next(ctx, cb)
		return
	}
	// todo 此处的遍历树的手段会导致先生成后加入的节点，需要修改
	ctx.SymbolStack.Pop()
	for _, v := range *children {
		fmt.Println(v.content)
		ctx.SymbolStack.Push(v)
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
	return GrammarProduction | GrammarExpr
}
