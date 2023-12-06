// Code generated from ./PathQuery.g4 by ANTLR 4.13.0. DO NOT EDIT.

package pathQuery // PathQuery
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

type PathQueryParser struct {
	*antlr.BaseParser
}

var PathQueryParserStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	LiteralNames           []string
	SymbolicNames          []string
	RuleNames              []string
	PredictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func pathqueryParserInit() {
	staticData := &PathQueryParserStaticData
	staticData.LiteralNames = []string{
		"", "'*'", "'/'", "'//'",
	}
	staticData.SymbolicNames = []string{
		"", "", "", "", "NODE_ID", "WS",
	}
	staticData.RuleNames = []string{
		"query", "segment", "rootNode", "pathSeparator",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 5, 33, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 1, 0, 3, 0,
		10, 8, 0, 1, 0, 1, 0, 1, 0, 1, 0, 5, 0, 16, 8, 0, 10, 0, 12, 0, 19, 9,
		0, 1, 0, 1, 0, 1, 1, 1, 1, 3, 1, 25, 8, 1, 1, 2, 1, 2, 1, 3, 1, 3, 3, 3,
		31, 8, 3, 1, 3, 0, 0, 4, 0, 2, 4, 6, 0, 0, 32, 0, 9, 1, 0, 0, 0, 2, 24,
		1, 0, 0, 0, 4, 26, 1, 0, 0, 0, 6, 30, 1, 0, 0, 0, 8, 10, 3, 4, 2, 0, 9,
		8, 1, 0, 0, 0, 9, 10, 1, 0, 0, 0, 10, 11, 1, 0, 0, 0, 11, 17, 3, 2, 1,
		0, 12, 13, 3, 6, 3, 0, 13, 14, 3, 2, 1, 0, 14, 16, 1, 0, 0, 0, 15, 12,
		1, 0, 0, 0, 16, 19, 1, 0, 0, 0, 17, 15, 1, 0, 0, 0, 17, 18, 1, 0, 0, 0,
		18, 20, 1, 0, 0, 0, 19, 17, 1, 0, 0, 0, 20, 21, 5, 0, 0, 1, 21, 1, 1, 0,
		0, 0, 22, 25, 5, 4, 0, 0, 23, 25, 5, 1, 0, 0, 24, 22, 1, 0, 0, 0, 24, 23,
		1, 0, 0, 0, 25, 3, 1, 0, 0, 0, 26, 27, 5, 2, 0, 0, 27, 5, 1, 0, 0, 0, 28,
		31, 5, 2, 0, 0, 29, 31, 5, 3, 0, 0, 30, 28, 1, 0, 0, 0, 30, 29, 1, 0, 0,
		0, 31, 7, 1, 0, 0, 0, 4, 9, 17, 24, 30,
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

// PathQueryParserInit initializes any static state used to implement PathQueryParser. By default the
// static state used to implement the parser is lazily initialized during the first call to
// NewPathQueryParser(). You can call this function if you wish to initialize the static state ahead
// of time.
func PathQueryParserInit() {
	staticData := &PathQueryParserStaticData
	staticData.once.Do(pathqueryParserInit)
}

// NewPathQueryParser produces a new parser instance for the optional input antlr.TokenStream.
func NewPathQueryParser(input antlr.TokenStream) *PathQueryParser {
	PathQueryParserInit()
	this := new(PathQueryParser)
	this.BaseParser = antlr.NewBaseParser(input)
	staticData := &PathQueryParserStaticData
	this.Interpreter = antlr.NewParserATNSimulator(this, staticData.atn, staticData.decisionToDFA, staticData.PredictionContextCache)
	this.RuleNames = staticData.RuleNames
	this.LiteralNames = staticData.LiteralNames
	this.SymbolicNames = staticData.SymbolicNames
	this.GrammarFileName = "PathQuery.g4"

	return this
}

// PathQueryParser tokens.
const (
	PathQueryParserEOF     = antlr.TokenEOF
	PathQueryParserT__0    = 1
	PathQueryParserT__1    = 2
	PathQueryParserT__2    = 3
	PathQueryParserNODE_ID = 4
	PathQueryParserWS      = 5
)

// PathQueryParser rules.
const (
	PathQueryParserRULE_query         = 0
	PathQueryParserRULE_segment       = 1
	PathQueryParserRULE_rootNode      = 2
	PathQueryParserRULE_pathSeparator = 3
)

// IQueryContext is an interface to support dynamic dispatch.
type IQueryContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllSegment() []ISegmentContext
	Segment(i int) ISegmentContext
	EOF() antlr.TerminalNode
	RootNode() IRootNodeContext
	AllPathSeparator() []IPathSeparatorContext
	PathSeparator(i int) IPathSeparatorContext

	// IsQueryContext differentiates from other interfaces.
	IsQueryContext()
}

type QueryContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyQueryContext() *QueryContext {
	var p = new(QueryContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PathQueryParserRULE_query
	return p
}

func InitEmptyQueryContext(p *QueryContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PathQueryParserRULE_query
}

func (*QueryContext) IsQueryContext() {}

func NewQueryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *QueryContext {
	var p = new(QueryContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PathQueryParserRULE_query

	return p
}

func (s *QueryContext) GetParser() antlr.Parser { return s.parser }

func (s *QueryContext) AllSegment() []ISegmentContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ISegmentContext); ok {
			len++
		}
	}

	tst := make([]ISegmentContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ISegmentContext); ok {
			tst[i] = t.(ISegmentContext)
			i++
		}
	}

	return tst
}

