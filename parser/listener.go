package parser

import (
	"fmt"
	"os"
	"strconv"

	"github.com/CUHK-SE-Group/ebnf-based-generator/parser/ebnf"

	"github.com/golang/glog"
)

type ebnfListener struct {
	*ebnf.BaseEBNFParserListener
	context           *Context
	stack             []*Grammar
	currentSymbolId   int
	currentProduction *Grammar
	productions       map[string]*Grammar
	forkList          map[*Grammar][]int
	config            *Config
}

func newEbnfListener(conf *Config) *ebnfListener {
	listener := &ebnfListener{
		context:           NewContext(),
		currentSymbolId:   0,
		currentProduction: &Grammar{},
		productions:       map[string]*Grammar{},
		stack:             []*Grammar{},
		forkList:          map[*Grammar][]int{},
		config:            conf,
	}
	return listener
}

func (l *ebnfListener) addToForkList(g *Grammar, id int) {
	l.forkList[g] = append(l.forkList[g], id)
}

func (l *ebnfListener) forkAll() {
	for g, idxList := range l.forkList {
		for _, idx := range idxList {
			// id := g.GetSymbol(idx).GetID()
			fork, err := g.GetSymbol(idx).ForkContext("")
			if err != nil {
				glog.Fatalf("non-production grammar in forkList: %s", g.GetID())
			}
			g.SetSymbol(idx, fork)
		}
	}
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

func (l *ebnfListener) getCurrentCtx() *Context {
	if l.currentProduction.Ctx == nil {
		return NewContext()
	}
	return l.currentProduction.Ctx
}

func (l *ebnfListener) EnterProduction(ctx *ebnf.ProductionContext) {
	l.currentSymbolId = 0
	name := ctx.ID().GetText()
	if p, ok := l.productions[name]; ok {
		// if the production has been parsed, then ...
		l.currentProduction = p
		l.clear()
		l.push(p)
	} else {
		// a new grammar tree
		new := NewGrammar(l.getCurrentCtx(), GrammarProduction, systemOperators["Catenate"], name, "", l.config)
		// save the production and clean the stack for a new round
		l.currentProduction = new
		l.clear()
		l.save(new)
		l.push(new)
	}
}

func (l *ebnfListener) EnterSymbolWithUOp(ctx *ebnf.SymbolWithUOpContext) {
	var op Operator
	switch ctx.UnaryOp().GetStart().GetTokenType() {
	case ebnf.EBNFLexerREP:
		op = systemOperators["Repeat"]
	case ebnf.EBNFLexerEXT:
		op = systemOperators["Exist"]
	case ebnf.EBNFLexerPLUS:
		op = systemOperators["Plus"]
	}
	// generate a new node
	expr := NewGrammar(l.getCurrentCtx(), GrammarInner, op, l.generateId(), "", l.config)
	// add node to its parent
	l.top().AddSymbol(expr)
	// traverse children
	l.push(expr)
	fmt.Fprintf(os.Stderr, "UOP the symbol is %s, push %s\n", ctx.GetText(), expr.ID)
}

func (l *ebnfListener) ExitSymbolWithUOp(ctx *ebnf.SymbolWithUOpContext) {
	l.pop()
}

var pushed []bool = make([]bool, 0)

func (l *ebnfListener) EnterSymbolWithBOp(ctx *ebnf.SymbolWithBOpContext) {
	var op Operator
	switch ctx.BinaryOp().GetStart().GetTokenType() {
	case ebnf.EBNFLexerOR:
		op = systemOperators["Or"]
	}
	if l.top().GetOperator() != op {
		expr := NewGrammar(l.getCurrentCtx(), GrammarInner, op, l.generateId(), "", l.config)
		l.top().AddSymbol(expr)
		l.push(expr)
		pushed = append(pushed, true)
		fmt.Fprintf(os.Stderr, "BOP the symbol is %s, left is %s, right is %s, push %s\n", ctx.GetText(), ctx.Expr(0).GetText(), ctx.Expr(1).GetText(), expr.ID)
	} else {
		pushed = append(pushed, false)
	}
}

func (l *ebnfListener) ExitSymbolWithBOp(ctx *ebnf.SymbolWithBOpContext) {
	if pushed[len(pushed)-1] == true {
		l.pop()
	}
	pushed = pushed[:len(pushed)-1]
}

func (l *ebnfListener) EnterSubSymbol(ctx *ebnf.SubSymbolContext) {
	fmt.Fprintf(os.Stderr, "SUBSYM the symbol is %s\n", ctx.GetText())
	expr := NewGrammar(l.getCurrentCtx(), GrammarInner, systemOperators["Catenate"], l.generateId(), "", l.config)
	// add the current node to its parent
	l.top().AddSymbol(expr)
	// then it is time to traverse this tree's subtrees.
	l.push(expr)
}

func (l *ebnfListener) ExitSubSymbol(ctx *ebnf.SubSymbolContext) {
	l.pop()
}

func (l *ebnfListener) EnterSymbolWithCat(ctx *ebnf.SymbolWithCatContext) {
	fmt.Fprintf(os.Stderr, "CAT the symbol is %s\n", ctx.GetText())
	// new a node that express the current node
	expr := NewGrammar(l.getCurrentCtx(), GrammarInner, systemOperators["Catenate"], l.generateId(), "", l.config)
	// add current node to its parent
	l.top().AddSymbol(expr)
	// traverse its children
	l.push(expr)
}

func (l *ebnfListener) ExitSymbolWithCat(ctx *ebnf.SymbolWithCatContext) {
	l.pop()
}

func (l *ebnfListener) EnterSubProduction(ctx *ebnf.SubProductionContext) {
	proName := ctx.ID().GetText()
	fmt.Fprintf(os.Stderr, "SUBP the symbol is %s\n", proName)
	if g, ok := l.productions[proName]; ok {
		new, err := g.ForkContext(l.generateId() + "#" + proName)
		if err != nil {
			glog.Errorf("failed to bind proudction to existing definition for %s.", g.ID)
			glog.Errorf("some non-production grammars may be present in the symbol table.")
			glog.Fatal("Terminating...")
		}
		l.top().AddSymbol(new)
	} else {
		new := NewGrammar(l.getCurrentCtx(), GrammarProduction, systemOperators["Catenate"], proName, "", l.config)
		placeholder := new.ShallowCopy().SetID(l.generateId() + "#" + proName)
		l.productions[proName] = new
		id := l.top().AddSymbol(placeholder)
		l.addToForkList(l.top(), id)
	}
}

func (l *ebnfListener) ExitEbnf(ctx *ebnf.EbnfContext) {
	l.forkAll()
}

func (l *ebnfListener) EnterTerminal(ctx *ebnf.TerminalContext) {
	fmt.Fprintf(os.Stderr, "TER the symbol is %s\n", ctx.GetText())
	expr := NewGrammar(l.getCurrentCtx(), GrammarTerminal, systemOperators["Regex"], l.generateId()+"#"+ctx.GetText(), ctx.GetText(), l.config)
	l.top().AddSymbol(expr)
}
