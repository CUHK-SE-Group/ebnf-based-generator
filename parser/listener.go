package parser

import (
	"github.com/CUHK-SE-Group/ebnf-based-generator/log"
	"github.com/CUHK-SE-Group/ebnf-based-generator/parser/ebnf"
	"github.com/CUHK-SE-Group/ebnf-based-generator/schemas"
	"github.com/antlr4-go/antlr/v4"
	"log/slog"
	"os"
	"strconv"
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
	stack             []*schemas.Grammar
	currentSymbolId   int
	currentProduction *schemas.Grammar
	productions       map[string]*schemas.Grammar
	logger            *slog.Logger
}

func newEbnfListener() *ebnfListener {
	textHandler := slog.NewTextHandler(os.Stderr, nil)
	callerInfoHandler := log.NewCallerInfoHandler(textHandler)
	logger := slog.New(callerInfoHandler)

	listener := &ebnfListener{
		currentSymbolId:   0,
		currentProduction: &schemas.Grammar{},
		productions:       map[string]*schemas.Grammar{},
		stack:             []*schemas.Grammar{},
		logger:            logger,
	}
	return listener
}

func (l *ebnfListener) generateId() string {
	id := l.currentProduction.GetID() + "#" + strconv.Itoa(l.currentSymbolId)
	l.currentSymbolId++
	return id
}

func (l *ebnfListener) top() *schemas.Grammar {
	return l.stack[len(l.stack)-1]
}

func (l *ebnfListener) pop() {
	l.stack = l.stack[:len(l.stack)-1]
}

func (l *ebnfListener) push(g *schemas.Grammar) {
	l.stack = append(l.stack, g)
}

func (l *ebnfListener) save(g *schemas.Grammar) {
	l.productions[g.GetID()] = g
}

func (l *ebnfListener) clear() {
	l.stack = []*schemas.Grammar{}
}

func (l *ebnfListener) EnterProduction(c *ebnf.ProductionContext) {
	l.logger.Info(c.GetText())
	name := c.ID().GetText()

	if p, ok := l.productions[name]; ok {
		l.currentProduction = p
		l.clear()
		l.push(p)
	} else {
		g := schemas.NewGrammar(schemas.GrammarProduction, name, c.GetText(), nil)
		l.currentProduction = g
		l.clear()
		l.save(g)
		l.push(g)
	}
}
func (l *ebnfListener) ExitProduction(c *ebnf.ProductionContext) {
	l.pop()
}

func (l *ebnfListener) EnterExpr(c *ebnf.ExprContext) {
	g := schemas.NewGrammar(schemas.GrammarExpr, l.generateId(), c.GetText(), nil)
	l.top().AddSymbol(g)
	l.push(g)
}
func (l *ebnfListener) ExitExpr(c *ebnf.ExprContext) {
	l.pop()
}

func (l *ebnfListener) EnterTerm(c *ebnf.TermContext) {
	l.logger.Info(c.GetText())
	g := schemas.NewGrammar(schemas.GrammarTerm, l.generateId(), c.GetText(), nil)
	l.top().AddSymbol(g)
	l.push(g)
}

func (l *ebnfListener) ExitTerm(c *ebnf.TermContext) {
	l.pop()
}

// EnterID one of factor branch
func (l *ebnfListener) EnterID(c *ebnf.IDContext) {
	l.logger.Info(c.GetText())
	g := schemas.NewGrammar(schemas.GrammarID, l.generateId(), c.GetText(), nil)
	l.top().AddSymbol(g)
	l.push(g)
}
func (l *ebnfListener) ExitID(c *ebnf.IDContext) {
	l.pop()
}

//// EnterSTR one of factor branch
//func (l *ebnfListener) EnterSTR(c *ebnf.STRContext) {
//	l.logger.Info(c.GetText())
//	g := schemas.NewGrammar( schemas.GrammarTerminal, l.generateId(), c.GetText(), nil)
//	l.top().AddSymbol(g)
//	l.push(g)
//}
//func (l *ebnfListener) ExitSTR(c *ebnf.STRContext) {
//	l.pop()
//}

// EnterQUOTE one of factor branch
func (l *ebnfListener) EnterQUOTE(c *ebnf.QUOTEContext) {
	l.logger.Info(c.GetText())
	g := schemas.NewGrammar(schemas.GrammarTerminal, l.generateId(), c.GetText(), nil)
	l.top().AddSymbol(g)
	l.push(g)
}
func (l *ebnfListener) ExitQUOTE(c *ebnf.QUOTEContext) {
	l.pop()
}

