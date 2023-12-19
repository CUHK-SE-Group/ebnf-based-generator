package schemas

import (
	"errors"
	"fmt"
	"github.com/lucasjones/reggen"
	"math/rand"
	"strings"
)

var (
	DefinedBeforeUse Constraint
	MaxLimit         Constraint
)
var (
	ErrSymbolNotFound = errors.New("notfound")
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
		FirstOp: func(ctx *Context) *Context {
			cur := ctx.SymbolStack.Top()
			if cur.GetType() != GrammarTerminal {
				return nil
			}
			result, err := reggen.Generate(cur.GetContent(), 1)
			if err != nil {
				return nil
			}
			ctx.Result += result

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
				fmt.Println("add a new symbol", trimNumber(cur.GetID()), cur.GetContent())
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
			return ctx
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
				ctx.Result += randomKey(node.SampledValue)
				ctx.SymbolStack.Pop()
				return ctx, nil
			},
		},
	}
	MaxLimit = Constraint{
		FirstOp:  nil,
		SecondOp: Action{},
	}
}
