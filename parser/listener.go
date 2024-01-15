package parser

import (
	"fmt"
	"github.com/CUHK-SE-Group/generic-generator/log"
	"github.com/CUHK-SE-Group/generic-generator/parser/ebnf"
	"github.com/CUHK-SE-Group/generic-generator/schemas"
	"github.com/antlr4-go/antlr/v4"
	"log/slog"
	"os"
	"strconv"
)

type ebnfListener struct {
	*ebnf.BaseEBNFParserListener
	stack             []*schemas.Node
	currentSymbolId   int
	currentProduction *schemas.Node
	logger            *slog.Logger
	grammar           *schemas.Grammar
	productions       map[string]*schemas.Node
	popStack          []int
}

func newEbnfListener(startSym string) *ebnfListener {
	textHandler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelError,
	})
	callerInfoHandler := log.NewCallerInfoHandler(textHandler)
	logger := slog.New(callerInfoHandler)

	listener := &ebnfListener{
		currentSymbolId:   0,
		currentProduction: &schemas.Node{},
		stack:             []*schemas.Node{},
		logger:            logger,
		grammar:           schemas.NewGrammar(schemas.WithStartSym(startSym)),
		productions:       map[string]*schemas.Node{},
	}
	return listener
}

func (l *ebnfListener) generateId() string {
	id := l.currentProduction.GetID() + "#" + strconv.Itoa(l.currentSymbolId)
	l.currentSymbolId++
	return id
}

func (l *ebnfListener) top() *schemas.Node {
	return l.stack[len(l.stack)-1]
}

func (l *ebnfListener) pop() {
	l.stack = l.stack[:len(l.stack)-1]
}

func (l *ebnfListener) enter(status int) {
	l.popStack = append(l.popStack, status)
}

func (l *ebnfListener) exit() {
	last := l.popStack[len(l.popStack)-1]
	l.popStack = l.popStack[:len(l.popStack)-1]
	if last != 0 {
		l.pop()
	}
}

func (l *ebnfListener) push(g *schemas.Node) {
	l.stack = append(l.stack, g)
}

func (l *ebnfListener) clear() {
	l.stack = []*schemas.Node{}
}
func (l *ebnfListener) empty() bool {
	return len(l.stack) == 0
}

func (l *ebnfListener) addSymbolTop(n *schemas.Node) {
	l.top().AddSymbol(n)
}

func (l *ebnfListener) addThenPush(n *schemas.Node) {
	l.addSymbolTop(n)
	l.push(n)
}

func (l *ebnfListener) EnterProduction(c *ebnf.ProductionContext) {
	l.currentSymbolId = 0
	l.logger.Debug("production", "id", c.ID().GetText(), "expr", c.Expr().GetText())
	name := c.ID().GetText()
	cur, ok := l.productions[name]
	if !ok {
		cur = schemas.NewNode(l.grammar, schemas.GrammarProduction, name, c.Expr().GetText())
	}
	l.currentProduction = cur
	l.productions[name] = cur
	l.clear()
	l.push(cur)
}
func (l *ebnfListener) ExitProduction(c *ebnf.ProductionContext) {
	l.pop()
	for !l.empty() {
		l.logger.Error("stack do not equal to 0", "id", l.top().GetID(), "content", l.top().GetContent())
		l.pop()
	}
}

func (l *ebnfListener) EnterExpr(c *ebnf.ExprContext) {
	l.logger.Debug("entered expr", fmt.Sprint(c.GetRuleIndex()), c.GetText())
	if len(c.AllCOMMA()) != 0 {
		l.addThenPush(schemas.NewNode(l.grammar, schemas.GrammarCatenate, l.generateId(), c.GetText()))
		l.enter(1)
	} else {
		l.enter(0)
	}
}

func (l *ebnfListener) ExitExpr(c *ebnf.ExprContext) {
	l.exit()
}

func (l *ebnfListener) EnterTerm(c *ebnf.TermContext) {
	l.logger.Debug("entered term", fmt.Sprint(c.GetRuleIndex()), c.GetText())
	if len(c.AllOR()) != 0 {
		l.addThenPush(schemas.NewNode(l.grammar, schemas.GrammarOR, l.generateId(), c.GetText()))
		l.enter(1)
	} else {
		l.enter(0)
	}
}

