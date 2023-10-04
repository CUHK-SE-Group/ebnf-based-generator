package parser

import (
	"context"
	"fmt"
	"github.com/CUHK-SE-Group/ebnf-based-generator/parser/ebnf"
	"log/slog"
	"os"
	"runtime"
	"strconv"
	"strings"
)

/*
node搜索流程：
1. 能控制生成的只有 OR，REP。因此生成限制仅通过指定 某个Production下的操作是什么

e.g., 我要限制 SKIP 语句下的 expression 的 OR 生成是随机的。或 我要指定 SKIP 语句下的 REP 次数小于 3

2. 定位节点

先将生成逻辑挂载到对应的node类型中，例如是针对OR的。并且声明他的作用域是SKIP。
则在真正生成时，对这类节点的生成逻辑进行匹配，如果这个类型匹配到是SKIP，则应用这类生成逻辑。（可以保存生成路径，则一旦检测到生成路径里有对应的作用域，则应用该逻辑）

3. 生成变量的约束

同样，约束SKIP下的变量生成是从前面生成的某语句里sample，类似于上一步
*/

type CallerInfoHandler struct {
	innerHandler slog.Handler
}

func (h *CallerInfoHandler) Handle(ctx context.Context, r slog.Record) error {
	pc, file, _, ok := runtime.Caller(3) // Adjust the skip value as needed
	if ok {
		shortFile := file[strings.LastIndex(file, "/")+1:]
		funcName := runtime.FuncForPC(pc).Name()
		shortFuncName := funcName[strings.LastIndex(funcName, ".")+1:]
		r.Message = fmt.Sprintf("%s:%s: %s", shortFile, shortFuncName, r.Message)
	}
	return h.innerHandler.Handle(ctx, r)
}

func (h *CallerInfoHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.innerHandler.Enabled(ctx, level)
}

func (h *CallerInfoHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &CallerInfoHandler{innerHandler: h.innerHandler.WithAttrs(attrs)}
}

func (h *CallerInfoHandler) WithGroup(name string) slog.Handler {
	return &CallerInfoHandler{innerHandler: h.innerHandler.WithGroup(name)}
}

func NewCallerInfoHandler(innerHandler slog.Handler) *CallerInfoHandler {
	return &CallerInfoHandler{innerHandler: innerHandler}
}

type ebnfListener struct {
	*ebnf.BaseEBNFParserListener
	stack             []*Grammar
	currentSymbolId   int
	currentProduction *Grammar
	productions       map[string]*Grammar
	logger            *slog.Logger
}

func newEbnfListener() *ebnfListener {
	textHandler := slog.NewTextHandler(os.Stdout, nil)
	callerInfoHandler := NewCallerInfoHandler(textHandler)
	logger := slog.New(callerInfoHandler)

	listener := &ebnfListener{
		currentSymbolId:   0,
		currentProduction: &Grammar{},
		productions:       map[string]*Grammar{},
		stack:             []*Grammar{},
		logger:            logger,
	}
	return listener
}

func (l *ebnfListener) generateId() string {
	id := l.currentProduction.ID + "#" + strconv.Itoa(l.currentSymbolId)
	l.currentSymbolId++
	return id
}

func (l *ebnfListener) top() *Grammar {
	return l.stack[len(l.stack)-1]
}

func (l *ebnfListener) pop() {
	l.stack = l.stack[:len(l.stack)-1]
}

func (l *ebnfListener) push(g *Grammar) {
	l.stack = append(l.stack, g)
}

func (l *ebnfListener) save(g *Grammar) {
	l.productions[g.ID] = g
}

func (l *ebnfListener) clear() {
	l.stack = []*Grammar{}
}

