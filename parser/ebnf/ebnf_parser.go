// Code generated from ./EBNFParser.g4 by ANTLR 4.13.0. DO NOT EDIT.

package ebnf // EBNFParser
import (
	"fmt"
	"strconv"
	"sync"

	"github.com/antlr4-go/antlr/v4"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = strconv.Itoa
var _ = sync.Once{}

type EBNFParser struct {
	*antlr.BaseParser
}

var EBNFParserParserStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	LiteralNames           []string
	SymbolicNames          []string
	RuleNames              []string
	PredictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func ebnfparserParserInit() {
	staticData := &EBNFParserParserStaticData
	staticData.LiteralNames = []string{
		"", "", "'('", "')'", "'['", "']'", "'{'", "'}'", "';'", "", "'|'",
		"'-'", "'*'", "'+'", "'?'", "','",
	}
	staticData.SymbolicNames = []string{
		"", "LINE_COMMENT", "LPAREN", "RPAREN", "LBRACKET", "RBRACKET", "LBRACE",
		"RBRACE", "SEMICOLON", "EQUAL", "OR", "SUB", "REP", "PLUS", "EXT", "COMMA",
		"ID", "WHITESPACE", "QUOTE", "DOUBLEQUOTE", "TEXT", "REGTEXT",
	}
	staticData.RuleNames = []string{
		"ebnf", "production", "expr", "term", "factor", "choice", "identifier",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 21, 82, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
		4, 2, 5, 7, 5, 2, 6, 7, 6, 1, 0, 5, 0, 16, 8, 0, 10, 0, 12, 0, 19, 9, 0,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 1, 2, 1, 2, 5, 2, 29, 8, 2, 10, 2,
		12, 2, 32, 9, 2, 1, 3, 1, 3, 1, 3, 5, 3, 37, 8, 3, 10, 3, 12, 3, 40, 9,
		3, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 1,
		4, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 3, 4, 62, 8, 4, 1, 4,
		1, 4, 1, 4, 3, 4, 67, 8, 4, 5, 4, 69, 8, 4, 10, 4, 12, 4, 72, 9, 4, 1,
		5, 1, 5, 1, 5, 1, 5, 3, 5, 78, 8, 5, 1, 6, 1, 6, 1, 6, 0, 1, 8, 7, 0, 2,
		4, 6, 8, 10, 12, 0, 0, 87, 0, 17, 1, 0, 0, 0, 2, 20, 1, 0, 0, 0, 4, 25,
		1, 0, 0, 0, 6, 33, 1, 0, 0, 0, 8, 61, 1, 0, 0, 0, 10, 77, 1, 0, 0, 0, 12,
		79, 1, 0, 0, 0, 14, 16, 3, 2, 1, 0, 15, 14, 1, 0, 0, 0, 16, 19, 1, 0, 0,
		0, 17, 15, 1, 0, 0, 0, 17, 18, 1, 0, 0, 0, 18, 1, 1, 0, 0, 0, 19, 17, 1,
		0, 0, 0, 20, 21, 5, 16, 0, 0, 21, 22, 5, 9, 0, 0, 22, 23, 3, 4, 2, 0, 23,
		24, 5, 8, 0, 0, 24, 3, 1, 0, 0, 0, 25, 30, 3, 6, 3, 0, 26, 27, 5, 15, 0,
		0, 27, 29, 3, 6, 3, 0, 28, 26, 1, 0, 0, 0, 29, 32, 1, 0, 0, 0, 30, 28,
		1, 0, 0, 0, 30, 31, 1, 0, 0, 0, 31, 5, 1, 0, 0, 0, 32, 30, 1, 0, 0, 0,
		33, 38, 3, 8, 4, 0, 34, 35, 5, 10, 0, 0, 35, 37, 3, 8, 4, 0, 36, 34, 1,
		0, 0, 0, 37, 40, 1, 0, 0, 0, 38, 36, 1, 0, 0, 0, 38, 39, 1, 0, 0, 0, 39,
		7, 1, 0, 0, 0, 40, 38, 1, 0, 0, 0, 41, 42, 6, 4, -1, 0, 42, 62, 3, 12,
		6, 0, 43, 44, 5, 4, 0, 0, 44, 45, 3, 4, 2, 0, 45, 46, 5, 5, 0, 0, 46, 62,
		1, 0, 0, 0, 47, 48, 5, 6, 0, 0, 48, 49, 3, 4, 2, 0, 49, 50, 5, 7, 0, 0,
		50, 62, 1, 0, 0, 0, 51, 52, 5, 18, 0, 0, 52, 53, 5, 20, 0, 0, 53, 62, 5,
		18, 0, 0, 54, 55, 5, 19, 0, 0, 55, 56, 5, 21, 0, 0, 56, 62, 5, 19, 0, 0,
		57, 58, 5, 2, 0, 0, 58, 59, 3, 4, 2, 0, 59, 60, 5, 3, 0, 0, 60, 62, 1,
		0, 0, 0, 61, 41, 1, 0, 0, 0, 61, 43, 1, 0, 0, 0, 61, 47, 1, 0, 0, 0, 61,
		51, 1, 0, 0, 0, 61, 54, 1, 0, 0, 0, 61, 57, 1, 0, 0, 0, 62, 70, 1, 0, 0,
		0, 63, 64, 10, 6, 0, 0, 64, 66, 3, 10, 5, 0, 65, 67, 3, 8, 4, 0, 66, 65,
		1, 0, 0, 0, 66, 67, 1, 0, 0, 0, 67, 69, 1, 0, 0, 0, 68, 63, 1, 0, 0, 0,
		69, 72, 1, 0, 0, 0, 70, 68, 1, 0, 0, 0, 70, 71, 1, 0, 0, 0, 71, 9, 1, 0,
		0, 0, 72, 70, 1, 0, 0, 0, 73, 78, 5, 12, 0, 0, 74, 78, 5, 13, 0, 0, 75,
		78, 5, 14, 0, 0, 76, 78, 5, 11, 0, 0, 77, 73, 1, 0, 0, 0, 77, 74, 1, 0,
		0, 0, 77, 75, 1, 0, 0, 0, 77, 76, 1, 0, 0, 0, 78, 11, 1, 0, 0, 0, 79, 80,
		5, 16, 0, 0, 80, 13, 1, 0, 0, 0, 7, 17, 30, 38, 61, 66, 70, 77,
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

// EBNFParserInit initializes any static state used to implement EBNFParser. By default the
// static state used to implement the parser is lazily initialized during the first call to
// NewEBNFParser(). You can call this function if you wish to initialize the static state ahead
// of time.
func EBNFParserInit() {
	staticData := &EBNFParserParserStaticData
	staticData.once.Do(ebnfparserParserInit)
}

// NewEBNFParser produces a new parser instance for the optional input antlr.TokenStream.
func NewEBNFParser(input antlr.TokenStream) *EBNFParser {
	EBNFParserInit()
	this := new(EBNFParser)
	this.BaseParser = antlr.NewBaseParser(input)
	staticData := &EBNFParserParserStaticData
	this.Interpreter = antlr.NewParserATNSimulator(this, staticData.atn, staticData.decisionToDFA, staticData.PredictionContextCache)
	this.RuleNames = staticData.RuleNames
	this.LiteralNames = staticData.LiteralNames
	this.SymbolicNames = staticData.SymbolicNames
	this.GrammarFileName = "EBNFParser.g4"

	return this
}

// EBNFParser tokens.
const (
	EBNFParserEOF          = antlr.TokenEOF
	EBNFParserLINE_COMMENT = 1
	EBNFParserLPAREN       = 2
	EBNFParserRPAREN       = 3
	EBNFParserLBRACKET     = 4
	EBNFParserRBRACKET     = 5
	EBNFParserLBRACE       = 6
	EBNFParserRBRACE       = 7
	EBNFParserSEMICOLON    = 8
	EBNFParserEQUAL        = 9
	EBNFParserOR           = 10
	EBNFParserSUB          = 11
	EBNFParserREP          = 12
	EBNFParserPLUS         = 13
	EBNFParserEXT          = 14
	EBNFParserCOMMA        = 15
	EBNFParserID           = 16
	EBNFParserWHITESPACE   = 17
	EBNFParserQUOTE        = 18
	EBNFParserDOUBLEQUOTE  = 19
	EBNFParserTEXT         = 20
	EBNFParserREGTEXT      = 21
)

// EBNFParser rules.
const (
	EBNFParserRULE_ebnf       = 0
	EBNFParserRULE_production = 1
	EBNFParserRULE_expr       = 2
	EBNFParserRULE_term       = 3
	EBNFParserRULE_factor     = 4
	EBNFParserRULE_choice     = 5
	EBNFParserRULE_identifier = 6
)

// IEbnfContext is an interface to support dynamic dispatch.
type IEbnfContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllProduction() []IProductionContext
	Production(i int) IProductionContext

	// IsEbnfContext differentiates from other interfaces.
	IsEbnfContext()
}

type EbnfContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEbnfContext() *EbnfContext {
	var p = new(EbnfContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EBNFParserRULE_ebnf
	return p
}

func InitEmptyEbnfContext(p *EbnfContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EBNFParserRULE_ebnf
}

func (*EbnfContext) IsEbnfContext() {}

func NewEbnfContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EbnfContext {
	var p = new(EbnfContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EBNFParserRULE_ebnf

	return p
}

func (s *EbnfContext) GetParser() antlr.Parser { return s.parser }

func (s *EbnfContext) AllProduction() []IProductionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IProductionContext); ok {
			len++
		}
	}

	tst := make([]IProductionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IProductionContext); ok {
			tst[i] = t.(IProductionContext)
			i++
		}
	}

	return tst
}

