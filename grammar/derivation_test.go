package grammar

import (
	"fmt"
	"testing"
	"time"
)

func TestDerivationTree_ExpandNode(t *testing.T) {
	tree := DerivationTree{
		gram:        ExprGrammar,
		startSymbol: StartSymbol,
		traversal:   BFS,
		expandAlgo:  RandomExpand,
	}
	tree.Construct()
	cnt := 0
	for !tree.ExpandNode() {
		fmt.Println("-")
		tree.Visualize(fmt.Sprintf("tree%d.png", cnt))
		for _, v := range tree.GetNonTerminals() {
			fmt.Printf("%s ", v.GetName())
		}
		fmt.Println()
		time.Sleep(time.Second)
		cnt++
	}
	tree.Visualize(fmt.Sprintf("tree%d.png", -1))
}
