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
