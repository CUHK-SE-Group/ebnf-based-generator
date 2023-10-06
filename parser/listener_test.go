package parser

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	productions, _ := Parse("./testdata/simple.ebnf")
	for _, g := range productions {
		fmt.Printf("%+v\n", g)
	}
}
func TestParse1(t *testing.T) {
	productions, _ := Parse("./testdata/cypher.ebnf")
	for _, g := range productions {
		fmt.Printf("%+v\n", g)
	}
}