func (l *ebnfListener) EnterProduction(c *ebnf.ProductionContext) {
	l.logger.Info(c.GetText())
	name := c.ID().GetText()

	if p, ok := l.productions[name]; ok {
		l.currentProduction = p
		l.clear()
		l.push(p)
	} else {
		g := NewGrammar(NewContext(), GrammarProduction, name, c.GetText(), nil)
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
	// expr 这层暂时不管，因为只有单个选择。但为了保持尊重，新建一个symbol
	g := NewGrammar(NewContext(), GrammarExpr, l.generateId(), c.GetText(), nil)
	l.top().AddSymbol(g)
	l.push(g)
}
func (l *ebnfListener) ExitExpr(c *ebnf.ExprContext) {
	l.pop()
}

func (l *ebnfListener) EnterTerm(c *ebnf.TermContext) {
	l.logger.Info(c.GetText())
	g := NewGrammar(NewContext(), GrammarTerm, l.generateId(), c.GetText(), nil)
	l.top().AddSymbol(g)
	l.push(g)
}

func (l *ebnfListener) ExitTerm(c *ebnf.TermContext) {
	l.pop()
}

// EnterID one of factor branch
func (l *ebnfListener) EnterID(c *ebnf.IDContext) {
	l.logger.Info(c.GetText())
	g := NewGrammar(NewContext(), GrammarID, l.generateId(), c.GetText(), nil)
	l.top().AddSymbol(g)
	l.push(g)
}
func (l *ebnfListener) ExitID(c *ebnf.IDContext) {
	l.pop()
}

// EnterSTR one of factor branch
func (l *ebnfListener) EnterSTR(c *ebnf.STRContext) {
	l.logger.Info(c.GetText())
	g := NewGrammar(NewContext(), GrammarTerminal, l.generateId(), c.GetText(), nil)
	l.top().AddSymbol(g)
	l.push(g)
}
func (l *ebnfListener) ExitSTR(c *ebnf.STRContext) {
	l.pop()
}

// EnterQUOTE one of factor branch
func (l *ebnfListener) EnterQUOTE(c *ebnf.QUOTEContext) {
	l.logger.Info(c.GetText())
	g := NewGrammar(NewContext(), GrammarTerminal, l.generateId(), c.GetText(), nil)
	l.top().AddSymbol(g)
	l.push(g)
}
func (l *ebnfListener) ExitQUOTE(c *ebnf.QUOTEContext) {
	l.pop()
}

// one of factor branch
func (l *ebnfListener) EnterPAREN(c *ebnf.PARENContext) {
	l.logger.Info(c.GetText())
	g := NewGrammar(NewContext(), GrammarPAREN, l.generateId(), c.GetText(), nil)
	l.top().AddSymbol(g)
	l.push(g)
}
func (l *ebnfListener) ExitPAREN(c *ebnf.PARENContext) {
	l.pop()
}

// one of factor branch
func (l *ebnfListener) EnterBRACKET(c *ebnf.BRACKETContext) {
	l.logger.Info(c.GetText())
	g := NewGrammar(NewContext(), GrammarBRACKET, l.generateId(), c.GetText(), nil)
	l.top().AddSymbol(g)
	l.push(g)
}
func (l *ebnfListener) ExitBRACKET(c *ebnf.BRACKETContext) {
	l.pop()
}

// one of factor branch
func (l *ebnfListener) EnterBRACE(c *ebnf.BRACEContext) {
	l.logger.Info(c.GetText())
	g := NewGrammar(NewContext(), GrammarBRACE, l.generateId(), c.GetText(), nil)
	l.top().AddSymbol(g)
	l.push(g)
}
func (l *ebnfListener) ExitBRACE(c *ebnf.BRACEContext) {
	l.pop()
}

// one of factor branch
func (l *ebnfListener) EnterRECU(c *ebnf.RECUContext) {
	l.logger.Info(c.GetText())
	g := NewGrammar(NewContext(), GrammarRECU, l.generateId(), c.GetText(), nil)
	l.top().AddSymbol(g)
	l.push(g)
}

// one of factor branch
func (l *ebnfListener) ExitRECU(c *ebnf.RECUContext) {
	l.pop()
}
