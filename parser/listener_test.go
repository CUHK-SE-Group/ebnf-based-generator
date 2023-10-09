package parser

import (
	"fmt"
	"testing"
)

func TestParseCalculator(t *testing.T) {
	productions, _ := Parse("./testdata/simple.ebnf")
	for _, g := range productions {
		fmt.Printf("%+v\n", g)
	}
	root := MergeProduction(productions, "expression")
	root.Visualize(fmt.Sprintf("figures/root.gv"), true)
}

func TestParseCypher(t *testing.T) {
	productions, _ := Parse("./testdata/cypher.ebnf")
	for _, g := range productions {
		fmt.Printf("%+v\n", g)
	}
	root := MergeProduction(productions, "Cypher")
	root.Visualize(fmt.Sprintf("figures/root.gv"), true)
}
