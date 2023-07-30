package grammar

import "fmt"

func countCoverage(ctx *DerivationContext, tree *DerivationTree) {
	ctx.Next()(ctx, tree)
	for _, pair := range ctx.preExpansions {
		if _, ok := tree.gram.G[pair.From.GetName()]; ok {
			ctx.coverage[fmt.Sprintf("%s->%s", pair.From.GetName(), pair.To.GetName())]++
		}
	}
}
func checkFinish(ctx *DerivationContext, tree *DerivationTree) {
	ctx.Next()(ctx, tree)
	finish := true
	tree.traversal(tree.root, func(node *Node) {
		if node == nil {
			return
		}
		if node.Children == nil && IsNonTerminals(node.GetName()) {
			finish = false
		}
	})
	if finish {
		ctx.finish = true
	}
}
