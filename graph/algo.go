package graph

import (
	"fmt"
	"math"
	"time"
)

type graphMeta struct {
	cycleProbability float64
	edgeProbability  float64
}

func identifyCyclesAndInitializeProbabilities[T1 any, T2 any](graph Graph[T1, T2], startId string) {
	vertices := graph.GetAllVertices()
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

	for _, v := range vertices {
		visited = make(map[string]bool)
		recStack = make(map[string]bool)
		path = make(map[string]string)
		dfs(v, "")
	}
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

func TarjanSCC[EdgePropertyType any, VertexPropertyType any](graph Graph[EdgePropertyType, VertexPropertyType]) (map[string]string, Graph[EdgePropertyType, VertexPropertyType]) {
	var (
		index    int = 0
		stack    []*VertexImpl[VertexPropertyType]
		inStack  = make(map[string]bool)
		indices  = make(map[string]int)
		lowLinks = make(map[string]int)
		sccMap   = make(map[string]string) // Map from original vertex ID to SCC representative ID
	)

	var strongconnect func(v *VertexImpl[VertexPropertyType])
	strongconnect = func(v *VertexImpl[VertexPropertyType]) {
		indices[v.GetID()] = index
		lowLinks[v.GetID()] = index
		index++
		stack = append(stack, v)
		inStack[v.GetID()] = true

		for _, edge := range graph.GetOutEdges(v) {
			w := edge.GetTo().(*VertexImpl[VertexPropertyType])
			if _, ok := indices[w.GetID()]; !ok {
				strongconnect(w)
				lowLinks[v.GetID()] = min(lowLinks[v.GetID()], lowLinks[w.GetID()])
			} else if inStack[w.GetID()] {
				lowLinks[v.GetID()] = min(lowLinks[v.GetID()], indices[w.GetID()])
			}
		}

		if lowLinks[v.GetID()] == indices[v.GetID()] {
			var scc []string
			for {
				w := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				inStack[w.GetID()] = false
				scc = append(scc, w.GetID())
				sccMap[w.GetID()] = v.GetID()
				if w.GetID() == v.GetID() {
					break
				}
			}
		}
	}

	for _, v := range graph.GetAllVertices() {
		if _, ok := indices[v.GetID()]; !ok {
			strongconnect(v.(*VertexImpl[VertexPropertyType]))
		}
	}

	newGraph := NewGraph[EdgePropertyType, VertexPropertyType]()

	for _, v := range graph.GetAllVertices() {
		sccID := sccMap[v.GetID()]
		if newGraph.GetVertexById(sccID) == nil {
			newV := NewVertex[VertexPropertyType]()
			newV.SetID(sccID)
			newGraph.AddVertex(newV)
		}
	}

	for _, e := range graph.GetAllEdges() {
		fromSCC := sccMap[e.GetFrom().GetID()]
		toSCC := sccMap[e.GetTo().GetID()]
		if fromSCC != toSCC {
			newE := NewEdge[EdgePropertyType, VertexPropertyType]()
			newE.SetFrom(newGraph.GetVertexById(fromSCC))
			newE.SetTo(newGraph.GetVertexById(toSCC))
			newE.SetID(fmt.Sprintf("%s_%s", fromSCC, toSCC))
			newGraph.AddEdge(newE)
		}
	}

	return sccMap, newGraph
}
func FloydAlgorithm[EdgePropertyType any, VertexPropertyType any](graph Graph[EdgePropertyType, VertexPropertyType]) map[string]map[string]float64 {
	t1 := time.Now()
	defer func() {
		duration := time.Since(t1)
		fmt.Println(duration)
	}()
	vertices := graph.GetAllVertices()
	numVertices := len(vertices)
	vertexMap := make(map[string]int) // 用于映射顶点 ID 到其索引
	for i, vertex := range vertices {
		vertexMap[vertex.GetID()] = i
	}

	// 初始化距离矩阵
	dist := make([][]float64, numVertices)
	for i := range dist {
		dist[i] = make([]float64, numVertices)
		for j := range dist[i] {
			if i == j {
				dist[i][j] = 0
			} else {
				dist[i][j] = math.Inf(1)
			}
		}
	}

	// 设置初始边的权重
	for _, edge := range graph.GetAllEdges() {
		fromID := edge.GetFrom().GetID()
		toID := edge.GetTo().GetID()
		weight := 1.0 // 假设边的权重属性名为 "weight"
		dist[vertexMap[fromID]][vertexMap[toID]] = weight
	}

	// 应用 Floyd 算法
	for k := 0; k < numVertices; k++ {
		for i := 0; i < numVertices; i++ {
			for j := 0; j < numVertices; j++ {
				if dist[i][k]+dist[k][j] < dist[i][j] {
					dist[i][j] = dist[i][k] + dist[k][j]
				}
			}
		}
	}

	// 转换为 map[string]map[string]float64
	distMap := make(map[string]map[string]float64)
	for i, vi := range vertices {
		distMap[vi.GetID()] = make(map[string]float64)
		for j, vj := range vertices {
			distMap[vi.GetID()][vj.GetID()] = dist[i][j]
		}
	}

	return distMap
}
