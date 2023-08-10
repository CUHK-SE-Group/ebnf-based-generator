package parser

import (
	"fmt"
	"math/rand"

	"github.com/golang/glog"
	"github.com/lucasjones/reggen"
)

type handler func(*Context, *Grammar, *Result)

type Operator interface {
	BeforeGen(*Context, *Grammar, *Result)
	Gen(*Context, *Grammar, *Result)
	AfterGen(*Context, *Grammar, *Result)
	GetText() string
}

type GenricOperator struct {
	beforeGenHanders []handler
	genHandler       handler
	afterGenHandlers []handler
	text             string
}

func systemBeforeGen(ctx *Context, g *Grammar, r *Result) {
	fmt.Println("do something before gen")
}

func systemAfterGen(ctx *Context, g *Grammar, r *Result) {
	fmt.Println("do something after gen")
}

func (op *GenricOperator) BeforeGen(ctx *Context, g *Grammar, r *Result) {
	for _, f := range op.beforeGenHanders {
		f(ctx, g, r)
	}
}

func (op *GenricOperator) Gen(ctx *Context, g *Grammar, r *Result) {
	op.genHandler(ctx, g, r)
}

func (op *GenricOperator) AfterGen(ctx *Context, g *Grammar, r *Result) {
	for _, f := range op.afterGenHandlers {
		f(ctx, g, r)
	}
}

func (op *GenricOperator) GetText() string {
	return op.text
}

var Or = GenricOperator{
	beforeGenHanders: []handler{systemBeforeGen},
	afterGenHandlers: []handler{systemAfterGen},
	genHandler:       OrGen,
	text:             "Or",
}

func OrGen(ctx *Context, g *Grammar, r *Result) {
	selected := (*g.Symbols)[rand.Int()%len(*g.Symbols)]
	r.AddNode(selected)
	selected.Generate(r)
}

var Rep = GenricOperator{
	beforeGenHanders: []handler{systemBeforeGen},
	afterGenHandlers: []handler{systemAfterGen},
	genHandler:       RepGen,
	text:             "Repeat",
}

func RepGen(ctx *Context, g *Grammar, r *Result) {
	times := rand.Int() % 100
	selected := (*g.Symbols)[0]
	for i := 0; i < times; i++ {
		r.AddNode(selected)
		selected.Generate(r)
	}
}

var Ext = GenricOperator{
	beforeGenHanders: []handler{systemBeforeGen},
	afterGenHandlers: []handler{systemAfterGen},
	genHandler:       ExtGen,
	text:             "Exist",
}

func ExtGen(ctx *Context, g *Grammar, r *Result) {
	ok := rand.Int() % 2
	selected := (*g.Symbols)[0]
	if ok == 1 {
		r.AddNode(selected)
		selected.Generate(r)
	}
}

var Cat = GenricOperator{
	beforeGenHanders: []handler{systemBeforeGen},
	afterGenHandlers: []handler{systemAfterGen},
	genHandler:       CatGen,
	text:             "Catenate",
}

func CatGen(ctx *Context, g *Grammar, r *Result) {
	for _, selected := range *g.Symbols {
		r.AddNode(selected)
		selected.Generate(r)
	}
}

var Regex = GenricOperator{
	beforeGenHanders: []handler{systemBeforeGen},
	afterGenHandlers: []handler{systemAfterGen},
	genHandler:       RegexGen,
	text:             "Regex",
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
