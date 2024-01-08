package Generator

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/CUHK-SE-Group/generic-generator/graph"
	"github.com/CUHK-SE-Group/generic-generator/parser"
	"github.com/CUHK-SE-Group/generic-generator/schemas"
	"github.com/CUHK-SE-Group/generic-generator/schemas/query"
	"log"
	"log/slog"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"runtime/pprof"
	"testing"
	"time"
)

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
		//if (votes > 0 && ctx.VisitedEdge[schemas.GetEdgeID(cur.GetID(), sym[candidates[idx]].GetID())] > 3*votes) || (votes == 0 && ctx.VisitedEdge[schemas.GetEdgeID(cur.GetID(), sym[candidates[idx]].GetID())] > 20) {
		//	// if it goes into this branch, it means it chooses too much times this path, which indicates that there is a big probability of circle
		//	idx = rand.Intn(len(sym)) //then re-vote for all the branches
		//	ctx.VisitedEdge[schemas.GetEdgeID(cur.GetID(), sym[idx].GetID())]++
		//	ctx.SymbolStack.Push(sym[idx])
		//	ctx.Result.AddEdge(cur, sym[idx])
		//} else {
		ctx.VisitedEdge[schemas.GetEdgeID(cur.GetID(), sym[candidates[idx]].GetID())]++
		ctx.SymbolStack.Push(sym[candidates[idx]])
		ctx.Result.AddEdge(cur, sym[candidates[idx]])
		//}
	default:
		idx := rand.Int() % len(cur.GetSymbols())
		ctx.SymbolStack.Push((cur.GetSymbols())[idx])
		ctx.VisitedEdge[schemas.GetEdgeID(cur.GetID(), (cur.GetSymbols())[idx].GetID())]++
		ctx.Result.AddEdge(cur, (cur.GetSymbols())[idx])
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

type MonitorHandler struct {
}

