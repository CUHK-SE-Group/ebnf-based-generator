package main

import (
	"fmt"
	"github.com/CUHK-SE-Group/generic-generator/schemas"
	"log"
	"math/rand"
	"os"
	"runtime/pprof"
	"strings"
)

func OrGen(g *schemas.Grammar, r *schemas.Result) {
	fmt.Printf("traversaling %v\n", g.id)
	var selected *schemas.Grammar
	maxChoice := 2 * len(*g.symbols)
	for i := 0; i < maxChoice; i++ {
		selected = (*g.symbols)[rand.Int()%len(*g.symbols)]
		sym := strings.Split(selected.id, "#")[0]
		if cur, ok := g.Ctx.SymCount[sym]; ok {
			if limit, ex := g.confgi.SymbolLimit[sym]; ex {
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
	schemas.Init()

	var Or = schemas.GenericOperator{
		GenHandler: OrGen,
		Text:       "Or",
		BeforeGenHandlers: []schemas.Handler{func(g *schemas.Grammar, r *schemas.Result) {

		}},
	}
	schemas.RegisterOperator(&Or)

	conf := schemas.Config{SymbolLimit: map[string]int{
		"oC_Match":                10,
		"oC_Unwind":               10,
		"oC_ComparisonExpression": 3,
	}}
	grammar, err := schemas.ParseGrammarDefinition("examples/cypher/cypher.ebnf", "oC_Cypher", &conf)
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