func (s *EbnfContext) Production(i int) IProductionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IProductionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IProductionContext)
}

func (s *EbnfContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EbnfContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EbnfContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EBNFParserListener); ok {
		listenerT.EnterEbnf(s)
	}
}

func (s *EbnfContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EBNFParserListener); ok {
		listenerT.ExitEbnf(s)
	}
}

func (p *EBNFParser) Ebnf() (localctx IEbnfContext) {
	localctx = NewEbnfContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, EBNFParserRULE_ebnf)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(17)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == EBNFParserID {
		{
			p.SetState(14)
			p.Production()
		}

		p.SetState(19)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IProductionContext is an interface to support dynamic dispatch.
type IProductionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ID() antlr.TerminalNode
	EQUAL() antlr.TerminalNode
	Expr() IExprContext
	SEMICOLON() antlr.TerminalNode

	// IsProductionContext differentiates from other interfaces.
	IsProductionContext()
}

type ProductionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyProductionContext() *ProductionContext {
	var p = new(ProductionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EBNFParserRULE_production
	return p
}

func InitEmptyProductionContext(p *ProductionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EBNFParserRULE_production
}

func (*ProductionContext) IsProductionContext() {}

func NewProductionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ProductionContext {
	var p = new(ProductionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EBNFParserRULE_production

	return p
}

func (s *ProductionContext) GetParser() antlr.Parser { return s.parser }

func (s *ProductionContext) ID() antlr.TerminalNode {
	return s.GetToken(EBNFParserID, 0)
}

func (s *ProductionContext) EQUAL() antlr.TerminalNode {
	return s.GetToken(EBNFParserEQUAL, 0)
}

func (s *ProductionContext) Expr() IExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *ProductionContext) SEMICOLON() antlr.TerminalNode {
	return s.GetToken(EBNFParserSEMICOLON, 0)
}

func (s *ProductionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ProductionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ProductionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EBNFParserListener); ok {
		listenerT.EnterProduction(s)
	}
}