func (h *MonitorHandler) Handle(chain *schemas.Chain, ctx *schemas.Context, cb schemas.ResponseCallBack) {
	if ctx.Constraint == nil {
		chain.Next(ctx, cb)
		return
	}
	constraints := ctx.Constraint.GetConstraints()
	for _, v := range constraints {
		if query.MatchPattern(ctx.SymbolStack.ProductionTrace, v.FirstNode) {
			switch v.FirstOp.Type {
			case schemas.FUNC:
				ctx, _ = v.FirstOp.Func(ctx)
			case schemas.REGEX:

			}
		}
		if query.MatchPattern(ctx.SymbolStack.ProductionTrace, v.SecondNode) {
			switch v.SecondOp.Type {
			case schemas.FUNC:
				ctx, _ = v.SecondOp.Func(ctx)
			case schemas.REGEX:

			}
		}
	}
	chain.Next(ctx, cb)
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

func TestDefaultHandler(t *testing.T) {
	cons := schemas.DefinedBeforeUse
	cons.FirstNode = "expr/id"
	cons.SecondNode = "id"

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
	ctx, err := schemas.NewContext(g, "program", context.Background(), nil, nil)
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

func TestWeightedHandler(t *testing.T) {
	g, err := parser.Parse("./testdata/complete/tinyc.ebnf", "program")
	if err != nil {
		panic(err)
	}
	cons := schemas.MaxLimit
	cons.FirstNode = "paren_expr"
	cons.SecondNode = "paren_expr"
	consg := schemas.NewConstraintGraph()
	consg.AddBinaryConstraint(cons)
	g.MergeProduction()
	g.BuildShortestNotation()
	chain, err := schemas.CreateChain("test", &MonitorHandler{}, &schemas.CatHandler{}, &schemas.IDHandler{}, &schemas.TermHandler{}, &WeightedHandler{}, &schemas.OrHandler{}, &schemas.RepHandler{}, &schemas.BracketHandler{})
	if err != nil {
		panic(err)
	}
	ctx, err := schemas.NewContext(g, "program", context.Background(), consg, nil)
	if err != nil {
		panic(err)
	}
	for !ctx.GetFinish() {
		chain.Next(ctx, func(result *schemas.Result) {
			ctx = result.GetCtx()
			ctx.HandlerIndex = 0
		})
	}
	fmt.Println(ctx.Result.GetResult())
	fmt.Printf("edge coverage: %d/%d\n", len(ctx.VisitedEdge), len(ctx.Grammar.GetInternal().GetAllEdges()))
	err = graph.Visualize(ctx.Result.Grammar.GetInternal(), "fig.dot", nil)
	if err != nil {
		panic(err)
	}

	timeout := 1 * time.Second
	ctxtime, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel() // 确保所有路径上都调用了cancel

	input := ctx.Result.GetResult()
	cmd := exec.CommandContext(ctxtime, "./tinyc")
	var in bytes.Buffer
	in.Write([]byte(input))
	cmd.Stdin = &in
	output, err := cmd.Output()
	if err != nil {
		if errors.Is(ctxtime.Err(), context.DeadlineExceeded) {
			fmt.Println("命令执行超时")
			return
		}
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			log.Printf("命令执行失败: %s\n标准错误输出:%s\n", exitErr.Error(), exitErr.Stderr)
		}
	}
	fmt.Println(string(output))
}

func TestWeightHandlerManyTimes(t *testing.T) {
	cpuFile, err := os.Create("cpu.prof")
	if err != nil {
		t.Fatalf("could not create CPU profile: %v", err)
	}
	defer cpuFile.Close()

	if err := pprof.StartCPUProfile(cpuFile); err != nil {
		t.Fatalf("could not start CPU profile: %v", err)
	}
	defer pprof.StopCPUProfile()

	g, err := parser.Parse("./testdata/complete/tinyc.ebnf", "program")
	if err != nil {
		panic(err)
	}
	g.MergeProduction()
	g.BuildShortestNotation()
	chain, err := schemas.CreateChain("test", &MonitorHandler{}, &schemas.CatHandler{}, &schemas.IDHandler{}, &schemas.TermHandler{}, &WeightedHandler{}, &schemas.OrHandler{}, &schemas.RepHandler{}, &schemas.BracketHandler{})
	if err != nil {
		panic(err)
	}
	t1 := time.Now()
	num := 1000
	defer func() {
		duration := time.Since(t1)
		fmt.Printf("generated %d instances, use %s\n", num, duration)
	}()
	for i := 0; i < num; i++ {
		ctx, err := schemas.NewContext(g, "program", context.Background(), nil, nil)
		ctx.Mode = schemas.ShrinkMode
		if err != nil {
			panic(err)
		}
		for !ctx.GetFinish() {
			chain.Next(ctx, func(result *schemas.Result) {
				ctx = result.GetCtx()
				ctx.HandlerIndex = 0
			})
		}
		fmt.Println(ctx.Result.GetResult())
		fmt.Printf("edge coverage: %d/%d\n", len(ctx.VisitedEdge), len(ctx.Grammar.GetInternal().GetAllEdges()))
	}

	memFile, err := os.Create("mem.prof")
	if err != nil {
		t.Fatalf("could not create memory profile: %v", err)
	}
	defer memFile.Close()
	runtime.GC() // GC, to get a clean memory profile
	if err := pprof.WriteHeapProfile(memFile); err != nil {
		t.Fatalf("could not write memory profile: %v", err)
	}
}

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
	fmt.Println(ctx.Result.GetResult())
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
	fmt.Printf("%s\n", ctx.Result.GetResult())
}

type EBPFTermHandler struct {
}

func (h *EBPFTermHandler) isTermPreserve(g *schemas.Node) bool {
	content := g.GetContent()
	return (content[0] == content[len(content)-1]) && ((content[0] == '\'') || content[0] == '"')
}

func (h *EBPFTermHandler) stripQuote(content string) string {
	if content[0] == content[len(content)-1] {
		if (content[0] == '\'') || (content[0] == '"') {
			return content[1 : len(content)-1]
		}
	}
	return content
}

func (h *EBPFTermHandler) Handle(chain *schemas.Chain, ctx *schemas.Context, cb schemas.ResponseCallBack) {
	cur := ctx.SymbolStack.Top()
	ctx.SymbolStack.Pop()

	if len(cur.GetSymbols()) != 0 {
		slog.Error("Pattern mismatched[Terminal]")
		return
	}
	ctx.Result.AddEdge(cur, cur) // 用一个自环标记到达了最后的终结符节点
	chain.Next(ctx, cb)
}

func (h *EBPFTermHandler) HookRoute() []regexp.Regexp {
	return make([]regexp.Regexp, 0)
}

func (h *EBPFTermHandler) Name() string {
	return "ebpfterminalhandler"
}

func (h *EBPFTermHandler) Type() schemas.GrammarType {
	return schemas.GrammarTerminal
}

func TestEBPFHandler(t *testing.T) {
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
	fmt.Printf("%s\n", ctx.Result.GetResult())
}
