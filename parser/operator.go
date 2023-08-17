package parser

import (
	"github.com/CUHK-SE-Group/ebnf-based-generator/parser/ebnf"
	"github.com/antlr4-go/antlr/v4"
	"github.com/golang/glog"
	"github.com/lucasjones/reggen"
	"math/rand"
)

var systemOperators map[string]Operator

type Handler func(*Grammar, *Result)

type Operator interface {
	BeforeGen(*Grammar, *Result)
	Gen(*Grammar, *Result)
	AfterGen(*Grammar, *Result)
	GetText() string
}

type GenericOperator struct {
	BeforeGenHandlers []Handler
	GenHandler        Handler
	AfterGenHandlers  []Handler
	Text              string
}

func systemBeforeGen(g *Grammar, r *Result) {

}

func systemAfterGen(g *Grammar, r *Result) {
}

func (op *GenericOperator) BeforeGen(g *Grammar, r *Result) {
	for _, f := range op.BeforeGenHandlers {
		f(g, r)
	}
}

func (op *GenericOperator) Gen(g *Grammar, r *Result) {
	op.GenHandler(g, r)
}

func (op *GenericOperator) AfterGen(g *Grammar, r *Result) {
	for _, f := range op.AfterGenHandlers {
		f(g, r)
	}
}

func (op *GenericOperator) GetText() string {
	return op.Text
}

var Or = GenericOperator{
	BeforeGenHandlers: []Handler{systemBeforeGen},
	AfterGenHandlers:  []Handler{systemAfterGen},
	GenHandler:        OrGen,
	Text:              "Or",
}

func OrGen(g *Grammar, r *Result) {
	selected := (*g.Symbols)[rand.Int()%len(*g.Symbols)]
	r.AddNode(selected)
	selected.Generate(r)
}

var Rep = GenericOperator{
	BeforeGenHandlers: []Handler{systemBeforeGen},
	AfterGenHandlers:  []Handler{systemAfterGen},
	GenHandler:        RepGen,
	Text:              "Repeat",
}

func RepGen(g *Grammar, r *Result) {
	times := rand.Int() % 3
	selected := (*g.Symbols)[0]
	for i := 0; i < times; i++ {
		r.AddNode(selected)
		selected.Generate(r)
	}
}

var Plus = GenericOperator{
	BeforeGenHandlers: []Handler{systemBeforeGen},
	AfterGenHandlers:  []Handler{systemAfterGen},
	GenHandler:        PlusGen,
	Text:              "Plus",
}

func PlusGen(g *Grammar, r *Result) {
	times := rand.Int()%3 + 1
	selected := (*g.Symbols)[0]
	for i := 0; i < times; i++ {
		r.AddNode(selected)
		selected.Generate(r)
	}
}

var Ext = GenericOperator{
	BeforeGenHandlers: []Handler{systemBeforeGen},
	AfterGenHandlers:  []Handler{systemAfterGen},
	GenHandler:        ExtGen,
	Text:              "Exist",
}

func ExtGen(g *Grammar, r *Result) {
	ok := rand.Int() % 2
	selected := (*g.Symbols)[0]
	if ok == 1 {
		r.AddNode(selected)
		selected.Generate(r)
	}
}

var Cat = GenericOperator{
	BeforeGenHandlers: []Handler{systemBeforeGen},
	AfterGenHandlers:  []Handler{systemAfterGen},
	GenHandler:        CatGen,
	Text:              "Catenate",
}

func CatGen(g *Grammar, r *Result) {
	for _, selected := range *g.Symbols {
		r.AddNode(selected)
		selected.Generate(r)
	}
}

var Regex = GenericOperator{
	BeforeGenHandlers: []Handler{systemBeforeGen},
	AfterGenHandlers:  []Handler{systemAfterGen},
	GenHandler:        RegexGen,
	Text:              "Regex",
}

func RegexGen(g *Grammar, r *Result) {
	str, err := reggen.Generate(g.GetContent(), 10)
	//glog.Errorf("Generating: %s\n", g.GetContent())
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
func ParseGrammarDefinition(file string, startSymbol string, conf *Config) (*Grammar, error) {
	is, err := antlr.NewFileStream(file)
	if err != nil {
		return nil, err
	}
	lexer := ebnf.NewEBNFLexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	parser := ebnf.NewEBNFParser(stream)
	listener := newEbnfListener(conf)
	antlr.ParseTreeWalkerDefault.Walk(listener, parser.Ebnf())
	g := listener.productions[startSymbol]
	return g, nil
}
