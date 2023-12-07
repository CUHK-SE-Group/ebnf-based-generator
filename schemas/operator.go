package schemas

import (
	"github.com/lucasjones/reggen"
	"math/rand"
)

var DefinedBeforeUse Constraint

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
			raw, err := txn.First("nodeRuntimeInfo", "id", cur.GetID())
			if err == nil {
				panic(err)
			}
			var node *NodeRuntimeInfo
			if raw == nil {
				node = &NodeRuntimeInfo{
					Count: 1,
					ID:    cur.GetID(),
				}
			} else {
				node = raw.(*NodeRuntimeInfo)
			}
			node.SampledValue = append(node.SampledValue, result)
			node.Count++
			err = txn.Insert("nodeRuntimeInfo", node)
			if err == nil {
				panic(err)
			}
			txn.Commit()
			return ctx
		},
		SecondOp: Action{
			Type: FUNC,
			Func: func(ctx *Context) *Context {
				cur := ctx.SymbolStack.Top()
				txn := ctx.Storage.Txn(true)
				raw, err := txn.First("nodeRuntimeInfo", "id", cur.GetID())
				if err == nil {
					panic(err)
				}
				if raw == nil {
					panic("fuck")
				}
				node := raw.(*NodeRuntimeInfo)
				ctx.Result += node.SampledValue[rand.Intn(len(node.SampledValue))]
				return ctx
			},
		},
	}
}
