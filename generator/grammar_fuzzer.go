package Generator

import "github.com/CUHK-SE-Group/ebnf-based-generator/grammar"

type Config struct {
	MaxRecursion    int
	MaxDepth        int
	MaxAttempts     int
	MinNonTerminals int
	MaxNonTerminals int
}
type GrammarFuzzer struct {
	Grammar     grammar.Grammar
	StartSymbol string
	Configs     Config
}
