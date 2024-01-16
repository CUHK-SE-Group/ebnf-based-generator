package Generator

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/CUHK-SE-Group/generic-generator/graph"
	"github.com/CUHK-SE-Group/generic-generator/parser"
	"github.com/CUHK-SE-Group/generic-generator/schemas"
	"log"
	"math"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"runtime/pprof"
	"testing"
	"time"
)

type dummyHandler struct {
}

func (h *dummyHandler) HookRoute() []regexp.Regexp {
	return make([]regexp.Regexp, 0)
}

func (h *dummyHandler) Name() string {
	return "dummy"
}

func (h *dummyHandler) Type() schemas.GrammarType {
	return math.MaxInt
}
func (h *dummyHandler) Handle(chain *schemas.Chain, ctx *schemas.Context, cb schemas.ResponseCallBack) {
	ctx.CurrentNode = ctx.SymbolStack.Top() // clear the message exchange buffer
	ctx.ResultBuffer = make([]*schemas.Node, 0)

	chain.Next(ctx, cb)

	if ctx.Error != nil {
		panic(ctx.Error)
		return
	}
	ctx.SymbolStack.Pop()
	ctx.SymbolStack.Push(ctx.ResultBuffer...)
	for i := len(ctx.ResultBuffer) - 1; i >= 0; i-- {
		ctx.Result.AddNode(ctx.ResultBuffer[i])
		ctx.Result.AddEdge(ctx.CurrentNode, ctx.ResultBuffer[i])
	}
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
	chain, err := schemas.CreateChain("test", &dummyHandler{}, &schemas.TraceHandler{}, &schemas.CatHandler{}, &schemas.IDHandler{}, &schemas.SubHandler{}, &WeightedHandler{}, &schemas.TermHandler{}, &schemas.RepHandler{}, &schemas.BracketHandler{})
	if err != nil {
		panic(err)
	}

	parentCtx, _ := context.WithTimeout(context.Background(), time.Second)
	ctx, err := schemas.NewContext(g, "program", parentCtx, nil, nil)
	if err != nil {
		panic(err)
	}

	for !ctx.GetFinish() {
		select {
		case <-parentCtx.Done():
			// Handle the case when the context's deadline is exceeded
			fmt.Println("Operation timed out")
			return
		default:
			chain.Next(ctx, func(result *schemas.Result) {
				ctx = result.GetCtx()
			})
			ctx.HandlerIndex = 0
			fmt.Println(ctx.Result.GetResult(nil))
		}

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
	chain, err := schemas.CreateChain("test", &MonitorHandler{}, &dummyHandler{}, &schemas.IDHandler{}, &schemas.CatHandler{}, &WeightedHandler{}, &schemas.RepHandler{}, &schemas.BracketHandler{}, &schemas.TermHandler{})
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
		})
		ctx.HandlerIndex = 0
	}
	err = ctx.Result.Save("/tmp/grammarfile")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(ctx.Result.GetResult(nil))
	fmt.Printf("edge coverage: %d/%d\n", len(ctx.VisitedEdge), len(ctx.Grammar.GetInternal().GetAllEdges()))
	err = graph.Visualize(ctx.Result.Grammar.GetInternal(), "fig.dot", nil, nil)
	if err != nil {
		panic(err)
	}

	input := ctx.Result.GetResult(nil)
	validateInput(input)
}

func TestHandlerChainMux(t *testing.T) {
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

	routerHandler := &WrapHandler{Chain: map[schemas.GrammarType]*schemas.Chain{}}
	err = routerHandler.Register(wrapChain(&schemas.IDHandler{}), wrapChain(&schemas.CatHandler{}), wrapChain(&WeightedHandler{}), wrapChain(&schemas.RepHandler{}), wrapChain(&schemas.BracketHandler{}), wrapChain(&schemas.TermHandler{}))
	if err != nil {
		panic(err)
	}
	chain, err := schemas.CreateChain("main", &MonitorHandler{}, routerHandler)
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
		})
		ctx.HandlerIndex = 0
	}

	err = ctx.Result.Save("/tmp/grammarfile")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(ctx.Result.GetResult(nil))
	fmt.Printf("edge coverage: %d/%d\n", len(ctx.VisitedEdge), len(ctx.Grammar.GetInternal().GetAllEdges()))
	err = graph.Visualize(ctx.Result.Grammar.GetInternal(), "fig.dot", nil, nil)
	if err != nil {
		panic(err)
	}

	input := ctx.Result.GetResult(nil)
	validateInput(input)
}

func validateInput(input string) {
	timeout := 1 * time.Second
	ctxtime, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel() // 确保所有路径上都调用了cancel
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

func TestMutate(t *testing.T) {
	g := schemas.NewGrammar(schemas.WithLoadFromFile("/tmp/grammarfile"))
	err := graph.Visualize(g.GetInternal(), "fig1.dot", func(vertex graph.Vertex[schemas.Property]) string {
		return fmt.Sprintf("%s\n%s", vertex.GetID(), vertex.GetProperty(schemas.Prop).Content)
	}, nil)
	if err != nil {
		t.Error(err)
	}
	g.PrintTerminals("program#0")
}

func TestSaveAndLoad(t *testing.T) {
	g, err := parser.Parse("./testdata/complete/tinyc.ebnf", "program")
	if err != nil {
		panic(err)

	}
	err = g.Save("/tmp/grammarfile")
	if err != nil {
		t.Error(err)
	}
	newg := schemas.NewGrammar(schemas.WithLoadFromFile("/tmp/grammarfile"))
	n := newg.GetNode("program")
	n.GetSymbols() // refresh the cache

	if len(g.GetInternal().GetAllEdges()) != len(newg.GetInternal().GetAllEdges()) {
		t.Error("the edge num should be equal")
	}
	if len(g.GetInternal().GetAllVertices()) != len(newg.GetInternal().GetAllVertices()) {
		t.Error("the vertex num should be equal")
	}
	if len(g.GetInternal().GetAllMetadata()) != len(newg.GetInternal().GetAllMetadata()) {
		t.Error("the metadata num should be equal")
	}
	for _, v := range g.GetInternal().GetAllEdges() {
		e := newg.GetInternal().GetEdgeById(v.GetID())
		if e == nil || e.GetFrom().GetID() != v.GetFrom().GetID() || e.GetTo().GetID() != v.GetTo().GetID() {
			t.Errorf("%s not found in new grammar", v.GetID())
		}
	}
	for _, v := range g.GetInternal().GetAllVertices() {
		e := newg.GetInternal().GetVertexById(v.GetID())
		if e == nil || e.GetProperty(schemas.Prop).Content != v.GetProperty(schemas.Prop).Content || e.GetProperty(schemas.Prop).Type != v.GetProperty(schemas.Prop).Type || e.GetProperty(schemas.Prop).DistanceToTerminal != v.GetProperty(schemas.Prop).DistanceToTerminal {
			t.Errorf("%s not found in new grammar", v.GetID())
		}
	}
	for k, v := range g.GetInternal().GetAllMetadata() {
		if newg.GetInternal().GetMetadata(k) != v {
			t.Errorf("%s not found in new grammar", k)
		}
	}
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
		fmt.Println(ctx.Result.GetResult(nil))
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
