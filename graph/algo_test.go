package graph

import (
	"fmt"
	"testing"
)

func newGraph() (Graph[string, string], []Vertex[string]) {
	// 构建测试图
	g := NewGraph[string, string]()
	// 添加节点和边以形成一个具有环的图
	// 例如，构建一个简单的图：0 -> 1 -> 2 -> 3 -> 4 -> 5 -> 6 -> 7 -> 8 -> 9
	//                         2->10->11->12->7
	//                         12->13->14
	//                           4->1
	vertices := make([]Vertex[string], 0)
	for i := 0; i < 15; i++ {
		v := NewVertex[string]()
		v.SetID(fmt.Sprintf("vertex%d", i))
		g.AddVertex(v)
		vertices = append(vertices, v)
	}

	// 添加边
	eCnt := 0
	for i := 0; i < 10; i++ {
		e := NewEdge[string, string]()
		e.SetID(fmt.Sprintf("edge%d", eCnt))
		e.SetFrom(vertices[i])
		e.SetTo(vertices[i+1])
		g.AddEdge(e)
		eCnt++
	}

	for i := 10; i < 14; i++ {
		e := NewEdge[string, string]()
		e.SetID(fmt.Sprintf("edge%d", eCnt))
		e.SetFrom(vertices[i])
		e.SetTo(vertices[i+1])
		g.AddEdge(e)
		eCnt++
	}
	// 添加形成环的边
	e1 := NewEdge[string, string]()
	e1.SetID(fmt.Sprintf("edge%d", eCnt))
	e1.SetFrom(vertices[4])
	e1.SetTo(vertices[1])
	g.AddEdge(e1)
	eCnt++

	e2 := NewEdge[string, string]()
	e2.SetID(fmt.Sprintf("edge%d", eCnt))
	e2.SetFrom(vertices[2])
	e2.SetTo(vertices[10])
	g.AddEdge(e2)
	eCnt++

	e3 := NewEdge[string, string]()
	e3.SetID(fmt.Sprintf("edge%d", eCnt))
	e3.SetFrom(vertices[12])
	e3.SetTo(vertices[7])
	g.AddEdge(e3)
	eCnt++

	return g, vertices
}
func TestIdentifyCyclesAndInitializeProbabilities(t *testing.T) {
	g, vertices := newGraph()
	// 调用函数
	identifyCyclesAndInitializeProbabilities[string, string](g, "vertex0")

	// 检查结果
	for i, v := range vertices {
		expectedProb := 0.0
		if i >= 1 && i <= 4 { // 节点 1, 2, 3 在环中
			expectedProb = 1.0
		}
		actualProb := v.GetMeta().(graphMeta).cycleProbability
		if actualProb != expectedProb {
			t.Errorf("Vertex %s: expected cycleProbability %v, got %v", v.GetID(), expectedProb, actualProb)
		}
	}
}
func TestUpdateNodeProbability(t *testing.T) {
	g, _ := newGraph()
	// 更新 node1 的概率
	identifyCyclesAndInitializeProbabilities(g, "vertex0")

	updateProbabilitiesUntilConvergence(g, 1000)

	fmt.Println(g)
}
