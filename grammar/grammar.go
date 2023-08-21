package grammar

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"regexp"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
)

const StartSymbol = "<start>"

// A regular expression pattern to identify non-terminal symbols.
var ReNonterminal = regexp.MustCompile("(<[^<> ]*>)")
var SimpleNonterminalGrammar = Grammar{G: map[string][]ExpansionTuple{
	"<start>":       {{name: "<nonterminal>"}},
	"<nonterminal>": {{name: "<left-angle><identifier><right-angle>"}},
	"<left-angle>":  {{name: "<"}},
	"<right-angle>": {{name: ">"}},
	"<identifier>":  {{name: "id"}}, // for now
}}

type Option map[string]interface{}

// ExpansionTuple represents a specific expansion in the grammar;
// [name]: actual expansion like <num> + <num>;
// [info]: additional information;
type ExpansionTuple struct {
	name string
	info Option
}

// Decompose a complex expansion into its individual parts;
// For example, if the name of an ExpansionTuple is <expr> + <term>, the Expand method would break it down into three ExpansionTuple objects: one for <expr>, one for +, and one for <term>.
func (e *ExpansionTuple) Expand() []ExpansionTuple {
	if e.GetName() == "" {
		return []ExpansionTuple{{"", nil}}
	}

	// identify and extract all non-terminals
	re := regexp.MustCompile("<[^<> ]*>|[^<> ]+")
	results := re.FindAllString(e.GetName(), -1)

	var ExpansionList []ExpansionTuple
	for _, v := range results {
		// only set the name field
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

type Grammar struct {
	// The main grammar map.
	G map[string][]ExpansionTuple
	// A cache for symbol costs.
	symbolCache map[string]float64
	// A cache for expansion costs.
	expansionCache map[string]float64
}

// Returns a random expansion of a given symbol.
func (grammar *Grammar) ExpandSymbol(symbol string) []ExpansionTuple {
	expansions, ok := grammar.G[symbol] // expansions: a list of ExpansionTuple
	if !ok {
		return nil
	}
	newChildren := make([][]ExpansionTuple, 0)
	for _, v := range expansions {
		// <num> + <num> => [<num>, +, <num>]
		chs := v.Expand() // break one expansion into seperated ExpansionTuples
		tmp := make([]ExpansionTuple, 0)
		tmp = append(tmp, chs...)
		newChildren = append(newChildren, tmp)
	}
	// randomly select one expansion
	chosenElement := newChildren[rand.Intn(len(newChildren))]
	return chosenElement
}

// Returns the expansions associated with a symbol.
func (grammar *Grammar) GetSymbol(symbol string) []ExpansionTuple {
	return grammar.G[symbol]
}

// Extends the grammar with another grammar.
func (grammar *Grammar) Extend(extension Grammar) {
	for k, v := range extension.G {
		grammar.G[k] = v
	}
}

// Returns a set of used options.
func (grammar *Grammar) OptsUsed() Set {
	usedOpts := make(Set)
	for _, expansions := range grammar.G {
		for _, expansion := range expansions {
			for opt := range expansion.GetOpt() {
				usedOpts[opt] = true
			}
		}
	}
	return usedOpts
}

// Sets options for a specific symbol and expansion;
// If the expansion already has some options set, the new options are merged with the existing ones;
func (grammar *Grammar) SetOpts(symbol string, expainsion ExpansionTuple, opt Option) {
	expansions := grammar.GetSymbol(symbol) // type: []ExpansionTuple
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
			grammar.G[symbol][i] = ExpansionTuple{name: expainsion.GetName(), info: newOpt}
		} else {
			grammar.G[symbol][i] = ExpansionTuple{name: expainsion.GetName()}
		}
	}
}