func (s *ProductionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EBNFParserListener); ok {
		listenerT.ExitProduction(s)
	}
}

func (p *EBNFParser) Production() (localctx IProductionContext) {
	localctx = NewProductionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, EBNFParserRULE_production)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(20)
		p.Match(EBNFParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(21)
		p.Match(EBNFParserEQUAL)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(22)
		p.Expr()
	}
	{
		p.SetState(23)
		p.Match(EBNFParserSEMICOLON)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IExprContext is an interface to support dynamic dispatch.
type IExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllTerm() []ITermContext
	Term(i int) ITermContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsExprContext differentiates from other interfaces.
	IsExprContext()
}

type ExprContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyExprContext() *ExprContext {
	var p = new(ExprContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EBNFParserRULE_expr
	return p
}

func InitEmptyExprContext(p *ExprContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EBNFParserRULE_expr
}

func (*ExprContext) IsExprContext() {}

func NewExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExprContext {
	var p = new(ExprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EBNFParserRULE_expr

	return p
}

func (s *ExprContext) GetParser() antlr.Parser { return s.parser }

func (s *ExprContext) AllTerm() []ITermContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ITermContext); ok {
			len++
		}
	}

	tst := make([]ITermContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ITermContext); ok {
			tst[i] = t.(ITermContext)
			i++
		}
	}

	return tst
}

