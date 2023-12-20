package Generator

import (
	"context"
	"errors"
	"fmt"
	"github.com/CUHK-SE-Group/generic-generator/graph"
	"github.com/CUHK-SE-Group/generic-generator/parser"
	"github.com/CUHK-SE-Group/generic-generator/schemas"
	"github.com/CUHK-SE-Group/generic-generator/schemas/query"
	"math"
	"math/rand"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestDefaultHandler(t *testing.T) {
	cons := schemas.DefinedBeforeUse
	cons.FirstNode = "expr/id"
	cons.SecondNode = "id"

	schemas.Register(cons)
	g, err := parser.Parse("./testdata/complete/tinyc.ebnf", "program")
	if err != nil {
		panic(err)
	}
	g.MergeProduction()
	g.BuildShortestNotation()
	chain, err := schemas.CreateChain("test", &schemas.TraceHandler{}, &schemas.CatHandler{}, &schemas.IDHandler{}, &schemas.SubHandler{}, &WeightedHandler{}, &schemas.TermHandler{}, &schemas.RepHandler{}, &schemas.BracketHandler{})
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

	ctx.SymbolStack.Pop()
	switch ctx.Mode {
	case schemas.ShrinkMode:
		sym := cur.GetSymbols()
		candidates := make([]int, 0)
		repechage := make([]int, 0)
		for i, v := range sym {
			if v.GetDistance() < cur.GetDistance() {
				candidates = append(candidates, i)
			} else {
				repechage = append(repechage, i)
			}
		}
		if len(candidates) == 0 {
			candidates = repechage
		}
		idx := rand.Intn(len(candidates))
		votes := 0
		for i, v := range sym {
			if i != candidates[idx] {
				votes += ctx.VisitedEdge[schemas.GetEdgeID(cur.GetID(), v.GetID())]
			}
		}
		if votes > 10 && ctx.VisitedEdge[schemas.GetEdgeID(cur.GetID(), sym[candidates[idx]].GetID())] > 4*votes {
			// if it goes into this branch, it means it chooses too much times this path, which indicates that there is a big probability of circle
			idx = rand.Intn(len(sym)) //then re-vote for all the branches
			ctx.VisitedEdge[schemas.GetEdgeID(cur.GetID(), sym[idx].GetID())]++
			ctx.SymbolStack.Push(sym[idx])
		} else {
			ctx.VisitedEdge[schemas.GetEdgeID(cur.GetID(), sym[candidates[idx]].GetID())]++
			ctx.SymbolStack.Push(sym[candidates[idx]])
		}
	default:
		idx := rand.Int() % len(cur.GetSymbols())
		ctx.SymbolStack.Push((cur.GetSymbols())[idx])
		ctx.VisitedEdge[schemas.GetEdgeID(cur.GetID(), (cur.GetSymbols())[idx].GetID())]++
	}

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
	g, err := parser.Parse("./testdata/complete/tinyc.ebnf", "program")
	if err != nil {
		panic(err)
	}
	g.MergeProduction()
	g.BuildShortestNotation()
	graph.Visualize(g.GetInternal(), "fig.dot", func(vertex graph.Vertex[schemas.Property]) string {
		return fmt.Sprintf("%s,%d", vertex.GetID(), vertex.GetProperty(schemas.Prop).DistanceToTerminal)
	})
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
		})
	}
	fmt.Println(ctx.Result)
}

func TestWeightHandlerManyTimes(t *testing.T) {
	g, err := parser.Parse("./testdata/complete/tinyc.ebnf", "program")
	if err != nil {
		panic(err)
	}
	g.MergeProduction()
	g.BuildShortestNotation()
	chain, err := schemas.CreateChain("test", &schemas.CatHandler{}, &schemas.IDHandler{}, &schemas.TermHandler{}, &WeightedHandler{}, &schemas.OrHandler{})
	if err != nil {
		panic(err)
	}
	t1 := time.Now()
	num := 1000000
	defer func() {
		duration := time.Since(t1)
		fmt.Printf("generated %d instances, use %s\n", num, duration)
	}()
	results := make(map[string]int)
	for i := 0; i < num; i++ {
		ctx, err := schemas.NewContext(g, "program", context.Background(), nil)
		if err != nil {
			panic(err)
		}
		for !ctx.GetFinish() {
			chain.Next(ctx, func(result *schemas.Result) {
				ctx = result.GetCtx()
				ctx.HandlerIndex = 0
			})
		}
		results[ctx.Result]++
	}
	fmt.Println(results)

}

type MonitorHandler struct {
}

func (h *MonitorHandler) Handle(chain *schemas.Chain, ctx *schemas.Context, cb schemas.ResponseCallBack) {
	cur := ctx.SymbolStack.Top()
	if strings.Contains(cur.GetID(), "Expression") {
		fmt.Println(cur.GetContent())
	}
	for key, fn := range schemas.FirstPlaceNode {
		if len(ctx.SymbolStack.ProductionTrace) > 0 && !strings.Contains(cur.GetID(), ctx.SymbolStack.ProductionTrace[len(ctx.SymbolStack.ProductionTrace)-1]) {
			continue
		}
		if query.MatchPattern(ctx.SymbolStack.ProductionTrace, key) {
			ctx = fn(ctx)
		}
	}

	chain.Next(ctx, cb)

	for key, fn := range schemas.SecondPlaceNode {
		if len(ctx.SymbolStack.ProductionTrace) > 0 && !strings.Contains(cur.GetID(), ctx.SymbolStack.ProductionTrace[len(ctx.SymbolStack.ProductionTrace)-1]) {
			continue
		}
		if query.MatchPattern(ctx.SymbolStack.ProductionTrace, key) {
			var err error
			ctx, err = fn.Func(ctx)
			if err != nil {
				if errors.Is(err, schemas.ErrSymbolNotFound) {
					continue
				}
				panic(err)
			}
		}
	}

}

func (h *MonitorHandler) HookRoute() []regexp.Regexp {
	return make([]regexp.Regexp, 0)
}

func (h *MonitorHandler) Name() string {
	return "monitor"
}

func (h *MonitorHandler) Type() schemas.GrammarType {
	return math.MaxInt
}

func TestDefaultHandlerCypher(t *testing.T) {
	cons := schemas.MaxLimit
	cons.FirstNode = "Expression"
	cons.SecondNode = "Expression"
	cons2 := schemas.MaxLimit
	cons2.FirstNode = "Expression"
	cons2.SecondNode = "Expression"
	schemas.Register(cons)
	g, err := parser.Parse("./testdata/complete/Cypher.ebnf", "Cypher")
	if err != nil {
		panic(err)
	}
	g.MergeProduction()
	g.BuildShortestNotation()
	chain, err := schemas.CreateChain("test", &schemas.OptionHandler{}, &MonitorHandler{}, &schemas.CatHandler{}, &schemas.IDHandler{}, &schemas.SubHandler{}, &WeightedHandler{}, &schemas.TermHandler{}, &schemas.RepHandler{}, &schemas.BracketHandler{})
	if err != nil {
		panic(err)
	}
	for i := 0; i < 10; i++ {
		ctx, err := schemas.NewContext(g, "Cypher", context.Background(), nil)
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
}
