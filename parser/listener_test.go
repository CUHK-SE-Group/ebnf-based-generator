package parser

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/CUHK-SE-Group/generic-generator/graph"
	"github.com/CUHK-SE-Group/generic-generator/schemas"
)

func escapeQuotes(s string) string {
	// 使用 strings.Replace 替换所有的双引号
	return strings.ReplaceAll(s, "\"", "\\\"")
}
func parseAndVisualize(file string) {
	productions, _ := Parse(file, "program")
	filenameWithoutExt := file[:len(file)-len(filepath.Ext(file))]
	productions.MergeProduction()
	graph.Visualize(productions.GetInternal(), filepath.Join(filenameWithoutExt+".dot"), func(v graph.Vertex[schemas.Property]) string {
		return fmt.Sprintf("id: %s\n content: %s\n type: %s", v.GetID(), escapeQuotes(v.GetProperty(schemas.Prop).Content), schemas.GetGrammarTypeStr(v.GetProperty(schemas.Prop).Type))
	})
}
func TestBasicBuildPath(t *testing.T) {
	p, _ := Parse("./testdata/complete/tinyc.ebnf", "")
	tree := p.GetIndex("program")
	fmt.Println(tree.ToJSON())
}

func TestBasic(t *testing.T) {
	parseAndVisualize("./testdata/basic/basic_all.ebnf")
	parseAndVisualize("./testdata/basic/basic_comma.ebnf")
	parseAndVisualize("./testdata/basic/basic_or.ebnf")
}

func TestNested(t *testing.T) {
	parseAndVisualize("./testdata/nested/nested_paren.ebnf")
	parseAndVisualize("./testdata/nested/nested_brace.ebnf")
	parseAndVisualize("./testdata/nested/nested_bracket.ebnf")
	parseAndVisualize("./testdata/nested/nested_all.ebnf")
}

func TestChoice(t *testing.T) {
	parseAndVisualize("./testdata/choice/choice.ebnf")
}

func TestStrings(t *testing.T) {
	parseAndVisualize("./testdata/strings/single_quote.ebnf")
	parseAndVisualize("./testdata/strings/double_quote.ebnf")
}

func TestParseCalculator(t *testing.T) {
	parseAndVisualize("./testdata/complete/simple.ebnf")
}

func TestParseCypher(t *testing.T) {
	parseAndVisualize("./testdata/complete/cypher.ebnf")
}
func TestParseTinyC(t *testing.T) {
	parseAndVisualize("./testdata/complete/tinyc.ebnf")
}