func (s *QueryContext) Segment(i int) ISegmentContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISegmentContext); ok {
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

	return t.(ISegmentContext)
}

func (s *QueryContext) EOF() antlr.TerminalNode {
	return s.GetToken(PathQueryParserEOF, 0)
}

func (s *QueryContext) RootNode() IRootNodeContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRootNodeContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRootNodeContext)
}

func (s *QueryContext) AllPathSeparator() []IPathSeparatorContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IPathSeparatorContext); ok {
			len++
		}
	}

	tst := make([]IPathSeparatorContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IPathSeparatorContext); ok {
			tst[i] = t.(IPathSeparatorContext)
			i++
		}
	}

	return tst
}

func (s *QueryContext) PathSeparator(i int) IPathSeparatorContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPathSeparatorContext); ok {
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

	return t.(IPathSeparatorContext)
}

func (s *QueryContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *QueryContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *QueryContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PathQueryListener); ok {
		listenerT.EnterQuery(s)
	}
}

func (s *QueryContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PathQueryListener); ok {
		listenerT.ExitQuery(s)
	}
}

func (p *PathQueryParser) Query() (localctx IQueryContext) {
	localctx = NewQueryContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, PathQueryParserRULE_query)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(9)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == PathQueryParserT__1 {
		{
			p.SetState(8)
			p.RootNode()
		}

	}
	{
		p.SetState(11)
		p.Segment()
	}
	p.SetState(17)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == PathQueryParserT__1 || _la == PathQueryParserT__2 {
		{
			p.SetState(12)
			p.PathSeparator()
		}
		{
			p.SetState(13)
			p.Segment()
		}

		p.SetState(19)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(20)
		p.Match(PathQueryParserEOF)
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

// ISegmentContext is an interface to support dynamic dispatch.
type ISegmentContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsSegmentContext differentiates from other interfaces.
	IsSegmentContext()
}

type SegmentContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySegmentContext() *SegmentContext {
	var p = new(SegmentContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PathQueryParserRULE_segment
	return p
}

func InitEmptySegmentContext(p *SegmentContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PathQueryParserRULE_segment
}

func (*SegmentContext) IsSegmentContext() {}

func NewSegmentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SegmentContext {
	var p = new(SegmentContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PathQueryParserRULE_segment

	return p
}

func (s *SegmentContext) GetParser() antlr.Parser { return s.parser }

func (s *SegmentContext) CopyAll(ctx *SegmentContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *SegmentContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SegmentContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type NodeContext struct {
	SegmentContext
}

func NewNodeContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *NodeContext {
	var p = new(NodeContext)

	InitEmptySegmentContext(&p.SegmentContext)
	p.parser = parser
	p.CopyAll(ctx.(*SegmentContext))

	return p
}

func (s *NodeContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NodeContext) NODE_ID() antlr.TerminalNode {
	return s.GetToken(PathQueryParserNODE_ID, 0)
}

func (s *NodeContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PathQueryListener); ok {
		listenerT.EnterNode(s)
	}
}

func (s *NodeContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PathQueryListener); ok {
		listenerT.ExitNode(s)
	}
}

type AnyContext struct {
	SegmentContext
}

func NewAnyContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *AnyContext {
	var p = new(AnyContext)

	InitEmptySegmentContext(&p.SegmentContext)
	p.parser = parser
	p.CopyAll(ctx.(*SegmentContext))

	return p
}

func (s *AnyContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AnyContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PathQueryListener); ok {
		listenerT.EnterAny(s)
	}
}

