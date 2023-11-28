package Generator

import (
	"context"
	"fmt"
	"github.com/CUHK-SE-Group/generic-generator/parser"
	"github.com/CUHK-SE-Group/generic-generator/schemas"
	"math/rand"
	"regexp"
	"strings"
	"testing"
)

func TestDefaultHandler(t *testing.T) {
	g, err := parser.Parse("./testdata/complete/tinyc.ebnf", "program")
	if err != nil {
		panic(err)
	}
	chain, err := schemas.CreateChain("test", &schemas.CatHandler{}, &schemas.IDHandler{}, &schemas.SubHandler{}, &schemas.OrHandler{}, &schemas.TermHandler{}, &schemas.RepHandler{}, &schemas.BracketHandler{})
	if err != nil {
		panic(err)
	}
	ctx, err := schemas.NewContext(g, "program", context.Background())
	if err != nil {
		panic(err)
	}

	for !ctx.GetFinish() {
		chain.Next(ctx, func(result *schemas.Result) {
			ctx = result.GetCtx()
			ctx.HandlerIndex = 0
			fmt.Println(ctx.Result)
		})
	}
}

type WeightedHandler struct {
}

func (h *WeightedHandler) Handle(chain *schemas.Chain, ctx *schemas.Context, cb schemas.ResponseCallBack) {
	cur := ctx.SymbolStack.Top()
	if len(cur.GetSymbols()) == 0 {
		chain.Next(ctx, cb)
		return
	}
	var idx int
	ctx.SymbolStack.Pop()

	trace := ctx.SymbolStack.GetTrace()
	if len(trace)-8 >= 0 && trace[len(trace)-8].GetContent() == "'-'" && strings.Contains(cur.GetID(), "factor") {
		idx = 0
		fmt.Println(trace)
	} else {
		idx = rand.Int() % len(cur.GetSymbols())
	}
	ctx.SymCount[cur.GetID()]++
	ctx.SymbolStack.Push((cur.GetSymbols())[idx])
	chain.Next(ctx, cb)
}

func (h *WeightedHandler) HookRoute() []regexp.Regexp {
	return make([]regexp.Regexp, 0)
}

func (h *WeightedHandler) Name() string {
	return "weight"
}

func (h *WeightedHandler) Type() schemas.GrammarType {
	return schemas.GrammarOR
}
func TestWeightedHandler(t *testing.T) {
	g, err := parser.Parse("./testdata/complete/simple.ebnf", "")
	if err != nil {
		panic(err)
	}
	chain, err := schemas.CreateChain("test", &schemas.CatHandler{}, &schemas.IDHandler{}, &schemas.TermHandler{}, &WeightedHandler{}, &schemas.OrHandler{})
	if err != nil {
		panic(err)
	}
	ctx, err := schemas.NewContext(g, "expression", context.Background())
	if err != nil {
		panic(err)
	}
	for !ctx.GetFinish() {
		chain.Next(ctx, func(result *schemas.Result) {
			ctx = result.GetCtx()
			ctx.HandlerIndex = 0
			fmt.Println(ctx.Result)
		})
	}
}
