package Generator

import (
	"context"
	"fmt"
	"github.com/CUHK-SE-Group/generic-generator/parser"
	"github.com/CUHK-SE-Group/generic-generator/schemas"
	"testing"
)

func TestDefaultHandlerCypher(t *testing.T) {
	cons := schemas.MaxLimit
	cons.FirstNode = "Expression"
	cons.SecondNode = "Expression"
	g, err := parser.Parse("./testdata/complete/Cypher.ebnf", "Cypher")
	if err != nil {
		panic(err)
	}
	g.MergeProduction()
	g.BuildShortestNotation()
	consg := schemas.NewConstraintGraph()
	consg.AddBinaryConstraint(cons)
	chain, err := schemas.CreateChain("test", &MonitorHandler{}, &schemas.OptionHandler{}, &MonitorHandler{}, &schemas.CatHandler{}, &schemas.IDHandler{}, &schemas.SubHandler{}, &WeightedHandler{}, &schemas.TermHandler{}, &schemas.RepHandler{}, &schemas.BracketHandler{})
	if err != nil {
		panic(err)
	}
	ctx, err := schemas.NewContext(g, "Cypher", context.Background(), consg, nil)
	if err != nil {
		panic(err)
	}
	for !ctx.GetFinish() {
		chain.Next(ctx, func(result *schemas.Result) {
			ctx = result.GetCtx()
			ctx.HandlerIndex = 0
		})
	}
	fmt.Println(ctx.Result.GetResult(nil))
}

func TestLLVMIRHandler(t *testing.T) {
	g, err := parser.Parse("./testdata/complete/llvmir.ebnf", "module")
	if err != nil {
		panic(err)
	}
	g.MergeProduction()
	g.BuildShortestNotation()
	chain, err := schemas.CreateChain("test", &MonitorHandler{}, &schemas.PlusHandler{}, &schemas.CatHandler{}, &schemas.IDHandler{}, &schemas.TermHandler{}, &WeightedHandler{}, &schemas.OrHandler{}, &schemas.RepHandler{}, &schemas.BracketHandler{})
	if err != nil {
		panic(err)
	}
	ctx, err := schemas.NewContext(g, "module", context.Background(), nil, nil)
	if err != nil {
		panic(err)
	}
	for !ctx.GetFinish() {
		chain.Next(ctx, func(result *schemas.Result) {
			ctx = result.GetCtx()
			ctx.HandlerIndex = 0
		})
	}
	fmt.Printf("%s\n", ctx.Result.GetResult(nil))
}
