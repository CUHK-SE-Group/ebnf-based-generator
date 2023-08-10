package parser

import (
	"fmt"
	"github.com/CUHK-SE-Group/ebnf-based-generator/parser/ebnf"
	"github.com/antlr4-go/antlr/v4"
	"math/rand"

	"github.com/golang/glog"
	"github.com/lucasjones/reggen"
)

var systemOperators map[string]Operator

type handler func(*Context, *Grammar, *Result)

type Operator interface {
	BeforeGen(*Context, *Grammar, *Result)
	Gen(*Context, *Grammar, *Result)
	AfterGen(*Context, *Grammar, *Result)
	GetText() string
}

type GenericOperator struct {
	beforeGenHandlers []handler
	genHandler        handler
	afterGenHandlers  []handler
	text              string
}

func systemBeforeGen(ctx *Context, g *Grammar, r *Result) {
	fmt.Println("do something before gen")
}

func systemAfterGen(ctx *Context, g *Grammar, r *Result) {
	fmt.Println("do something after gen")
}

func (op *GenericOperator) BeforeGen(ctx *Context, g *Grammar, r *Result) {
	for _, f := range op.beforeGenHandlers {
		f(ctx, g, r)
	}
}

func (op *GenericOperator) Gen(ctx *Context, g *Grammar, r *Result) {
	op.genHandler(ctx, g, r)
}

func (op *GenericOperator) AfterGen(ctx *Context, g *Grammar, r *Result) {
	for _, f := range op.afterGenHandlers {
		f(ctx, g, r)
	}
}

func (op *GenericOperator) GetText() string {
	return op.text
}

var Or = GenericOperator{
	beforeGenHandlers: []handler{systemBeforeGen},
	afterGenHandlers:  []handler{systemAfterGen},
	genHandler:        OrGen,
	text:              "Or",
}

func OrGen(ctx *Context, g *Grammar, r *Result) {
	selected := (*g.Symbols)[rand.Int()%len(*g.Symbols)]
	r.AddNode(selected)
	selected.Generate(r)
}

var Rep = GenericOperator{
	beforeGenHandlers: []handler{systemBeforeGen},
	afterGenHandlers:  []handler{systemAfterGen},
	genHandler:        RepGen,
	text:              "Repeat",
}

func RepGen(ctx *Context, g *Grammar, r *Result) {
	times := rand.Int() % 100
	selected := (*g.Symbols)[0]
	for i := 0; i < times; i++ {
		r.AddNode(selected)
		selected.Generate(r)
	}
}

var Plus = GenericOperator{
	beforeGenHandlers: []handler{systemBeforeGen},
	afterGenHandlers:  []handler{systemAfterGen},
	genHandler:        PlusGen,
	text:              "Plus",
}

func PlusGen(ctx *Context, g *Grammar, r *Result) {
	times := rand.Int()%99 + 1
	selected := (*g.Symbols)[0]
	for i := 0; i < times; i++ {
		r.AddNode(selected)
		selected.Generate(r)
	}
}

var Ext = GenericOperator{
	beforeGenHandlers: []handler{systemBeforeGen},
	afterGenHandlers:  []handler{systemAfterGen},
	genHandler:        ExtGen,
	text:              "Exist",
}

func ExtGen(ctx *Context, g *Grammar, r *Result) {
	ok := rand.Int() % 2
	selected := (*g.Symbols)[0]
	if ok == 1 {
		r.AddNode(selected)
		selected.Generate(r)
	}
}

var Cat = GenericOperator{
	beforeGenHandlers: []handler{systemBeforeGen},
	afterGenHandlers:  []handler{systemAfterGen},
	genHandler:        CatGen,
	text:              "Catenate",
}

func CatGen(ctx *Context, g *Grammar, r *Result) {
	for _, selected := range *g.Symbols {
		r.AddNode(selected)
		selected.Generate(r)
	}
}

var Regex = GenericOperator{
	beforeGenHandlers: []handler{systemBeforeGen},
	afterGenHandlers:  []handler{systemAfterGen},
	genHandler:        RegexGen,
	text:              "Regex",
}

func RegexGen(ctx *Context, g *Grammar, r *Result) {
	str, err := reggen.Generate(g.GetContent(), 100)
	glog.Errorf("Generating: %s\n", g.GetContent())
	if err != nil {
		glog.Warningf("can not generate regex: %s", g.GetContent())
	}
	// r.AddNode(g)
	r.AddOutput(str)
}
func Init() {
	systemOperators = make(map[string]Operator)
	RegisterOperator(&Or)
	RegisterOperator(&Rep)
	RegisterOperator(&Ext)
	RegisterOperator(&Cat)
	RegisterOperator(&Regex)
	RegisterOperator(&Plus)
}
func RegisterOperator(op Operator) {
	if op.GetText() == "" {
		glog.Fatalf("Operator must have a name")
	}
	systemOperators[op.GetText()] = op
}
func ParseGrammarDefinition(file string, startSymbol string) (*Grammar, error) {
	is, err := antlr.NewFileStream(file)
	if err != nil {
		return nil, err
	}
	lexer := ebnf.NewEBNFLexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	parser := ebnf.NewEBNFParser(stream)
	listener := newEbnfListener()
	antlr.ParseTreeWalkerDefault.Walk(listener, parser.Ebnf())
	g := listener.productions[startSymbol]
	return g, nil
}
