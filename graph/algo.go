package graph

import "fmt"

type graphMeta struct {
	cycleProbability float64
	edgeProbability  float64
}

func identifyCyclesAndInitializeProbabilities[T1 any, T2 any](graph Graph[T1, T2], startId string) {
	vertices := graph.GetAllVertices()
	startSym := graph.GetVertexById(startId)
	visited := make(map[string]bool)
	recStack := make(map[string]bool)
	path := make(map[string]string) // 用于记录每个节点的前驱节点

	var dfs func(Vertex[T2], string) bool
	dfs = func(v Vertex[T2], parent string) bool {
		fmt.Println(v.GetID())
		visited[v.GetID()] = true
		recStack[v.GetID()] = true
		path[v.GetID()] = parent

		for _, edge := range graph.GetOutEdges(v) {
			to := edge.GetTo()
			if !visited[to.GetID()] {
				if dfs(to, v.GetID()) {
					return true
				}
			} else if recStack[to.GetID()] {
				// 发现环，回溯路径直到环的起点
				currentNode := v.GetID()
				for currentNode != to.GetID() && currentNode != "" {
					n := graph.GetVertexById(currentNode)
					n.SetMeta(graphMeta{cycleProbability: 1.0})
					currentNode = path[currentNode]
				}
				graph.GetVertexById(to.GetID()).SetMeta(graphMeta{cycleProbability: 1.0})
				return true
			}
		}

		recStack[v.GetID()] = false
		return false
	}

	// 初始化所有节点的概率为0
	for _, v := range vertices {
		v.SetMeta(graphMeta{cycleProbability: 0.0})
	}

	dfs(startSym, "")
}

// 更新单个节点的到达环概率
func updateNodeProbability[T1 any, T2 any](graph Graph[T1, T2], node Vertex[T2]) float64 {
	outEdges := graph.GetOutEdges(node)
	probability := 0.0

	for _, edge := range outEdges {
		targetNode := edge.GetTo()
		if edge.GetMeta() == nil {
			edge.SetMeta(graphMeta{})
		}
		edgeProbability := edge.GetMeta().(graphMeta).edgeProbability
		if edgeProbability == 0 {
			edgeProbability = float64(1.0 / len(outEdges))
		}
		targetProbability := targetNode.GetMeta().(graphMeta).cycleProbability
		probability += edgeProbability * targetProbability
	}
	ori := node.GetMeta().(graphMeta)
	ori.cycleProbability = probability
	node.SetMeta(ori)
	return probability
}
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// 迭代更新所有节点的概率，直到收敛或达到最大迭代次数
func updateProbabilitiesUntilConvergence[T1 any, T2 any](graph Graph[T1, T2], maxIterations int) {
	const Threshold = 0.0001

	for i := 0; i < maxIterations; i++ {
		var totalChange float64 = 0
		vertices := graph.GetAllVertices()

		for _, node := range vertices {
			oldProbability := node.GetMeta().(graphMeta).cycleProbability
			newProbability := updateNodeProbability(graph, node)
			totalChange += abs(newProbability - oldProbability)
		}

		if totalChange < Threshold {
			break
		}
	}
}
