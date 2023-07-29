package grammar

import (
	"fmt"
	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
	"math"
	"math/rand"
	"os"
	"regexp"
)

const StartSymbol = "<start>"

var ReNonterminal = regexp.MustCompile("(<[^<> ]*>)")
var SimpleNonterminalGrammar = Grammar{
	"<start>":       []ExpansionTuple{{name: "<nonterminal>"}},
	"<nonterminal>": []ExpansionTuple{{name: "<left-angle><identifier><right-angle>"}},
	"<left-angle>":  []ExpansionTuple{{name: "<"}},
	"<right-angle>": []ExpansionTuple{{name: ">"}},
	"<identifier>":  []ExpansionTuple{{name: "id"}}, // for now
}

type Option map[string]interface{}
type ExpansionTuple struct {
	name string
	info Option
}

func (e *ExpansionTuple) Expand() []ExpansionTuple {
	if e.GetName() == "" {
		return []ExpansionTuple{{"", nil}}
	}
	re := regexp.MustCompile("<[^<> ]*>|[^<> ]+")
	results := re.FindAllString(e.GetName(), -1)
	var ExpansionList []ExpansionTuple
	for _, v := range results {
		ExpansionList = append(ExpansionList, ExpansionTuple{name: v})
	}
	return ExpansionList
}
func (e *ExpansionTuple) GetName() string {
	return e.name
}
func (e *ExpansionTuple) GetOpt() Option {
	return e.info
}

type Grammar map[string][]ExpansionTuple

func (grammar Grammar) ExpandSymbol(symbol string) []ExpansionTuple {
	expansions, ok := grammar[symbol]
	if !ok {
		return nil
	}
	newChildren := make([][]ExpansionTuple, 0)
	for _, v := range expansions {
		chs := v.Expand()
		tmp := make([]ExpansionTuple, 0)
		for _, ch := range chs {
			tmp = append(tmp, ch)
		}
		newChildren = append(newChildren, tmp)
	}
	chosenElement := newChildren[rand.Intn(len(newChildren))]
	return chosenElement
}
func (grammar Grammar) GetSymbol(symbol string) []ExpansionTuple {
	return grammar[symbol]
}

func (grammar Grammar) Extend(extension Grammar) {
	for k, v := range extension {
		grammar[k] = v
	}
}
func (grammar Grammar) OptsUsed() Set {
	usedOpts := make(Set)
	for _, expansions := range grammar {
		for _, expansion := range expansions {
			for opt := range expansion.GetOpt() {
				usedOpts[opt] = true
			}
		}
	}
	return usedOpts
}
func (grammar Grammar) SetOpts(symbol string, expainsion ExpansionTuple, opt Option) {
	expansions := grammar.GetSymbol(symbol)
	for i, v := range expansions {
		if v.GetName() != expainsion.GetName() {
			continue
		}
		newOpt := v.GetOpt()
		if newOpt == nil {
			newOpt = make(Option)

		} else {
			for k, nv := range opt {
				newOpt[k] = nv
			}
		}
		if newOpt != nil {
			grammar[symbol][i] = ExpansionTuple{name: expainsion.GetName(), info: newOpt}
		} else {
			grammar[symbol][i] = ExpansionTuple{name: expainsion.GetName()}
		}
	}
	return
}

func (grammar Grammar) defUsedNonterminals(startSymbol string) (Set, Set, error) {
	definedNonterminals := make(Set)
	usedNonterminals := make(Set)
	usedNonterminals[startSymbol] = true

	for definedNonterminal, expansions := range grammar {
		definedNonterminals[definedNonterminal] = true

		if len(expansions) == 0 {
			return nil, nil, fmt.Errorf("%s: expansion list empty", definedNonterminal)
		}

		for _, expansion := range expansions {
			for _, usedNonterminal := range NonTerminals(expansion.GetName()) {
				usedNonterminals[usedNonterminal] = true
			}
		}
	}

	return definedNonterminals, usedNonterminals, nil
}

func (grammar Grammar) reachableNonterminals(startSymbol string) map[string]struct{} {
	reachable := make(map[string]struct{})

	var findReachableNonterminals func(grammar Grammar, symbol string)
	findReachableNonterminals = func(grammar Grammar, symbol string) {
		reachable[symbol] = struct{}{}
		for _, expansion := range grammar[symbol] {
			for _, nonterminal := range NonTerminals(expansion.GetName()) {
				if _, ok := reachable[nonterminal]; !ok {
					findReachableNonterminals(grammar, nonterminal)
				}
			}
		}
	}

	findReachableNonterminals(grammar, startSymbol)
	return reachable
}

func (grammar Grammar) unreachableNonterminals(startSymbol string) map[string]struct{} {
	reachable := grammar.reachableNonterminals(startSymbol)
	unreachable := make(map[string]struct{})
	for key := range grammar {
		if _, ok := reachable[key]; !ok {
			unreachable[key] = struct{}{}
		}
	}
	return unreachable
}

