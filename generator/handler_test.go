package Generator

import (
	"fmt"
	"github.com/CUHK-SE-Group/ebnf-based-generator/parser"
	"github.com/CUHK-SE-Group/ebnf-based-generator/schemas"
	"testing"
)

func TestHandler(t *testing.T) {
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
	for i := 0; i < 1000; i++ { // todo 此处会无限循环，并没有任何输出，因为CatHandler只hook了 GrammarProduction 和 GrammarExpr 两种type
		chain.Next(ctx, func(result *schemas.Result) {
			ctx = result.GetCtx()
			ctx.HandlerIndex = 0
			fmt.Println(ctx.Result)
		})
	}

}
