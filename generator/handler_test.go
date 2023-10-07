package Generator

import (
	"fmt"
	"github.com/CUHK-SE-Group/ebnf-based-generator/parser"
	"github.com/CUHK-SE-Group/ebnf-based-generator/schemas"
	"math/rand"
	"regexp"
	"strings"
	"testing"
)

func TestDefaultHandler(t *testing.T) {
	g, err := parser.Parse("parser/testdata/simple.ebnf")
	if err != nil {
		panic(err)
	}
	chain, err := schemas.CreateChain("test", &schemas.CatHandler{}, &schemas.OrHandler{}, &schemas.IDHandler{}, &schemas.TermHandler{})
	if err != nil {
		panic(err)
	}
	ctx, err := schemas.NewContext(g, "expression")
	if err != nil {
		panic(err)
	}
	for i := 0; i < 1000; i++ {
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
	if len(*cur.GetSymbols()) == 0 {
		chain.Next(ctx, cb)
		return
	}
	var idx int
	ctx.SymbolStack.Pop()
	if strings.Contains(cur.GetID(), "factor") {
		idx = 0

	} else if strings.Contains(cur.GetID(), "expression") {
		idx = 0

	} else if strings.Contains(cur.GetID(), "term") {
		idx = 0

	} else {
		idx = rand.Int() % len(*cur.GetSymbols())
	}
	ctx.SymCount[cur.GetID()]++
	ctx.SymbolStack.Push((*cur.GetSymbols())[idx])
	chain.Next(ctx, cb)
}

func (h *WeightedHandler) HookRoute() []regexp.Regexp {
	return make([]regexp.Regexp, 0)
}

func (h *WeightedHandler) Name() string {
	return "weight"
}

func (h *WeightedHandler) Type() schemas.GrammarType {
	return schemas.GrammarExpr
}
func TestWeightedHandler(t *testing.T) {

	g, err := parser.Parse("parser/testdata/simple.ebnf")
	if err != nil {
		panic(err)
	}
	chain, err := schemas.CreateChain("test", &schemas.CatHandler{}, &schemas.IDHandler{}, &schemas.TermHandler{}, &WeightedHandler{}, &schemas.OrHandler{})
	if err != nil {
		panic(err)
	}
	ctx, err := schemas.NewContext(g, "fake")
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
