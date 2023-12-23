package parser

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/CUHK-SE-Group/generic-generator/log"
	"github.com/CUHK-SE-Group/generic-generator/parser/ebnf"
	"github.com/CUHK-SE-Group/generic-generator/schemas"
	"github.com/antlr4-go/antlr/v4"
)

/*
node搜索流程：
1. 能控制生成的只有 OR，REP。因此生成限制仅通过指定 某个Production下的操作是什么
xpath ..    node1.. node2
[SKIP, expression] func () { random ... }
e.g., 我要限制 SKIP 语句下的 expression 的 OR 生成是随机的。或 我要指定 SKIP 语句下的 REP 次数小于 3

2. 定位节点

先将生成逻辑挂载到对应的node类型中，例如是针对OR的。并且声明他的作用域是SKIP。
则在真正生成时，对这类节点的生成逻辑进行匹配，如果这个类型匹配到是SKIP，则应用这类生成逻辑。（可以保存生成路径，则一旦检测到生成路径里有对应的作用域，则应用该逻辑）

3. 生成变量的约束

同样，约束SKIP下的变量生成是从前面生成的某语句里sample，类似于上一步
*/

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
	textHandler := slog.NewTextHandler(os.Stderr, nil)
	callerInfoHandler := log.NewCallerInfoHandler(textHandler)
	logger := slog.New(callerInfoHandler)

	listener := &ebnfListener{
		currentSymbolId:   0,
		currentProduction: &schemas.Node{},
		stack:             []*schemas.Node{},
		logger:            logger,
		grammar:           schemas.NewGrammar(startSym),
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
	content := strings.Trim(c.GetText(), "'\"")
	l.addSymbolTop(schemas.NewNode(l.grammar, schemas.GrammarTerminal, l.generateId(), content))
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