func (s *ExprContext) Term(i int) ITermContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITermContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITermContext)
}

func (s *ExprContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(EBNFParserCOMMA)
}

func (s *ExprContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(EBNFParserCOMMA, i)
}

func (s *ExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EBNFParserListener); ok {
		listenerT.EnterExpr(s)
	}
}

func (s *ExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EBNFParserListener); ok {
		listenerT.ExitExpr(s)
	}
}

func (p *EBNFParser) Expr() (localctx IExprContext) {
	localctx = NewExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, EBNFParserRULE_expr)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(25)
		p.Term()
	}
	p.SetState(30)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == EBNFParserCOMMA {
		{
			p.SetState(26)
			p.Match(EBNFParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(27)
			p.Term()
		}

		p.SetState(32)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ITermContext is an interface to support dynamic dispatch.
type ITermContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllFactor() []IFactorContext
	Factor(i int) IFactorContext
	AllOR() []antlr.TerminalNode
	OR(i int) antlr.TerminalNode

	// IsTermContext differentiates from other interfaces.
	IsTermContext()
}

type TermContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTermContext() *TermContext {
	var p = new(TermContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EBNFParserRULE_term
	return p
}

func InitEmptyTermContext(p *TermContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EBNFParserRULE_term
}

func (*TermContext) IsTermContext() {}

func NewTermContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TermContext {
	var p = new(TermContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EBNFParserRULE_term

	return p
}

func (s *TermContext) GetParser() antlr.Parser { return s.parser }

func (s *TermContext) AllFactor() []IFactorContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IFactorContext); ok {
			len++
		}
	}

	tst := make([]IFactorContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IFactorContext); ok {
			tst[i] = t.(IFactorContext)
			i++
		}
	}

	return tst
}

func (s *TermContext) Factor(i int) IFactorContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFactorContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFactorContext)
}

func (s *TermContext) AllOR() []antlr.TerminalNode {
	return s.GetTokens(EBNFParserOR)
}

func (s *TermContext) OR(i int) antlr.TerminalNode {
	return s.GetToken(EBNFParserOR, i)
}

func (s *TermContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TermContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TermContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EBNFParserListener); ok {
		listenerT.EnterTerm(s)
	}
}

func (s *TermContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EBNFParserListener); ok {
		listenerT.ExitTerm(s)
	}
}

func (p *EBNFParser) Term() (localctx ITermContext) {
	localctx = NewTermContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, EBNFParserRULE_term)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(33)
		p.factor(0)
	}
	p.SetState(38)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == EBNFParserOR {
		{
			p.SetState(34)
			p.Match(EBNFParserOR)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(35)
			p.factor(0)
		}

		p.SetState(40)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IFactorContext is an interface to support dynamic dispatch.
type IFactorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsFactorContext differentiates from other interfaces.
	IsFactorContext()
}

type FactorContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFactorContext() *FactorContext {
	var p = new(FactorContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EBNFParserRULE_factor
	return p
}

func InitEmptyFactorContext(p *FactorContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EBNFParserRULE_factor
}

func (*FactorContext) IsFactorContext() {}

func NewFactorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FactorContext {
	var p = new(FactorContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EBNFParserRULE_factor

	return p
}

func (s *FactorContext) GetParser() antlr.Parser { return s.parser }

func (s *FactorContext) CopyAll(ctx *FactorContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *FactorContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FactorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type QUOTEContext struct {
	FactorContext
}

func NewQUOTEContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *QUOTEContext {
	var p = new(QUOTEContext)

	InitEmptyFactorContext(&p.FactorContext)
	p.parser = parser
	p.CopyAll(ctx.(*FactorContext))

	return p
}

func (s *QUOTEContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *QUOTEContext) AllQUOTE() []antlr.TerminalNode {
	return s.GetTokens(EBNFParserQUOTE)
}

func (s *QUOTEContext) QUOTE(i int) antlr.TerminalNode {
	return s.GetToken(EBNFParserQUOTE, i)
}

func (s *QUOTEContext) TEXT() antlr.TerminalNode {
	return s.GetToken(EBNFParserTEXT, 0)
}

func (s *QUOTEContext) AllDOUBLEQUOTE() []antlr.TerminalNode {
	return s.GetTokens(EBNFParserDOUBLEQUOTE)
}

func (s *QUOTEContext) DOUBLEQUOTE(i int) antlr.TerminalNode {
	return s.GetToken(EBNFParserDOUBLEQUOTE, i)
}

func (s *QUOTEContext) REGTEXT() antlr.TerminalNode {
	return s.GetToken(EBNFParserREGTEXT, 0)
}

func (s *QUOTEContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EBNFParserListener); ok {
		listenerT.EnterQUOTE(s)
	}
}

func (s *QUOTEContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EBNFParserListener); ok {
		listenerT.ExitQUOTE(s)
	}
}

type CHOICEContext struct {
	FactorContext
}

func NewCHOICEContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *CHOICEContext {
	var p = new(CHOICEContext)

	InitEmptyFactorContext(&p.FactorContext)
	p.parser = parser
	p.CopyAll(ctx.(*FactorContext))

	return p
}

func (s *CHOICEContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CHOICEContext) AllFactor() []IFactorContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IFactorContext); ok {
			len++
		}
	}

	tst := make([]IFactorContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IFactorContext); ok {
			tst[i] = t.(IFactorContext)
			i++
		}
	}

	return tst
}

