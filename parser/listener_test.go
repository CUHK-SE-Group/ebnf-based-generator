package parser

import (
	"fmt"
	"github.com/CUHK-SE-Group/generic-generator/graph"
	"github.com/CUHK-SE-Group/generic-generator/schemas"
	"testing"
)

func TestParseCalculator(t *testing.T) {
	productions, _ := Parse("./testdata/simple.ebnf")
	graph.Visualize(productions.GetInternal(), "file.dot", func(v graph.Vertex[schemas.Property]) string {
		return fmt.Sprintf("id: %s content: %s type: %d", v.GetID(), v.GetProperty(schemas.Prop).Content, v.GetProperty(schemas.Prop).Type)
	})
	//graph.Visualize(productions.GetInternal(), "file.dot", nil, nil)
}

//func TestParseCypher(t *testing.T) {
//	productions, _ := Parse("./testdata/cypher.ebnf")
//	for _, g := range productions {
//		fmt.Printf("%+v\n", g)
//	}
//	root := MergeProduction(productions, "Cypher")
//	root.Visualize(fmt.Sprintf("figures/root.gv"), true)
//}
