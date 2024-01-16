package Generator

import (
	"fmt"
	"github.com/CUHK-SE-Group/generic-generator/schemas"
	"github.com/CUHK-SE-Group/generic-generator/schemas/query"
	"math"
	"math/rand"
	"regexp"
	"strings"
)

type WeightedHandler struct {
}

func (h *WeightedHandler) Handle(chain *schemas.Chain, ctx *schemas.Context, cb schemas.ResponseCallBack) {
	//cur := ctx.SymbolStack.Top()
	if len(ctx.CurrentNode.GetSymbols()) == 0 {
		chain.Next(ctx, cb)
		return
	}

	//ctx.SymbolStack.Pop()
	switch ctx.Mode {
	case schemas.ShrinkMode:
		sym := ctx.CurrentNode.GetSymbols()
		candidates := make([]int, 0)
		repechage := make([]int, 0)
		for i, v := range sym {
			if v.GetDistance() < ctx.CurrentNode.GetDistance() {
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
				votes += ctx.VisitedEdge[schemas.GetEdgeID(ctx.CurrentNode.GetID(), v.GetID())]
			}
		}
		ctx.VisitedEdge[schemas.GetEdgeID(ctx.CurrentNode.GetID(), sym[candidates[idx]].GetID())]++
		ctx.ResultBuffer = append(ctx.ResultBuffer, sym[candidates[idx]])
	default:
		idx := rand.Int() % len(ctx.CurrentNode.GetSymbols())
		ctx.ResultBuffer = append(ctx.ResultBuffer, ctx.CurrentNode.GetSymbol(idx))
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
	trace := append(ctx.SymbolStack.ProductionTrace, strings.Split(ctx.SymbolStack.Top().GetID(), "#")[0])
	for _, v := range constraints {
		if query.MatchPattern(trace, v.FirstNode) {
			switch v.FirstOp.Type {
			case schemas.FUNC:
				ctx, _ = v.FirstOp.Func(ctx)
			case schemas.REGEX:

			}
		}
		if query.MatchPattern(trace, v.SecondNode) {
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

type WrapHandler struct {
	Chain map[schemas.GrammarType]*schemas.Chain
}

func (h *WrapHandler) Handle(chain *schemas.Chain, ctx *schemas.Context, cb schemas.ResponseCallBack) {
	// save and restore the environment
	handlerIndex := ctx.HandlerIndex
	ctx.HandlerIndex = 0
	defer func() {
		ctx.HandlerIndex = handlerIndex
	}()

	ctx.CurrentNode = ctx.SymbolStack.Top() // clear the message exchange buffer
	ctx.ResultBuffer = make([]*schemas.Node, 0)
	fmt.Println("dealing with ", ctx.CurrentNode.GetID())
	// route the request to different Chain
	for k, c := range h.Chain {
		if k&ctx.CurrentNode.GetType() != 0 {
			c.Next(ctx, cb)
			if ctx.Error != nil {
				panic(ctx.Error)
				return
			}
			if len(ctx.ResultBuffer) == 0 && ctx.CurrentNode.GetType() == schemas.GrammarTerminal {
				if len(ctx.Tmp1) == 0 {
					ctx.Tmp1 = make([]string, 0)
				}
				ctx.Tmp1 = append(ctx.Tmp1, strings.Trim(ctx.CurrentNode.GetContent(), "'"))
				//fmt.Println(strings.Join(ctx.Tmp1, ""))
				fmt.Printf("============in the final handler: ")
				for i := len(ctx.Tmp1) - 1; i >= 0; i-- {
					fmt.Printf("%s", ctx.Tmp1[i])
				}
				fmt.Println()
				ctx.Result.AddEdge(ctx.CurrentNode, ctx.CurrentNode)
			}
			ctx.SymbolStack.Pop()
			ctx.SymbolStack.Push(ctx.ResultBuffer...)
			for i := len(ctx.ResultBuffer) - 1; i >= 0; i-- {
				ctx.Result.AddNode(ctx.ResultBuffer[i])
				ctx.Result.AddEdge(ctx.CurrentNode, ctx.ResultBuffer[i])
			}
			return
		}
	}
	panic(fmt.Errorf("there is no such handler that can deal with %v", ctx.CurrentNode.GetType()))
}

func (h *WrapHandler) HookRoute() []regexp.Regexp {
	return make([]regexp.Regexp, 0)
}

func (h *WrapHandler) Name() string {
	return "mux"
}

func (h *WrapHandler) Type() schemas.GrammarType {
	return math.MaxInt
}

func (h *WrapHandler) Register(chain ...*schemas.Chain) error {
	for _, v := range chain {
		if len(v.Handlers) == 0 {
			return fmt.Errorf("the length of chain should not be zero")
		}
		h.Chain[v.Handlers[0].Type()] = v // note: every grammarType should only has one handler chain
	}
	return nil
}

func wrapChain(h schemas.Handler) *schemas.Chain {
	chain, _ := schemas.CreateChain(h.Name(), h)
	return chain
}
