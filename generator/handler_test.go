package Generator

import (
	"fmt"
	"github.com/CUHK-SE-Group/generic-generator/parser"
	"github.com/CUHK-SE-Group/generic-generator/schemas"
	"math/rand"
	"regexp"
	"strings"
	"testing"
)

func TestDefaultHandler(t *testing.T) {
	g, err := parser.Parse("parser/testdata/cypher.ebnf")
	if err != nil {
		panic(err)
	}
	chain, err := schemas.CreateChain("test", &schemas.CatHandler{}, &schemas.OrHandler{}, &schemas.IDHandler{}, &schemas.TermHandler{}, &schemas.BracketHandler{}, &schemas.ParenHandler{}, &schemas.BraceHandler{}, &schemas.SubHandler{})
	if err != nil {
		panic(err)
	}
	ctx, err := schemas.NewContext(g, "Cypher")
	if err != nil {
		panic(err)
	}

	for !ctx.GetFinish() {
		//start := ctx.SymbolStack.Top()
		chain.Next(ctx, func(result *schemas.Result) {
			ctx = result.GetCtx()
			ctx.HandlerIndex = 0
			fmt.Println(ctx.Result)
		})
		//if ctx.SymbolStack.Top() == start {
		//	glog.Errorf("generate failed, type:[%v], content:[%v]", ctx.SymbolStack.Top().GetType(), ctx.SymbolStack.Top().GetContent())
		//	break
		//}
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
	ctx, err := schemas.NewContext(g, "factor")
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
