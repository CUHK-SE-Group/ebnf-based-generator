package Generator

import (
	"context"
	"fmt"
	"github.com/CUHK-SE-Group/generic-generator/parser"
	"github.com/CUHK-SE-Group/generic-generator/schemas"
	"math/rand"
	"regexp"
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
	ctx, err := schemas.NewContext(g, "program", context.Background(), nil)
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

	trace := ctx.SymbolStack.GetStack()
	ctx.SymbolStack.Pop()
	candidates := make([]int, 0)
	sym := cur.GetSymbols()
	for i, v := range sym {
		flag := 0
		for _, his := range trace {
			if v.GetID() == his.GetID() {
				flag++
			}
		}
		if flag == 0 {
			candidates = append(candidates, i)
		}
	}
	idx := rand.Intn(len(candidates))
	ctx.SymbolStack.Push(sym[candidates[idx]])
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
	g, err := parser.Parse("./testdata/complete/tinyc.ebnf", "")
	if err != nil {
		panic(err)
	}
	chain, err := schemas.CreateChain("test", &schemas.CatHandler{}, &schemas.IDHandler{}, &schemas.TermHandler{}, &WeightedHandler{}, &schemas.OrHandler{})
	if err != nil {
		panic(err)
	}
	ctx, err := schemas.NewContext(g, "program", context.Background(), nil)
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
