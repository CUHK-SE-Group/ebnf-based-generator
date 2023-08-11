package parser

import (
	"testing"

	"github.com/CUHK-SE-Group/ebnf-based-generator/parser/ebnf"

	"github.com/antlr4-go/antlr/v4"
)

func TestGeneration(t *testing.T) {
	is, err := antlr.NewFileStream("./testdata/generate.ebnf")
	if err != nil {
		t.Fatalf("Can not open testdata")
	}
	Init()
	lexer := ebnf.NewEBNFLexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	parser := ebnf.NewEBNFParser(stream)
	listener := newEbnfListener()
	antlr.ParseTreeWalkerDefault.Walk(listener, parser.Ebnf())
	g := listener.productions["genSimple"]
	g.Visualize("./testdata/generate.output.png", true)
	r := g.Generate(nil)
	for _, s := range r.output {
		t.Logf("%s ", s)
	}
	r.Visualize("./testdata/generate.path.output.png")
}
