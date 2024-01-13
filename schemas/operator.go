package schemas

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"

	"github.com/lucasjones/reggen"
)

var (
	DefinedBeforeUse Constraint
	MaxLimit         Constraint
)
var (
	ErrSymbolNotFound = errors.New("notfound")
	ErrPassThrough    = errors.New("pass through the process")
	ErrIntercept      = errors.New("intercept the handlers")
)

func trimNumber(id string) string {
	if !strings.Contains(id, "#") {
		return id
	}
	return strings.Split(id, "#")[0]
}

func randomKey(m map[string]int) string {
	if len(m) == 0 {
		return ""
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys[rand.Intn(len(keys))]
}

func init() {
	DefinedBeforeUse = Constraint{
		FirstOp: Action{
			Type: FUNC,
			Func: func(ctx *Context) (*Context, error) {
				cur := ctx.SymbolStack.Top()
				if cur.GetType() != GrammarTerminal {
					return ctx, errors.New("the symbol type is not GrammarTerminal")
				}
				result, err := reggen.Generate(cur.GetContent(), 1)
				if err != nil {
					return ctx, errors.Join(err, errors.New("reggen failed"))
				}
				//ctx.Result += result
				fmt.Println(result)

				txn := ctx.Storage.Txn(true)
				raw, err := txn.First("nodeRuntimeInfo", "id", trimNumber(cur.GetID()))
				if err != nil {
					panic(err)
				}
				var node *NodeRuntimeInfo
				if raw == nil {
					node = &NodeRuntimeInfo{
						Count:        1,
						ID:           trimNumber(cur.GetID()),
						SampledValue: make(map[string]int),
					}
				} else {
					node = raw.(*NodeRuntimeInfo)
				}
				node.SampledValue[result]++
				node.Count++
				err = txn.Insert("nodeRuntimeInfo", node)
				if err != nil {
					panic(err)
				}
				txn.Commit()
				ctx.SymbolStack.Pop()
				return ctx, nil
			},
		},
		SecondOp: Action{
			Type: FUNC,
			Func: func(ctx *Context) (*Context, error) {
				cur := ctx.SymbolStack.Top()
				txn := ctx.Storage.Txn(false)
				raw, err := txn.First("nodeRuntimeInfo", "id", trimNumber(cur.GetID()))
				if err != nil {
					panic(err)
				}
				if raw == nil {
					return ctx, ErrSymbolNotFound
				}
				node := raw.(*NodeRuntimeInfo)
				if node.Count > 5 {
					ctx.Mode = ShrinkMode
				}
				//ctx.Result += randomKey(node.SampledValue)
				fmt.Println(randomKey(node.SampledValue))
				ctx.SymbolStack.Pop()
				return ctx, nil
			},
		},
	}
	MaxLimit = Constraint{
		FirstOp: Action{
			Type: FUNC,
			Func: func(ctx *Context) (*Context, error) {
				cur := ctx.SymbolStack.Top()
				txn := ctx.Storage.Txn(true)
				raw, err := txn.First("nodeRuntimeInfo", "id", trimNumber(cur.GetID()))
				if err != nil {
					panic(err)
				}
				var node *NodeRuntimeInfo
				if raw == nil {
					node = &NodeRuntimeInfo{
						Count:        1,
						ID:           trimNumber(cur.GetID()),
						SampledValue: make(map[string]int),
					}
				} else {
					node = raw.(*NodeRuntimeInfo)
				}
				node.Count++
				if node.Count > 3 {
					ctx.Mode = ShrinkMode
				}
				err = txn.Insert("nodeRuntimeInfo", node)
				if err != nil {
					panic(err)
				}
				txn.Commit()
				return ctx, ErrPassThrough
			},
		},
		SecondOp: Action{
			Type: FUNC,
			Func: func(ctx *Context) (*Context, error) {
				cur := ctx.SymbolStack.Top()
				txn := ctx.Storage.Txn(false)
				raw, err := txn.First("nodeRuntimeInfo", "id", trimNumber(cur.GetID()))
				if err != nil {
					panic(err)
				}
				if raw == nil {
					return ctx, nil
				}
				node := raw.(*NodeRuntimeInfo)
				if node.Count > 10 {
					ctx.Mode = ShrinkMode
				}
				return ctx, ErrPassThrough
			},
		},
	}
}
