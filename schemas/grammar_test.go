package schemas_test

import (
	"fmt"
	"testing"

	"github.com/CUHK-SE-Group/generic-generator/graph"
	"github.com/CUHK-SE-Group/generic-generator/schemas"
)

func TestNewGrammar(t *testing.T) {
	t.Run("CreateGrammarAndVerifyNodes", func(t *testing.T) {
		g := createSampleGrammar()

		// Visualize the graph for debugging purposes
		graph.Visualize(g.GetInternal(), "./fig.dot", nil)

		verifyGrammarNodes(t, g)
	})

	t.Run("NodePropertiesAndContent", func(t *testing.T) {
		g := createSampleGrammar()

		// Test properties and content of specific nodes
		verifyNodePropertiesAndContent(t, g)
	})
}

func createSampleGrammar() *schemas.Grammar {
	g := schemas.NewGrammar()
	for i := 0; i < 10; i++ {
		father := schemas.NewNode(g, schemas.GrammarProduction, fmt.Sprintf("node%d", i), "test")
		for j := 0; j < 3; j++ {
			child := schemas.NewNode(g, schemas.GrammarProduction, fmt.Sprintf("child%d-%d", i, j), "test")
			father.AddSymbol(child)
			child.SetRoot(father)
		}
	}
	return g
}

func verifyGrammarNodes(t *testing.T, g *schemas.Grammar) {
	// Verify that grammar nodes are created correctly
	for i := 0; i < 10; i++ {
		nodeID := fmt.Sprintf("node%d", i)
		if g != g.GetNode(nodeID).GetGrammar() {
			t.Errorf("Grammar of node %s does not match", nodeID)
		}
	}
}

func verifyNodePropertiesAndContent(t *testing.T, g *schemas.Grammar) {
	// Verify properties and content of specific nodes
	node2Symbols := g.GetNode("node2").GetSymbols()
	if len(node2Symbols) != 3 {
		t.Error("Length does not match")
	}
	for j, v := range node2Symbols {
		expectedID := fmt.Sprintf("child2-%d", j)
		if v.GetID() != expectedID {
			t.Logf("%s does not equal to %s, it is ok\n", v.GetID(), expectedID)
		}
	}

	node3 := g.GetNode("node3")
	id := node3.GetSymbol(2).GetID()
	if id != "child3-0" && id != "child3-1" && id != "child3-2" {
		t.Errorf("%s should equal to one of child3-0, child3-1, child3-2\n", id)
	}
	if node3.GetContent() != "test" {
		t.Errorf("%s should equal to test \n", node3.GetContent())
	}

	if g.GetNode("node2").GetType() != schemas.GrammarProduction {
		t.Errorf("%d should equal to schemas.GrammarProduction \n", g.GetNode("node2").GetType())
	}

	if g.GetNode("node2").GetSymbol(10) != nil {
		t.Errorf("node2's 10 child should equal to nil \n")
	}
}
