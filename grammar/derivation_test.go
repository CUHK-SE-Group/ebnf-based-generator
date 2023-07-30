package grammar

import (
	"context"
	"fmt"
	"testing"
)

func TestDerivationTree_ExpandNode(t *testing.T) {
	tree := NewDerivationTree(URLGrammar, StartSymbol, DFS, MinimalCostExpand)
	cnt := 0
	finish := false
	ctx := NewDerivationContext(context.Background())
	ctx = tree.Configure(ctx)
	for !finish {
		_, finish = tree.ExpandNode(ctx)
		ctx.handlerIndex = 0 // reset the handler index for test
		tree.Visualize(fmt.Sprintf("tree%d.png", cnt))
		cnt++
	}
	_, res := tree.GetLeafNodes()
	if res != "http://cispa.saarland:80" {
		t.Errorf("ExpandNode error, actual: %s, want: %s", res, "http://cispa.saarland:80")
	}
}

func TestDerivationTree_RandomExpand(t *testing.T) {
	tree := NewDerivationTree(URLGrammar, StartSymbol, DFS, RandomExpand)
	cnt := 0
	finish := false
	ctx := NewDerivationContext(context.Background())
	ctx = tree.Configure(ctx)
	for !finish {
		_, finish = tree.ExpandNode(ctx)
		ctx.handlerIndex = 0 // reset the handler index for test
		cnt++
		_, res := tree.GetLeafNodes()
		fmt.Println(res)
	}

}

func TestDerivationTree_GetCoverage(t *testing.T) {
	tree := NewDerivationTree(URLGrammar, StartSymbol, DFS, MinimalCostExpand)
	cnt := 0
	finish := false
	ctx := NewDerivationContext(context.Background())
	ctx = tree.Configure(ctx)
	for !finish {
		_, finish = tree.ExpandNode(ctx)
		ctx.handlerIndex = 0 // reset the handler index for test
		tree.Visualize(fmt.Sprintf("tree%d.png", cnt))
		cnt++
	}
	_, res := tree.GetLeafNodes()
	if res != "http://cispa.saarland:80" {
		t.Errorf("ExpandNode error, actual: %s, want: %s", res, "http://cispa.saarland:80")
	}
	_, cov := ctx.GetCoverage()
	if cov != 0.17073170731707318 {
		t.Errorf("GetCoverage error, actual: %f, want: %f", cov, 0.17073170731707318)
	}
}