func (s *CHOICEContext) Factor(i int) IFactorContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFactorContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFactorContext)
}

func (s *CHOICEContext) Choice() IChoiceContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IChoiceContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IChoiceContext)
}

func (s *CHOICEContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EBNFParserListener); ok {
		listenerT.EnterCHOICE(s)
	}
}

func (s *CHOICEContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EBNFParserListener); ok {
		listenerT.ExitCHOICE(s)
	}
}

type BRACKETContext struct {
	FactorContext
}

func NewBRACKETContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *BRACKETContext {
	var p = new(BRACKETContext)

	InitEmptyFactorContext(&p.FactorContext)
	p.parser = parser
	p.CopyAll(ctx.(*FactorContext))

	return p
}

func (s *BRACKETContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BRACKETContext) LBRACKET() antlr.TerminalNode {
	return s.GetToken(EBNFParserLBRACKET, 0)
}

func (s *BRACKETContext) Expr() IExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *BRACKETContext) RBRACKET() antlr.TerminalNode {
	return s.GetToken(EBNFParserRBRACKET, 0)
}

func (s *BRACKETContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EBNFParserListener); ok {
		listenerT.EnterBRACKET(s)
	}
}

func (s *BRACKETContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EBNFParserListener); ok {
		listenerT.ExitBRACKET(s)
	}
}

type IDContext struct {
	FactorContext
}

func NewIDContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *IDContext {
	var p = new(IDContext)

	InitEmptyFactorContext(&p.FactorContext)
	p.parser = parser
	p.CopyAll(ctx.(*FactorContext))

	return p
}

func (s *IDContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IDContext) Identifier() IIdentifierContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentifierContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentifierContext)
}

func (s *IDContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EBNFParserListener); ok {
		listenerT.EnterID(s)
	}
}

func (s *IDContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EBNFParserListener); ok {
		listenerT.ExitID(s)
	}
}

type BRACEContext struct {
	FactorContext
}

func NewBRACEContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *BRACEContext {
	var p = new(BRACEContext)

	InitEmptyFactorContext(&p.FactorContext)
	p.parser = parser
	p.CopyAll(ctx.(*FactorContext))

	return p
}

func (s *BRACEContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BRACEContext) LBRACE() antlr.TerminalNode {
	return s.GetToken(EBNFParserLBRACE, 0)
}

func (s *BRACEContext) Expr() IExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *BRACEContext) RBRACE() antlr.TerminalNode {
	return s.GetToken(EBNFParserRBRACE, 0)
}

func (s *BRACEContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EBNFParserListener); ok {
		listenerT.EnterBRACE(s)
	}
}

func (s *BRACEContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EBNFParserListener); ok {
		listenerT.ExitBRACE(s)
	}
}

type NoneContext struct {
	FactorContext
}

func NewNoneContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *NoneContext {
	var p = new(NoneContext)

	InitEmptyFactorContext(&p.FactorContext)
	p.parser = parser
	p.CopyAll(ctx.(*FactorContext))

	return p
}

func (s *NoneContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NoneContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(EBNFParserLPAREN, 0)
}

func (s *NoneContext) Expr() IExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *NoneContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(EBNFParserRPAREN, 0)
}