// Return: definedNonterminals, usedNonterminals, error
func (grammar *Grammar) defUsedNonterminals(startSymbol string) (Set, Set, error) {
	definedNonterminals := make(Set)
	usedNonterminals := make(Set)
	usedNonterminals[startSymbol] = true

	for definedNonterminal, expansions := range grammar.G {
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

// Returns non-terminals that are reachable from a given start symbol.
func (grammar *Grammar) reachableNonterminals(startSymbol string) map[string]struct{} {
	reachable := make(map[string]struct{})

	var findReachableNonterminals func(grammar *Grammar, symbol string)
	findReachableNonterminals = func(grammar *Grammar, symbol string) {
		reachable[symbol] = struct{}{}
		for _, expansion := range grammar.G[symbol] {
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

// Returns non-terminals that are not reachable from a given start symbol.
func (grammar *Grammar) unreachableNonterminals(startSymbol string) map[string]struct{} {
	reachable := grammar.reachableNonterminals(startSymbol)
	unreachable := make(map[string]struct{})
	for key := range grammar.G {
		if _, ok := reachable[key]; !ok {
			unreachable[key] = struct{}{}
		}
	}
	return unreachable
}

func (grammar *Grammar) IsValidGrammar(startSymbol string, supportedOpts Set) bool {
	if startSymbol == "" {
		startSymbol = StartSymbol
	}

	definedNonterminals, usedNonterminals, err := grammar.defUsedNonterminals(startSymbol)
	if err != nil {
		return false
	}

	if _, ok := grammar.G[StartSymbol]; ok {
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

	if _, ok := grammar.G[StartSymbol]; ok {
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

	usedButNotSupportedOpts := NewSet(nil) // Warning: value not used
	if len(supportedOpts) > 0 {
		usedOpts := grammar.OptsUsed()
		usedButNotSupportedOpts = SetDifference(usedOpts, supportedOpts)
		for opt := range usedButNotSupportedOpts {
			fmt.Fprintf(os.Stderr, "warning: option %q is not supported\n", opt)
		}
	}

	return len(usedNonterminals) == len(definedNonterminals) && len(unreachable) == 0
}

// Visualizes the grammar and saves it as an image.
func (grammar *Grammar) Visualize(filename string) {
	g := graphviz.New()
	graph, _ := g.Graph()

	nodes := make(map[string]*cgraph.Node)

	// Create nodes
	for key := range grammar.G {
		node, _ := graph.CreateNode(key)
		nodes[key] = node
		if key != StartSymbol {
			node.SetColor("red")
		} else {
			node.SetColor("green")
		}
	}

	// Create edges
	for key, expansions := range grammar.G {
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

// The next two methods recursively call each other:
// SymbolCost is the cost of a symbol cost, to calculate, you must know
// the cost of every expansion of this symbol, so it calls ExpansionCost
// For ExpansionCost, it sums up the cost of symbols in this expansion
// Edge/Reduce case:
// 1. The symbol is a terminal / Expansion contains no non-terminals
// 2. The symbol has been seen -> infinite recursion -> cost = +Inf

// Computes the cost of expanding a symbol;
// The method can compute either the minimum or maximum cost based on the max argument.
func (grammar *Grammar) SymbolCost(symbol ExpansionTuple, seen map[string]struct{}, max bool) float64 {
	// Cache the cost of a symbol
	if grammar.symbolCache == nil {
		grammar.symbolCache = make(map[string]float64)
	}
	// Has been cached
	if cost, ok := grammar.symbolCache[symbol.GetName()]; ok {
		return cost
	}
	// Get all symbols in this ExpansionTuple
	symbols := symbol.Expand()
	totalCost := 0.0
	for _, s := range symbols {
		// The total cost of a symbol is the sum of the costs of its constituent symbols.
		// For each symbol, expand using the strategy and calculate ExpansionCost
		expansions, ok := grammar.G[s.GetName()]
		if !ok {
			return 0
		}
		minCost := math.Inf(1)
		maxCost := math.Inf(-1)
		for _, expansion := range expansions {
			seen[s.GetName()] = struct{}{}
			cost := grammar.ExpansionCost(expansion, seen, max)
			if cost < minCost {
				minCost = cost
			}
			if cost > maxCost {
				maxCost = cost
			}
		}
		if max {
			if maxCost == math.Inf(-1) {
				return maxCost
			}
			totalCost += maxCost
		}
		if !max {
			if minCost == math.Inf(1) {
				return minCost
			}
			totalCost += minCost
		}
	}
	grammar.symbolCache[symbol.GetName()] = totalCost
	return totalCost
}

// Computes the cost of an expansion.
func (grammar *Grammar) ExpansionCost(expansion ExpansionTuple, seen map[string]struct{}, max bool) float64 {
	if grammar.expansionCache == nil {
		grammar.expansionCache = make(map[string]float64)
	}
	if cost, ok := grammar.expansionCache[expansion.GetName()]; ok {
		return cost
	}

	//  Identifies all non-terminals in the expansion.
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
		total += grammar.SymbolCost(ExpansionTuple{name: s}, newSeen, max)
	}
	grammar.expansionCache[expansion.GetName()] = total
	return total
}
