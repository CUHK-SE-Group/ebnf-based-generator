package grammar

import (
	"math"
	"math/rand"
)

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

// MinimalCostExpand 找出非终端符号中成本最小的一个，并将它扩展为其孩子节点。
func MinimalCostExpand(tree *DerivationTree) {
	// 获取叶子节点上的所有的非终端符号
	nonTerminals := tree.GetNonTerminals()
	var minCostSymbol *Node // 用来记录成本最小的非终端符号
	var cost = 100000.0     // 初始成本设为一个较大的数值

	// 遍历所有非终端符号，找出成本最小的那个
	for _, symbol := range nonTerminals {
		children := symbol.ExpansionTuple.Expand() // 扩展当前符号, 当前符号有可能是 `<expr><expr>` 这种形式, 所以扩展后可能有多个孩子节点
		for _, child := range children {
			m := make(map[string]struct{})
			// 计算当前子节点的成本
			curCost := tree.gram.SymbolCost(child, m)
			if cost > curCost {
				cost = curCost
				minCostSymbol = symbol
			}
		}
	}

	// 如果没有找到非终端符号，直接返回
	if minCostSymbol == nil {
		return
	}
	//log.Default().Println("Choose minCostSymbol: ", minCostSymbol.GetName())
	// 开始扩展最小成本的非终端符号
	children := minCostSymbol.ExpansionTuple.Expand()
	//log.Default().Printf("Expand children: %+v\n", children)
	minCostSymbol.Children = make([]*Node, 0)

	// 遍历所有子节点，如果是非终端节点则进一步扩展，否则直接添加为孩子节点
	for _, child := range children {
		if IsNonTerminals(child.GetName()) {
			minCost := math.Inf(1)
			var tmpChildren *Node
			expansions := tree.gram[child.GetName()]
			for _, expansion := range expansions { // 寻找要扩展的子节点中成本最小的那个
				cost := tree.gram.SymbolCost(expansion, make(map[string]struct{}))
				if minCost > cost {
					minCost = cost
					tmpChildren = &Node{ExpansionTuple: expansion}
				}
			}

			// 若子节点的名字与当前符号的名字不同，则将其及其子节点一起添加为孩子
			// 否则，将子节点的子节点添加为孩子
			if child.GetName() != minCostSymbol.GetName() {
				minCostSymbol.Children = append(minCostSymbol.Children, &Node{ExpansionTuple: child, Children: []*Node{tmpChildren}})
			} else {
				minCostSymbol.Children = append(minCostSymbol.Children, tmpChildren)
			}
		} else {
			// 如果是终端节点，直接添加为孩子
			minCostSymbol.Children = append(minCostSymbol.Children, &Node{ExpansionTuple: child})
		}
	}
}
