package grammar

import (
	"fmt"
	"testing"
)

func TestDerivationTree_ExpandNode(t *testing.T) {
	tree := DerivationTree{
		gram:        SimpleNonterminalGrammar,
		startSymbol: StartSymbol,
		traversal:   BFS,
		expandAlgo:  MinimalCostExpand,
	}
	tree.Construct()
	cnt := 0
	for !tree.ExpandNode() {
		tree.Visualize(fmt.Sprintf("tree%d.png", cnt))
		tree.GetNonTerminals()
		//for _, v := range tree.GetNonTerminals() {
		//	fmt.Printf("%s ", v.GetName())
		//}
		fmt.Println()
		cnt++
	}
	tree.Visualize(fmt.Sprintf("tree%d.png", -1))
}
