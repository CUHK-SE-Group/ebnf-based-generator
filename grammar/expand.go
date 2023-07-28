package grammar

import "math/rand"

type ExpansionAlgo func(*DerivationTree)
type Expansion interface {
	func(*Node, func(node *Node))
}

func RandomExpand(tree *DerivationTree) {
	nonTerminals := tree.GetNonTerminals()
	symbol := nonTerminals[rand.Int()%len(nonTerminals)]

	children := symbol.ExpansionTuple.Expand()
	symbol.Children = make([]*Node, 0)
	for _, v := range children {
		if IsNonTerminals(v.GetName()) {
			tmpChildren := make([]*Node, 0)
			expansions := tree.gram[v.GetName()]
			for _, expansion := range expansions {
				tmpChildren = append(tmpChildren, &Node{ExpansionTuple: expansion})
			}
			symbol.Children = append(symbol.Children, &Node{ExpansionTuple: v, Children: tmpChildren})
		} else {
			symbol.Children = append(symbol.Children, &Node{ExpansionTuple: v})
		}
	}
}
func MinimalCostExpand(tree *DerivationTree) {
	nonTerminals := tree.GetNonTerminals()
	var minCostSymbol *Node
	var cost = 100000.0

	for _, v := range nonTerminals {
		children := v.ExpansionTuple.Expand()
		for _, child := range children {
			m := make(map[string]struct{})
			curCost := tree.gram.SymbolCost(child.GetName(), m)
			if cost > curCost {
				cost = curCost
				minCostSymbol = v
			}
		}
	}
	if minCostSymbol == nil {
		return
	}
	symbol := minCostSymbol

	children := symbol.ExpansionTuple.Expand()
	symbol.Children = make([]*Node, 0)
	for _, v := range children {
		if IsNonTerminals(v.GetName()) {
			tmpChildren := make([]*Node, 0)
			expansions := tree.gram[v.GetName()]
			for _, expansion := range expansions {
				tmpChildren = append(tmpChildren, &Node{ExpansionTuple: expansion})
			}
			if v.GetName() != symbol.GetName() {
				symbol.Children = append(symbol.Children, &Node{ExpansionTuple: v, Children: tmpChildren})
			} else {
				for _, child := range tmpChildren {
					symbol.Children = append(symbol.Children, child)
				}
			}

		} else {
			symbol.Children = append(symbol.Children, &Node{ExpansionTuple: v})
		}
	}
}