func (s *AnyContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PathQueryListener); ok {
		listenerT.ExitAny(s)
	}
}

func (p *PathQueryParser) Segment() (localctx ISegmentContext) {
	localctx = NewSegmentContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, PathQueryParserRULE_segment)
	p.SetState(24)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case PathQueryParserNODE_ID:
		localctx = NewNodeContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(22)
			p.Match(PathQueryParserNODE_ID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case PathQueryParserT__0:
		localctx = NewAnyContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(23)
			p.Match(PathQueryParserT__0)
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

// IRootNodeContext is an interface to support dynamic dispatch.
type IRootNodeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsRootNodeContext differentiates from other interfaces.
	IsRootNodeContext()
}

type RootNodeContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyRootNodeContext() *RootNodeContext {
	var p = new(RootNodeContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PathQueryParserRULE_rootNode
	return p
}

func InitEmptyRootNodeContext(p *RootNodeContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PathQueryParserRULE_rootNode
}

func (*RootNodeContext) IsRootNodeContext() {}

func NewRootNodeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RootNodeContext {
	var p = new(RootNodeContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PathQueryParserRULE_rootNode

	return p
}

func (s *RootNodeContext) GetParser() antlr.Parser { return s.parser }
func (s *RootNodeContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RootNodeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *RootNodeContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PathQueryListener); ok {
		listenerT.EnterRootNode(s)
	}
}

func (s *RootNodeContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PathQueryListener); ok {
		listenerT.ExitRootNode(s)
	}
}

func (p *PathQueryParser) RootNode() (localctx IRootNodeContext) {
	localctx = NewRootNodeContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, PathQueryParserRULE_rootNode)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(26)
		p.Match(PathQueryParserT__1)
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

// IPathSeparatorContext is an interface to support dynamic dispatch.
type IPathSeparatorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsPathSeparatorContext differentiates from other interfaces.
	IsPathSeparatorContext()
}

type PathSeparatorContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPathSeparatorContext() *PathSeparatorContext {
	var p = new(PathSeparatorContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PathQueryParserRULE_pathSeparator
	return p
}

func InitEmptyPathSeparatorContext(p *PathSeparatorContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = PathQueryParserRULE_pathSeparator
}

func (*PathSeparatorContext) IsPathSeparatorContext() {}

func NewPathSeparatorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PathSeparatorContext {
	var p = new(PathSeparatorContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = PathQueryParserRULE_pathSeparator

	return p
}

func (s *PathSeparatorContext) GetParser() antlr.Parser { return s.parser }

func (s *PathSeparatorContext) CopyAll(ctx *PathSeparatorContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *PathSeparatorContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PathSeparatorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type AllContext struct {
	PathSeparatorContext
}

func NewAllContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *AllContext {
	var p = new(AllContext)

	InitEmptyPathSeparatorContext(&p.PathSeparatorContext)
	p.parser = parser
	p.CopyAll(ctx.(*PathSeparatorContext))

	return p
}

func (s *AllContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AllContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PathQueryListener); ok {
		listenerT.EnterAll(s)
	}
}

func (s *AllContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PathQueryListener); ok {
		listenerT.ExitAll(s)
	}
}

type ChildContext struct {
	PathSeparatorContext
}

func NewChildContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ChildContext {
	var p = new(ChildContext)

	InitEmptyPathSeparatorContext(&p.PathSeparatorContext)
	p.parser = parser
	p.CopyAll(ctx.(*PathSeparatorContext))

	return p
}

func (s *ChildContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ChildContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PathQueryListener); ok {
		listenerT.EnterChild(s)
	}
}

func (s *ChildContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PathQueryListener); ok {
		listenerT.ExitChild(s)
	}
}

func (p *PathQueryParser) PathSeparator() (localctx IPathSeparatorContext) {
	localctx = NewPathSeparatorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, PathQueryParserRULE_pathSeparator)
	p.SetState(30)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case PathQueryParserT__1:
		localctx = NewChildContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(28)
			p.Match(PathQueryParserT__1)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case PathQueryParserT__2:
		localctx = NewAllContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(29)
			p.Match(PathQueryParserT__2)
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
