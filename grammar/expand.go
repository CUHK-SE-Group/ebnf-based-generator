package grammar

import (
	"log"
	"math"
	"math/rand"
)

// ExpansionAlgo  第一个返回值是要拓展的symbol, 第二个返回值是拓展后的symbol
type ExpansionAlgo Handler
type Expansion interface {
	func(*DerivationContext, *DerivationTree)
}

func RandomExpand(ctx *DerivationContext, tree *DerivationTree) {
	nonTerminals := tree.GetNonTerminals()
	symbol := nonTerminals[rand.Int()%len(nonTerminals)]

	children := symbol.ExpansionTuple.Expand()
	symbol.Children = make([]*Node, 0)
	expPair := make([]ExpansionPair, 0)
	for _, v := range children {
		if IsNonTerminals(v.GetName()) {
			expansions := tree.gram.G[v.GetName()]
			symbol.Children = append(symbol.Children, &Node{ExpansionTuple: expansions[rand.Int()%len(expansions)]})
		} else {
			symbol.Children = append(symbol.Children, &Node{ExpansionTuple: v})
		}
		expPair = append(expPair, ExpansionPair{From: symbol.ExpansionTuple, To: v})
	}
	ctx.preExpansions = expPair
}

// MinimalCostExpand 找出非终端符号中成本最小的一个，并将它扩展为其孩子节点。
// Expand a non-terminal symbol in the derivation tree with the least cost.
func MinimalCostExpand(ctx *DerivationContext, tree *DerivationTree) {
	// 获取叶子节点上的所有的非终端符号
	nonTerminals := tree.GetNonTerminals() // []*Node
	var minCostSymbol *Node                // 用来记录成本最小的非终端符号
	var cost = math.Inf(1)                 // 初始成本设为一个较大的数值

	// 遍历所有非终端符号，找出成本最小的那个
	for _, symbol := range nonTerminals {
		children := symbol.ExpansionTuple.Expand() // 扩展当前符号, 当前符号有可能是 `<expr><expr>` 这种形式, 所以扩展后可能有多个孩子节点
		for _, child := range children {
			// 计算当前子节点的成本
			curCost := tree.gram.SymbolCost(child, make(map[string]struct{}), false)
			if cost > curCost {
				cost = curCost
				minCostSymbol = symbol
			}
		}
	}

	// 如果没有找到非终端符号，直接返回
	if minCostSymbol == nil {
		ctx.preExpansions = nil
		return
	}
	log.Default().Println("Choose minCostSymbol: ", minCostSymbol.GetName())

	// 开始扩展最小成本的非终端符号
	children := minCostSymbol.ExpansionTuple.Expand()
	log.Default().Printf("Expand children: %+v\n", children)
	minCostSymbol.Children = make([]*Node, 0)
	expPair := make([]ExpansionPair, 0) // 用来记录扩展的过程
	// 遍历所有子节点，如果是非终端节点则进一步扩展，否则直接添加为孩子节点
	for _, child := range children {
		if IsNonTerminals(child.GetName()) {
			minCost := math.Inf(1)
			var tmpChildren *Node
			expansions := tree.gram.G[child.GetName()]
			for _, expansion := range expansions { // 寻找要扩展的子节点中成本最小的那个
				cost := tree.gram.SymbolCost(expansion, make(map[string]struct{}), false)
				if minCost > cost {
					minCost = cost
					tmpChildren = &Node{ExpansionTuple: expansion}
				}
			}

			if child.GetName() != minCostSymbol.GetName() {
				minCostSymbol.Children = append(minCostSymbol.Children, &Node{ExpansionTuple: child, Children: []*Node{tmpChildren}})
				expPair = append(expPair, ExpansionPair{minCostSymbol.ExpansionTuple, child})
				expPair = append(expPair, ExpansionPair{child, tmpChildren.ExpansionTuple})
			} else {
				expPair = append(expPair, ExpansionPair{minCostSymbol.ExpansionTuple, tmpChildren.ExpansionTuple})
				minCostSymbol.Children = append(minCostSymbol.Children, tmpChildren)
			}
		} else {
			minCostSymbol.Children = append(minCostSymbol.Children, &Node{ExpansionTuple: child})
			expPair = append(expPair, ExpansionPair{minCostSymbol.ExpansionTuple, child})
		}
	}
	ctx.preExpansions = expPair
}

func MaximumCostExpand(ctx *DerivationContext, tree *DerivationTree) {
	nonTerminals := tree.GetNonTerminals()
	var maxCostSymbol *Node
	var cost = math.Inf(-1) // Initialize to negative infinity for maximum cost search

	for _, symbol := range nonTerminals {
		children := symbol.ExpansionTuple.Expand()
		for _, child := range children {
			curCost := tree.gram.SymbolCost(child, make(map[string]struct{}), true)
			if cost < curCost {
				cost = curCost
				maxCostSymbol = symbol
			}
		}
	}

	if maxCostSymbol == nil {
		ctx.preExpansions = nil
		return
	}
	log.Default().Println("Choose maxCostSymbol: ", maxCostSymbol.GetName())
	children := maxCostSymbol.ExpansionTuple.Expand()
	log.Default().Printf("Expand children: %+v\n", children)
	maxCostSymbol.Children = make([]*Node, 0)
	expPair := make([]ExpansionPair, 0)
	for _, child := range children {
		if IsNonTerminals(child.GetName()) {
			maxCost := math.Inf(-1)
			var tmpChildren *Node
			expansions := tree.gram.G[child.GetName()]
			for _, expansion := range expansions {
				cost := tree.gram.SymbolCost(expansion, make(map[string]struct{}), true) // Use 'true' for max cost
				if maxCost < cost {
					maxCost = cost
					tmpChildren = &Node{ExpansionTuple: expansion}
				}
			}

			if child.GetName() != maxCostSymbol.GetName() {
				maxCostSymbol.Children = append(maxCostSymbol.Children, &Node{ExpansionTuple: child, Children: []*Node{tmpChildren}})
				expPair = append(expPair, ExpansionPair{maxCostSymbol.ExpansionTuple, child})
				expPair = append(expPair, ExpansionPair{child, tmpChildren.ExpansionTuple})
			} else {
				expPair = append(expPair, ExpansionPair{maxCostSymbol.ExpansionTuple, tmpChildren.ExpansionTuple})
				maxCostSymbol.Children = append(maxCostSymbol.Children, tmpChildren)
			}
		} else {
			maxCostSymbol.Children = append(maxCostSymbol.Children, &Node{ExpansionTuple: child})
			expPair = append(expPair, ExpansionPair{maxCostSymbol.ExpansionTuple, child})
		}
	}
	ctx.preExpansions = expPair
}

// ThreePhaseExpansion https://www.fuzzingbook.org/html/GrammarFuzzer.html#Three-Expansion-Phases
func ThreePhaseExpansion(ctx *DerivationContext, tree *DerivationTree) {
}

// CoverageBasedExpasion https://www.fuzzingbook.org/html/GrammarCoverageFuzzer.html#Tracking-Grammar-Coverage
func CoverageBasedExpansion(ctx *DerivationContext, tree *DerivationTree) {
}
