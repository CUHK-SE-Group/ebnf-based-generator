package grammar

import (
	"fmt"
	"testing"
)

func TestDerivationTree_ExpandNode(t *testing.T) {
	tree := DerivationTree{
		gram:        URLGrammar,
		startSymbol: StartSymbol,
		traversal:   DFS,
		expandAlgo:  MinimalCostExpand,
	}
	tree.Construct()
	cnt := 0
	finish := false
	for !finish {
		finish = tree.ExpandNode()
		tree.Visualize(fmt.Sprintf("tree%d.png", cnt))
		cnt++
	}
	_, res := tree.GetLeafNodes()
	if res != "http://cispa.saarland:80" {
		t.Errorf("ExpandNode error, actual: %s, want: %s", res, "http://cispa.saarland:80")
	}
}
