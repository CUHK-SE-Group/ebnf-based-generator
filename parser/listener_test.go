package parser

import (
	"fmt"
	"testing"

	"github.com/CUHK-SE-Group/ebnf-based-generator/parser/ebnf"

	"github.com/antlr4-go/antlr/v4"
)

func TestParse(t *testing.T) {
	Init()
	is, err := antlr.NewFileStream("./testdata/parse.ebnf")
	if err != nil {
		t.Fatalf("Can not open testdata")
	}
	lexer := ebnf.NewEBNFLexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	parser := ebnf.NewEBNFParser(stream)
	listener := newEbnfListener()
	antlr.ParseTreeWalkerDefault.Walk(listener, parser.Ebnf())
	for k, g := range listener.productions {
		if k == "test_production" || k == "like" || k == "parens" || k == "HexLetter" {
			path := fmt.Sprintf("./testdata/parse.%s.output.png", k)
			g.Visualize(path, false)
		}
	}
}