func (grammar Grammar) IsValidGrammar(startSymbol string, supportedOpts Set) bool {
	if startSymbol == "" {
		startSymbol = StartSymbol
	}

	definedNonterminals, usedNonterminals, err := grammar.defUsedNonterminals(startSymbol)
	if err != nil {
		return false
	}

	if _, ok := grammar[StartSymbol]; ok {
		usedNonterminals[StartSymbol] = true
	}

	for unusedNonterminal := range SetDifference(definedNonterminals, usedNonterminals) {
		fmt.Fprintf(os.Stderr, "%q: defined, but not used\n", unusedNonterminal)
	}
	for undefinedNonterminal := range SetDifference(usedNonterminals, definedNonterminals) {
		fmt.Fprintf(os.Stderr, "%q: used, but not defined\n", undefinedNonterminal)
	}

	unreachable := grammar.unreachableNonterminals(startSymbol)
	msgStartSymbol := startSymbol

	if _, ok := grammar[StartSymbol]; ok {
		reachableFromStart := grammar.reachableNonterminals(StartSymbol)
		for terminal := range reachableFromStart {
			delete(unreachable, terminal)
		}
		if startSymbol != StartSymbol {
			msgStartSymbol += " or " + StartSymbol
		}
	}

	for unreachableNonterminal := range unreachable {
		fmt.Fprintf(os.Stderr, "%q: unreachable from %s\n", unreachableNonterminal, msgStartSymbol)
	}

	usedButNotSupportedOpts := NewSet(nil)
	if len(supportedOpts) > 0 {
		usedOpts := grammar.OptsUsed()
		usedButNotSupportedOpts = SetDifference(usedOpts, supportedOpts)
		for opt := range usedButNotSupportedOpts {
			fmt.Fprintf(os.Stderr, "warning: option %q is not supported\n", opt)
		}
	}

	return len(usedNonterminals) == len(definedNonterminals) && len(unreachable) == 0
}

func (grammar Grammar) Visualize(filename string) {
	g := graphviz.New()
	graph, _ := g.Graph()

	nodes := make(map[string]*cgraph.Node)

	// Create nodes
	for key := range grammar {
		node, _ := graph.CreateNode(key)
		nodes[key] = node
		if key != StartSymbol {
			node.SetColor("red")
		} else {
			node.SetColor("green")
		}
	}

	// Create edges
	for key, expansions := range grammar {
		for _, expansion := range expansions {

			target, exists := nodes[expansion.GetName()]
			if !exists {
				target, _ = graph.CreateNode(expansion.GetName())
				if IsNonTerminals(expansion.GetName()) {
					target.SetColor("red")
					for _, nonterminal := range NonTerminals(expansion.GetName()) {
						node, ok := nodes[nonterminal]
						if !ok {
							node, _ = graph.CreateNode(nonterminal)
							node.SetColor("red")
							nodes[nonterminal] = node
						}
						_, err := graph.CreateEdge("", target, node)
						if err != nil {
							panic(err)
						}

					}
				}
			}
			_, err := graph.CreateEdge("", nodes[key], target)
			if err != nil {
				panic(err)
			}
		}
	}

	// render the graph to a file
	err := g.RenderFilename(graph, graphviz.PNG, filename)
	if err != nil {
		panic(err)
	}
}

func (grammar Grammar) SymbolCost(symbol ExpansionTuple, seen map[string]struct{}) float64 {
	symbols := symbol.Expand()
	totalCost := 0.0
	for _, s := range symbols {
		expansions, ok := grammar[s.GetName()]
		if !ok {
			return 0 // terminal
		}
		minCost := math.Inf(1)
		for _, expansion := range expansions {
			seen[s.GetName()] = struct{}{}
			cost := grammar.ExpansionCost(expansion, seen)
			if cost < minCost {
				minCost = cost
			}
		}
		if minCost == math.Inf(1) {
			return minCost
		}
		totalCost += minCost
	}
	return totalCost
}

func (grammar Grammar) ExpansionCost(expansion ExpansionTuple, seen map[string]struct{}) float64 {
	symbols := NonTerminals(expansion.GetName())
	if len(symbols) == 0 {
		return 1 // no symbol
	}

	// Check if any symbol is in the seen set
	for _, s := range symbols {
		if _, ok := seen[s]; ok {
			return math.Inf(1)
		}
	}

	// Copy seen and add current symbol to the set
	newSeen := make(map[string]struct{})
	for k := range seen {
		newSeen[k] = struct{}{}
	}

	// The value of an expansion is the sum of all expandable variables
	// inside + 1
	total := 1.0
	for _, s := range symbols {
		total += grammar.SymbolCost(ExpansionTuple{name: s}, newSeen)
	}
	return total
}
