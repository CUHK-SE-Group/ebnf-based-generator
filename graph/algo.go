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

func TarjanSCC[EdgePropertyType any, VertexPropertyType any](graph Graph[EdgePropertyType, VertexPropertyType]) (map[string]string, Graph[EdgePropertyType, VertexPropertyType]) {
	var (
		index    int = 0
		stack    []*VertexImpl[VertexPropertyType]
		inStack  = make(map[string]bool)
		indices  = make(map[string]int)
		lowLinks = make(map[string]int)
		sccMap   = make(map[string]string) // Map From original vertex ID To SCC representative ID
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
