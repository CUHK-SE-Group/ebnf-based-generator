// Code generated from ./PathQuery.g4 by ANTLR 4.13.0. DO NOT EDIT.

package pathQuery

import (
	"fmt"
	"github.com/antlr4-go/antlr/v4"
	"sync"
	"unicode"
)

// Suppress unused import error
var _ = fmt.Printf
var _ = sync.Once{}
var _ = unicode.IsLetter

type PathQueryLexer struct {
	*antlr.BaseLexer
	channelNames []string
	modeNames    []string
	// TODO: EOF string
}

var PathQueryLexerLexerStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	ChannelNames           []string
	ModeNames              []string
	LiteralNames           []string
	SymbolicNames          []string
	RuleNames              []string
	PredictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func pathquerylexerLexerInit() {
	staticData := &PathQueryLexerLexerStaticData
	staticData.ChannelNames = []string{
		"DEFAULT_TOKEN_CHANNEL", "HIDDEN",
	}
	staticData.ModeNames = []string{
		"DEFAULT_MODE",
	}
	staticData.LiteralNames = []string{
		"", "'*'", "'/'", "'//'",
	}
	staticData.SymbolicNames = []string{
		"", "", "", "", "NODE_ID", "WS",
	}
	staticData.RuleNames = []string{
		"T__0", "T__1", "T__2", "NODE_ID", "WS",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 0, 5, 30, 6, -1, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2,
		4, 7, 4, 1, 0, 1, 0, 1, 1, 1, 1, 1, 2, 1, 2, 1, 2, 1, 3, 4, 3, 20, 8, 3,
		11, 3, 12, 3, 21, 1, 4, 4, 4, 25, 8, 4, 11, 4, 12, 4, 26, 1, 4, 1, 4, 0,
		0, 5, 1, 1, 3, 2, 5, 3, 7, 4, 9, 5, 1, 0, 2, 5, 0, 35, 35, 48, 57, 65,
		90, 95, 95, 97, 122, 3, 0, 9, 10, 13, 13, 32, 32, 31, 0, 1, 1, 0, 0, 0,
		0, 3, 1, 0, 0, 0, 0, 5, 1, 0, 0, 0, 0, 7, 1, 0, 0, 0, 0, 9, 1, 0, 0, 0,
		1, 11, 1, 0, 0, 0, 3, 13, 1, 0, 0, 0, 5, 15, 1, 0, 0, 0, 7, 19, 1, 0, 0,
		0, 9, 24, 1, 0, 0, 0, 11, 12, 5, 42, 0, 0, 12, 2, 1, 0, 0, 0, 13, 14, 5,
		47, 0, 0, 14, 4, 1, 0, 0, 0, 15, 16, 5, 47, 0, 0, 16, 17, 5, 47, 0, 0,
		17, 6, 1, 0, 0, 0, 18, 20, 7, 0, 0, 0, 19, 18, 1, 0, 0, 0, 20, 21, 1, 0,
		0, 0, 21, 19, 1, 0, 0, 0, 21, 22, 1, 0, 0, 0, 22, 8, 1, 0, 0, 0, 23, 25,
		7, 1, 0, 0, 24, 23, 1, 0, 0, 0, 25, 26, 1, 0, 0, 0, 26, 24, 1, 0, 0, 0,
		26, 27, 1, 0, 0, 0, 27, 28, 1, 0, 0, 0, 28, 29, 6, 4, 0, 0, 29, 10, 1,
		0, 0, 0, 3, 0, 21, 26, 1, 6, 0, 0,
	}
	deserializer := antlr.NewATNDeserializer(nil)
	staticData.atn = deserializer.Deserialize(staticData.serializedATN)
	atn := staticData.atn
	staticData.decisionToDFA = make([]*antlr.DFA, len(atn.DecisionToState))
	decisionToDFA := staticData.decisionToDFA
	for index, state := range atn.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(state, index)
	}
}

// PathQueryLexerInit initializes any static state used to implement PathQueryLexer. By default the
// static state used to implement the lexer is lazily initialized during the first call to
// NewPathQueryLexer(). You can call this function if you wish to initialize the static state ahead
// of time.
func PathQueryLexerInit() {
	staticData := &PathQueryLexerLexerStaticData
	staticData.once.Do(pathquerylexerLexerInit)
}

// NewPathQueryLexer produces a new lexer instance for the optional input antlr.CharStream.
func NewPathQueryLexer(input antlr.CharStream) *PathQueryLexer {
	PathQueryLexerInit()
	l := new(PathQueryLexer)
	l.BaseLexer = antlr.NewBaseLexer(input)
	staticData := &PathQueryLexerLexerStaticData
	l.Interpreter = antlr.NewLexerATNSimulator(l, staticData.atn, staticData.decisionToDFA, staticData.PredictionContextCache)
	l.channelNames = staticData.ChannelNames
	l.modeNames = staticData.ModeNames
	l.RuleNames = staticData.RuleNames
	l.LiteralNames = staticData.LiteralNames
	l.SymbolicNames = staticData.SymbolicNames
	l.GrammarFileName = "PathQuery.g4"
	// TODO: l.EOF = antlr.TokenEOF

	return l
}

// PathQueryLexer tokens.
const (
	PathQueryLexerT__0    = 1
	PathQueryLexerT__1    = 2
	PathQueryLexerT__2    = 3
	PathQueryLexerNODE_ID = 4
	PathQueryLexerWS      = 5
)