func (s *NoneContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EBNFParserListener); ok {
		listenerT.EnterNone(s)
	}
}

func (s *NoneContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EBNFParserListener); ok {
		listenerT.ExitNone(s)
	}
}

func (p *EBNFParser) Factor() (localctx IFactorContext) {
	return p.factor(0)
}

func (p *EBNFParser) factor(_p int) (localctx IFactorContext) {
	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()

	_parentState := p.GetState()
	localctx = NewFactorContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx IFactorContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 8
	p.EnterRecursionRule(localctx, 8, EBNFParserRULE_factor, _p)
	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(61)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case EBNFParserID:
		localctx = NewIDContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx

		{
			p.SetState(42)
			p.Identifier()
		}

	case EBNFParserLBRACKET:
		localctx = NewBRACKETContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(43)
			p.Match(EBNFParserLBRACKET)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(44)
			p.Expr()
		}
		{
			p.SetState(45)
			p.Match(EBNFParserRBRACKET)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case EBNFParserLBRACE:
		localctx = NewBRACEContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(47)
			p.Match(EBNFParserLBRACE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(48)
			p.Expr()
		}
		{
			p.SetState(49)
			p.Match(EBNFParserRBRACE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case EBNFParserQUOTE:
		localctx = NewQUOTEContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(51)
			p.Match(EBNFParserQUOTE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(52)
			p.Match(EBNFParserTEXT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(53)
			p.Match(EBNFParserQUOTE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case EBNFParserDOUBLEQUOTE:
		localctx = NewQUOTEContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(54)
			p.Match(EBNFParserDOUBLEQUOTE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(55)
			p.Match(EBNFParserREGTEXT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(56)
			p.Match(EBNFParserDOUBLEQUOTE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case EBNFParserLPAREN:
		localctx = NewNoneContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(57)
			p.Match(EBNFParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(58)
			p.Expr()
		}
		{
			p.SetState(59)
			p.Match(EBNFParserRPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(70)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 5, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			localctx = NewCHOICEContext(p, NewFactorContext(p, _parentctx, _parentState))
			p.PushNewRecursionContext(localctx, _startState, EBNFParserRULE_factor)
			p.SetState(63)

			if !(p.Precpred(p.GetParserRuleContext(), 6)) {
				p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 6)", ""))
				goto errorExit
			}
			{
				p.SetState(64)
				p.Choice()
			}
			p.SetState(66)
			p.GetErrorHandler().Sync(p)

			if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 4, p.GetParserRuleContext()) == 1 {
				{
					p.SetState(65)
					p.factor(0)
				}

			} else if p.HasError() { // JIM
				goto errorExit
			}

		}
		p.SetState(72)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 5, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.UnrollRecursionContexts(_parentctx)
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IChoiceContext is an interface to support dynamic dispatch.
type IChoiceContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsChoiceContext differentiates from other interfaces.
	IsChoiceContext()
}

type ChoiceContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyChoiceContext() *ChoiceContext {
	var p = new(ChoiceContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EBNFParserRULE_choice
	return p
}

func InitEmptyChoiceContext(p *ChoiceContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EBNFParserRULE_choice
}

func (*ChoiceContext) IsChoiceContext() {}

func NewChoiceContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ChoiceContext {
	var p = new(ChoiceContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EBNFParserRULE_choice

	return p
}

func (s *ChoiceContext) GetParser() antlr.Parser { return s.parser }

func (s *ChoiceContext) CopyAll(ctx *ChoiceContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *ChoiceContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ChoiceContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type EXTContext struct {
	ChoiceContext
}

func NewEXTContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *EXTContext {
	var p = new(EXTContext)

	InitEmptyChoiceContext(&p.ChoiceContext)
	p.parser = parser
	p.CopyAll(ctx.(*ChoiceContext))

	return p
}

func (s *EXTContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EXTContext) EXT() antlr.TerminalNode {
	return s.GetToken(EBNFParserEXT, 0)
}

func (s *EXTContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EBNFParserListener); ok {
		listenerT.EnterEXT(s)
	}
}

func (s *EXTContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EBNFParserListener); ok {
		listenerT.ExitEXT(s)
	}
}

type SUBContext struct {
	ChoiceContext
}

func NewSUBContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SUBContext {
	var p = new(SUBContext)

	InitEmptyChoiceContext(&p.ChoiceContext)
	p.parser = parser
	p.CopyAll(ctx.(*ChoiceContext))

	return p
}

func (s *SUBContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SUBContext) SUB() antlr.TerminalNode {
	return s.GetToken(EBNFParserSUB, 0)
}

func (s *SUBContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EBNFParserListener); ok {
		listenerT.EnterSUB(s)
	}
}

func (s *SUBContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EBNFParserListener); ok {
		listenerT.ExitSUB(s)
	}
}

type REPContext struct {
	ChoiceContext
}

func NewREPContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *REPContext {
	var p = new(REPContext)

	InitEmptyChoiceContext(&p.ChoiceContext)
	p.parser = parser
	p.CopyAll(ctx.(*ChoiceContext))

	return p
}

func (s *REPContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *REPContext) REP() antlr.TerminalNode {
	return s.GetToken(EBNFParserREP, 0)
}

func (s *REPContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EBNFParserListener); ok {
		listenerT.EnterREP(s)
	}
}

func (s *REPContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EBNFParserListener); ok {
		listenerT.ExitREP(s)
	}
}

type PLUSContext struct {
	ChoiceContext
}

func NewPLUSContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *PLUSContext {
	var p = new(PLUSContext)

	InitEmptyChoiceContext(&p.ChoiceContext)
	p.parser = parser
	p.CopyAll(ctx.(*ChoiceContext))

	return p
}

func (s *PLUSContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PLUSContext) PLUS() antlr.TerminalNode {
	return s.GetToken(EBNFParserPLUS, 0)
}

func (s *PLUSContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EBNFParserListener); ok {
		listenerT.EnterPLUS(s)
	}
}

func (s *PLUSContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EBNFParserListener); ok {
		listenerT.ExitPLUS(s)
	}
}

func (p *EBNFParser) Choice() (localctx IChoiceContext) {
	localctx = NewChoiceContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, EBNFParserRULE_choice)
	p.SetState(77)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case EBNFParserREP:
		localctx = NewREPContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(73)
			p.Match(EBNFParserREP)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case EBNFParserPLUS:
		localctx = NewPLUSContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(74)
			p.Match(EBNFParserPLUS)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case EBNFParserEXT:
		localctx = NewEXTContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(75)
			p.Match(EBNFParserEXT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case EBNFParserSUB:
		localctx = NewSUBContext(p, localctx)
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(76)
			p.Match(EBNFParserSUB)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IIdentifierContext is an interface to support dynamic dispatch.
type IIdentifierContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ID() antlr.TerminalNode

	// IsIdentifierContext differentiates from other interfaces.
	IsIdentifierContext()
}

type IdentifierContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyIdentifierContext() *IdentifierContext {
	var p = new(IdentifierContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EBNFParserRULE_identifier
	return p
}

func InitEmptyIdentifierContext(p *IdentifierContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EBNFParserRULE_identifier
}

func (*IdentifierContext) IsIdentifierContext() {}

func NewIdentifierContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IdentifierContext {
	var p = new(IdentifierContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EBNFParserRULE_identifier

	return p
}

func (s *IdentifierContext) GetParser() antlr.Parser { return s.parser }

func (s *IdentifierContext) ID() antlr.TerminalNode {
	return s.GetToken(EBNFParserID, 0)
}

func (s *IdentifierContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IdentifierContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *IdentifierContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EBNFParserListener); ok {
		listenerT.EnterIdentifier(s)
	}
}

func (s *IdentifierContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EBNFParserListener); ok {
		listenerT.ExitIdentifier(s)
	}
}

func (p *EBNFParser) Identifier() (localctx IIdentifierContext) {
	localctx = NewIdentifierContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, EBNFParserRULE_identifier)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(79)
		p.Match(EBNFParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

func (p *EBNFParser) Sempred(localctx antlr.RuleContext, ruleIndex, predIndex int) bool {
	switch ruleIndex {
	case 4:
		var t *FactorContext = nil
		if localctx != nil {
			t = localctx.(*FactorContext)
		}
		return p.Factor_Sempred(t, predIndex)

	default:
		panic("No predicate with index: " + fmt.Sprint(ruleIndex))
	}
}

func (p *EBNFParser) Factor_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	switch predIndex {
	case 0:
		return p.Precpred(p.GetParserRuleContext(), 6)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}