// one of factor branch
func (l *ebnfListener) EnterPAREN(c *ebnf.PARENContext) {
	l.logger.Info(c.GetText())
	g := schemas.NewGrammar(schemas.GrammarPAREN, l.generateId(), c.GetText(), nil)
	l.top().AddSymbol(g)
	l.push(g)
}
func (l *ebnfListener) ExitPAREN(c *ebnf.PARENContext) {
	l.pop()
}

// one of factor branch
func (l *ebnfListener) EnterBRACKET(c *ebnf.BRACKETContext) {
	l.logger.Info(c.GetText())
	g := schemas.NewGrammar(schemas.GrammarBRACKET, l.generateId(), c.GetText(), nil)
	l.top().AddSymbol(g)
	l.push(g)
}
func (l *ebnfListener) ExitBRACKET(c *ebnf.BRACKETContext) {
	l.pop()
}

// one of factor branch
func (l *ebnfListener) EnterBRACE(c *ebnf.BRACEContext) {
	l.logger.Info(c.GetText())
	g := schemas.NewGrammar(schemas.GrammarBRACE, l.generateId(), c.GetText(), nil)
	l.top().AddSymbol(g)
	l.push(g)
}
func (l *ebnfListener) ExitBRACE(c *ebnf.BRACEContext) {
	l.pop()
}

// one of factor branch
func (l *ebnfListener) EnterREP(c *ebnf.REPContext) {
	l.logger.Info(c.GetText())
	g := schemas.NewGrammar(schemas.GrammarREP, l.generateId(), c.GetText(), nil)
	l.top().AddSymbol(g)
	l.push(g)
}

// one of factor branch
func (l *ebnfListener) ExitREP(c *ebnf.REPContext) {
	l.pop()
}

// one of factor branch
func (l *ebnfListener) EnterPLUS(c *ebnf.PLUSContext) {
	l.logger.Info(c.GetText())
	g := schemas.NewGrammar(schemas.GrammarPLUS, l.generateId(), c.GetText(), nil)
	l.top().AddSymbol(g)
	l.push(g)
}

// one of factor branch
func (l *ebnfListener) ExitPLUS(c *ebnf.PLUSContext) {
	l.pop()
}

// one of factor branch
func (l *ebnfListener) EnterEXT(c *ebnf.EXTContext) {
	l.logger.Info(c.GetText())
	g := schemas.NewGrammar(schemas.GrammarEXT, l.generateId(), c.GetText(), nil)
	l.top().AddSymbol(g)
	l.push(g)
}

// one of factor branch
func (l *ebnfListener) ExitEXT(c *ebnf.EXTContext) {
	l.pop()
}

// one of factor branch
func (l *ebnfListener) EnterSUB(c *ebnf.SUBContext) {
	l.logger.Info(c.GetText())
	g := schemas.NewGrammar(schemas.GrammarSUB, l.generateId(), c.GetText(), nil)
	l.top().AddSymbol(g)
	l.push(g)
}

// one of factor branch
func (l *ebnfListener) ExitSUB(c *ebnf.SUBContext) {
	l.pop()
}

func Parse(file string) (map[string]*schemas.Grammar, error) {
	is, err := antlr.NewFileStream(file)
	if err != nil {
		return nil, err
	}
	lexer := ebnf.NewEBNFLexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	parser := ebnf.NewEBNFParser(stream)
	listener := newEbnfListener()
	antlr.ParseTreeWalkerDefault.Walk(listener, parser.Ebnf())
	return listener.productions, nil
}

func MergeProduction(p map[string]*schemas.Grammar, startSymbol string) *schemas.Grammar {
	root, ok := p[startSymbol]
	if !ok {
		return nil
	}

	queue := make([]*schemas.Grammar, 0)
	queue = append(queue, root)
	visited := make([]*schemas.Grammar, 0)
	for len(queue) != 0 {
		cur := queue[0]
		queue = queue[1:]
		if len(*cur.GetSymbols()) == 0 {
			if cur.GetType() == schemas.GrammarID {
				visited = append(visited, cur)
			}
		}
		for _, v := range *cur.GetSymbols() {
			queue = append(queue, v)
		}
	}
	for _, v := range visited {
		*v.GetSymbols() = append(*v.GetSymbols(), p[v.GetContent()])
	}
	return root
}
