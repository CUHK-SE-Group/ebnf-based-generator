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
	exprCnt           int
	termCnt           int
}

func newEbnfListener() *ebnfListener {
	textHandler := slog.NewTextHandler(os.Stderr, nil)
	callerInfoHandler := log.NewCallerInfoHandler(textHandler)
	logger := slog.New(callerInfoHandler)

	listener := &ebnfListener{
		currentSymbolId:   0,
		currentProduction: &schemas.Node{},
		stack:             []*schemas.Node{},
		logger:            logger,
		grammar:           schemas.NewGrammar(),
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

func (l *ebnfListener) push(g *schemas.Node) {
	l.stack = append(l.stack, g)
}

func (l *ebnfListener) clear() {
	l.stack = []*schemas.Node{}
}
func (l *ebnfListener) empty() bool {
	return len(l.stack) == 0
}

func (l *ebnfListener) EnterProduction(c *ebnf.ProductionContext) {
	l.currentSymbolId = 0
	l.logger.Info("production", "id", c.ID().GetText(), "expr", c.Expr().GetText())
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
	children := make([]*schemas.Node, 0)
	for i, v := range c.AllTerm() {
		l.logger.Info("expr", fmt.Sprint(i), v.GetText())
		n := schemas.NewNode(l.grammar, schemas.GrammarExpr, l.generateId(), v.GetText())
		l.top().AddSymbol(n)
		children = append(children, n)
	}

	for i := len(children) - 1; i >= 0; i-- {
		l.push(children[i])
		l.exprCnt++
	}

}
func (l *ebnfListener) ExitExpr(c *ebnf.ExprContext) {
	for l.exprCnt > 0 {
		l.pop()
		l.exprCnt--
	}
}

func (l *ebnfListener) EnterTerm(c *ebnf.TermContext) {
	children := make([]*schemas.Node, 0)
	for i, v := range c.AllFactor() {
		l.logger.Info("term", fmt.Sprint(i), v.GetText())
		n := schemas.NewNode(l.grammar, schemas.GrammarTerm, l.generateId(), v.GetText())
		l.top().AddSymbol(n)
		children = append(children, n)
	}
	for i := len(children) - 1; i >= 0; i-- {
		l.push(children[i])
		l.termCnt++
	}
}

func (l *ebnfListener) ExitTerm(c *ebnf.TermContext) {
	for l.termCnt > 0 {
		l.pop()
		l.termCnt--
	}
}

//
//// EnterID one of factor branch
//func (l *ebnfListener) EnterID(c *ebnf.IDContext) {
//	l.logger.Info("id", "val", c.GetText())
//	n := schemas.NewNode(l.grammar, schemas.GrammarID, l.generateId(), c.GetText())
//	l.top().AddSymbol(n)
//	l.push(n)
//}
//func (l *ebnfListener) ExitID(c *ebnf.IDContext) {
//	l.pop()
//}
//
//// EnterQUOTE one of factor branch
//func (l *ebnfListener) EnterQUOTE(c *ebnf.QUOTEContext) {
//	l.logger.Info("quote", "val", c.GetText())
//}
//func (l *ebnfListener) ExitQUOTE(c *ebnf.QUOTEContext) {
//}
//
//// one of factor branch
//func (l *ebnfListener) EnterPAREN(c *ebnf.PARENContext) {
//	l.logger.Info("paren", "val", c.GetText())
//}
//func (l *ebnfListener) ExitPAREN(c *ebnf.PARENContext) {
//}
//
//// one of factor branch
//func (l *ebnfListener) EnterBRACKET(c *ebnf.BRACKETContext) {
//	l.logger.Info("bracket", "val", c.GetText())
//}
//func (l *ebnfListener) ExitBRACKET(c *ebnf.BRACKETContext) {
//}
//
//// one of factor branch
//func (l *ebnfListener) EnterBRACE(c *ebnf.BRACEContext) {
//	l.logger.Info("brace", "val", c.GetText())
//}
//func (l *ebnfListener) ExitBRACE(c *ebnf.BRACEContext) {
//}
//
//// one of factor branch
//func (l *ebnfListener) EnterREP(c *ebnf.REPContext) {
//	l.logger.Info("rep", "val", c.GetText())
//}
//
//// one of factor branch
//func (l *ebnfListener) ExitREP(c *ebnf.REPContext) {
//}
//
//// one of factor branch
//func (l *ebnfListener) EnterPLUS(c *ebnf.PLUSContext) {
//	l.logger.Info("plus", "val", c.GetText())
//}
//
//// one of factor branch
//func (l *ebnfListener) ExitPLUS(c *ebnf.PLUSContext) {
//}
//
//// one of factor branch
//func (l *ebnfListener) EnterEXT(c *ebnf.EXTContext) {
//	l.logger.Info("ext", "val", c.GetText())
//}
//
//// one of factor branch
//func (l *ebnfListener) ExitEXT(c *ebnf.EXTContext) {
//}
//
//// one of factor branch
//func (l *ebnfListener) EnterSUB(c *ebnf.SUBContext) {
//	l.logger.Info("sub", "val", c.GetText())
//}
//
//// one of factor branch
//func (l *ebnfListener) ExitSUB(c *ebnf.SUBContext) {
//}

func Parse(file string) (*schemas.Grammar, error) {
	is, err := antlr.NewFileStream(file)
	if err != nil {
		return nil, err
	}
	lexer := ebnf.NewEBNFLexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	parser := ebnf.NewEBNFParser(stream)
	listener := newEbnfListener()
	antlr.ParseTreeWalkerDefault.Walk(listener, parser.Ebnf())
	return listener.grammar, nil
}

func MergeProduction(p map[string]*schemas.Node, startSymbol string) *schemas.Grammar {
	return nil
	//root, ok := p[startSymbol]
	//if !ok {
	//	return nil
	//}
	//
	//queue := make([]*schemas.Grammar, 0)
	//queue = append(queue, root)
	//visited := make([]*schemas.Grammar, 0)
	//for len(queue) != 0 {
	//	cur := queue[0]
	//	queue = queue[1:]
	//	if len(*cur.GetSymbols()) == 0 {
	//		if cur.GetType() == schemas.GrammarID {
	//			visited = append(visited, cur)
	//		}
	//	}
	//	for _, v := range *cur.GetSymbols() {
	//		queue = append(queue, v)
	//	}
	//}
	//for _, v := range visited {
	//	*v.GetSymbols() = append(*v.GetSymbols(), p[v.GetContent()])
	//}
	//return root
}