func (l *ebnfListener) ExitTerm(c *ebnf.TermContext) {
	l.exit()
}

func (l *ebnfListener) EnterID(c *ebnf.IDContext) {
	l.logger.Debug("encountered id", "val", c.GetText())
	l.addSymbolTop(schemas.NewNode(l.grammar, schemas.GrammarID, l.generateId(), c.GetText()))
}

func (l *ebnfListener) EnterQUOTE(c *ebnf.QUOTEContext) {
	l.logger.Debug("encountered quote", "val", c.GetText())
	l.addSymbolTop(schemas.NewNode(l.grammar, schemas.GrammarTerminal, l.generateId(), c.GetText()))
}

func (l *ebnfListener) EnterCHOICE(c *ebnf.CHOICEContext) {
	l.logger.Debug("entered choice", fmt.Sprint(c.GetRuleIndex()), c.GetText())
	l.addThenPush(schemas.NewNode(l.grammar, schemas.GrammarChoice, l.generateId(), c.GetText()))
}

func (l *ebnfListener) ExitCHOICE(c *ebnf.CHOICEContext) {
	l.pop()
}

func (l *ebnfListener) EnterBRACE(c *ebnf.BRACEContext) {
	l.logger.Debug("entered brace", fmt.Sprint(c.GetRuleIndex()), c.GetText())
	l.addThenPush(schemas.NewNode(l.grammar, schemas.GrammarREP, l.generateId(), c.GetText()))
}

func (l *ebnfListener) ExitBRACE(c *ebnf.BRACEContext) {
	l.pop()
}

func (l *ebnfListener) EnterREP(c *ebnf.REPContext) {
	l.logger.Debug("entered rep", fmt.Sprint(c.GetRuleIndex()), c.GetText())
	if l.top().GetType() != schemas.GrammarChoice {
		l.logger.Error("parent is not choice", "id", l.top().GetID(), "content", l.top().GetContent())
		os.Exit(1)
	}
	l.top().SetType(schemas.GrammarREP)
}

func (l *ebnfListener) EnterPLUS(c *ebnf.PLUSContext) {
	l.logger.Debug("entered rep", fmt.Sprint(c.GetRuleIndex()), c.GetText())
	if l.top().GetType() != schemas.GrammarChoice {
		l.logger.Error("parent is not choice", "id", l.top().GetID(), "content", l.top().GetContent())
		os.Exit(1)
	}
	l.top().SetType(schemas.GrammarPLUS)
}

func (l *ebnfListener) EnterEXT(c *ebnf.EXTContext) {
	l.logger.Debug("entered rep", fmt.Sprint(c.GetRuleIndex()), c.GetText())
	if l.top().GetType() != schemas.GrammarChoice {
		l.logger.Error("parent is not choice", "id", l.top().GetID(), "content", l.top().GetContent())
		os.Exit(1)
	}
	l.top().SetType(schemas.GrammarEXT)
}

func (l *ebnfListener) EnterSUB(c *ebnf.SUBContext) {
	l.logger.Debug("entered rep", fmt.Sprint(c.GetRuleIndex()), c.GetText())
	if l.top().GetType() != schemas.GrammarChoice {
		l.logger.Error("parent is not choice", "id", l.top().GetID(), "content", l.top().GetContent())
		os.Exit(1)
	}
	l.top().SetType(schemas.GrammarSUB)
}

func (l *ebnfListener) EnterBRACKET(c *ebnf.BRACKETContext) {
	l.logger.Debug("entered bracket", fmt.Sprint(c.GetRuleIndex()), c.GetText())
	l.addThenPush(schemas.NewNode(l.grammar, schemas.GrammarOptional, l.generateId(), c.GetText()))
}

func (l *ebnfListener) ExitBRACKET(c *ebnf.BRACKETContext) {
	l.pop()
}

func Parse(file string, startSym string) (*schemas.Grammar, error) {
	is, err := antlr.NewFileStream(file)
	if err != nil {
		return nil, err
	}
	lexer := ebnf.NewEBNFLexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	parser := ebnf.NewEBNFParser(stream)
	listener := newEbnfListener(startSym)
	antlr.ParseTreeWalkerDefault.Walk(listener, parser.Ebnf())

	return listener.grammar, nil
}
