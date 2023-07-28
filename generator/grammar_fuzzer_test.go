package Generator

import (
	"github.com/CUHK-SE-Group/ebnf-based-generator/grammar"
	"testing"
)

func TestGrammarFuzzer(t *testing.T) {
	f := GrammarFuzzer{Grammar: grammar.SimpleNonterminalGrammar, StartSymbol: grammar.StartSymbol}
	f.Grammar.Visualize("test.png")
}
