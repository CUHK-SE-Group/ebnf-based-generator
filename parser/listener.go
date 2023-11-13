package parser

import (
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
	productions       map[string]*schemas.Node
	logger            *slog.Logger
	grammar           *schemas.Grammar
}

func newEbnfListener() *ebnfListener {
	textHandler := slog.NewTextHandler(os.Stderr, nil)
	callerInfoHandler := log.NewCallerInfoHandler(textHandler)
	logger := slog.New(callerInfoHandler)

	listener := &ebnfListener{
		currentSymbolId:   0,
		currentProduction: &schemas.Node{},
		productions:       map[string]*schemas.Node{},
		stack:             []*schemas.Node{},
		logger:            logger,
		grammar:           schemas.NewGrammar(),
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

func (l *ebnfListener) save(g *schemas.Node) {
	l.productions[g.GetID()] = g
}

func (l *ebnfListener) clear() {
	l.stack = []*schemas.Node{}
}

func (l *ebnfListener) EnterProduction(c *ebnf.ProductionContext) {

}
func (l *ebnfListener) ExitProduction(c *ebnf.ProductionContext) {
	l.pop()
}

func (l *ebnfListener) EnterExpr(c *ebnf.ExprContext) {

}
func (l *ebnfListener) ExitExpr(c *ebnf.ExprContext) {
	l.pop()
}

func (l *ebnfListener) EnterTerm(c *ebnf.TermContext) {

}

func (l *ebnfListener) ExitTerm(c *ebnf.TermContext) {
	l.pop()
}

// EnterID one of factor branch
func (l *ebnfListener) EnterID(c *ebnf.IDContext) {

}
func (l *ebnfListener) ExitID(c *ebnf.IDContext) {
	l.pop()
}

// EnterQUOTE one of factor branch
func (l *ebnfListener) EnterQUOTE(c *ebnf.QUOTEContext) {

}
func (l *ebnfListener) ExitQUOTE(c *ebnf.QUOTEContext) {
	l.pop()
}

// one of factor branch
func (l *ebnfListener) EnterPAREN(c *ebnf.PARENContext) {

}
func (l *ebnfListener) ExitPAREN(c *ebnf.PARENContext) {
	l.pop()
}

// one of factor branch
func (l *ebnfListener) EnterBRACKET(c *ebnf.BRACKETContext) {

}
func (l *ebnfListener) ExitBRACKET(c *ebnf.BRACKETContext) {
	l.pop()
}

// one of factor branch
func (l *ebnfListener) EnterBRACE(c *ebnf.BRACEContext) {

}
func (l *ebnfListener) ExitBRACE(c *ebnf.BRACEContext) {
	l.pop()
}

// one of factor branch
func (l *ebnfListener) EnterREP(c *ebnf.REPContext) {

}

// one of factor branch
func (l *ebnfListener) ExitREP(c *ebnf.REPContext) {
	l.pop()
}

// one of factor branch
func (l *ebnfListener) EnterPLUS(c *ebnf.PLUSContext) {

}

// one of factor branch
func (l *ebnfListener) ExitPLUS(c *ebnf.PLUSContext) {
	l.pop()
}

// one of factor branch
func (l *ebnfListener) EnterEXT(c *ebnf.EXTContext) {

}

// one of factor branch
func (l *ebnfListener) ExitEXT(c *ebnf.EXTContext) {
	l.pop()
}

// one of factor branch
func (l *ebnfListener) EnterSUB(c *ebnf.SUBContext) {

}

// one of factor branch
func (l *ebnfListener) ExitSUB(c *ebnf.SUBContext) {
	l.pop()
}

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
