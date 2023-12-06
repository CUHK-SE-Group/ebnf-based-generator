package schemas

var DefinedBeforeUse Constraint

func init() {

	DefinedBeforeUse = Constraint{
		FirstOp: func(ctx *Context) *Context {
			cur := ctx.SymbolStack.Top()
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
			}
			node = raw.(*NodeRuntimeInfo)
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
				return ctx
			},
		},
	}
}
