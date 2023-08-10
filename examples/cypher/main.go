package main

import (
	"fmt"
	"github.com/CUHK-SE-Group/ebnf-based-generator/parser"
	"log"
	"os"
	"runtime/pprof"
)

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
	grammar, err := parser.ParseGrammarDefinition("examples/cypher/cypher.ebnf", "oC_Cypher")
	if err != nil {
		panic(err)
	}
	fmt.Println(grammar.ID)

	grammar.Visualize("examples/cypher/tmp.png", true)
	res := grammar.Generate(nil)
	res.Visualize("./tmp.png")
}
