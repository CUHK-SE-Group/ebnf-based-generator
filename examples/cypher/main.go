package main

import (
	"fmt"
	"github.com/CUHK-SE-Group/ebnf-based-generator/parser"
	"log"
	"math/rand"
	"os"
	"runtime/pprof"
	"strings"
)

func OrGen(g *parser.Grammar, r *parser.Result) {
	fmt.Printf("traversaling %v\n", g.ID)
	var selected *parser.Grammar
	maxChoice := 2 * len(*g.Symbols)
	for i := 0; i < maxChoice; i++ {
		selected = (*g.Symbols)[rand.Int()%len(*g.Symbols)]
		sym := strings.Split(selected.ID, "#")[0]
		if cur, ok := g.Ctx.SymCount[sym]; ok {
			if limit, ex := g.Config.SymbolLimit[sym]; ex {
				if cur < limit {
					r.AddNode(selected)
					selected.Generate(r)
					g.Ctx.SymCount[sym]++
					return
				}
			}
		} else {
			r.AddNode(selected)
			selected.Generate(r)
			g.Ctx.SymCount[sym]++
		}
	}
}
func main() {
	f, err := os.Create("examples/cypher/cpu.prof")
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	defer f.Close()
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
	defer pprof.StopCPUProfile()
	parser.Init()

	var Or = parser.GenericOperator{
		GenHandler: OrGen,
		Text:       "Or",
		BeforeGenHandlers: []parser.Handler{func(g *parser.Grammar, r *parser.Result) {

		}},
	}
	parser.RegisterOperator(&Or)

	conf := parser.Config{SymbolLimit: map[string]int{
		"oC_Match":                10,
		"oC_Unwind":               10,
		"oC_ComparisonExpression": 3,
	}}
	grammar, err := parser.ParseGrammarDefinition("examples/cypher/cypher.ebnf", "oC_Cypher", &conf)
	if err != nil {
		panic(err)
	}
	fmt.Println(grammar.ID)

	//grammar.Visualize("examples/cypher/tmp.png", true)
	res := grammar.Generate(nil)
	res.Visualize("./tmp.png")
	for _, s := range res.GetOutput() {
		fmt.Printf("%v", s)
	}
}
