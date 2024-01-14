package graph

import (
	"encoding/gob"
	"fmt"
	"math/rand"
	"reflect"
	"sort"
	"testing"
)

func sortVerticesByID[A interface{ GetID() string }](arr []A) {
	sort.Slice(arr, func(i, j int) bool {
		return arr[i].GetID() < arr[j].GetID()
	})
}

func isEqualByID[A interface{ GetID() string }](arr1, arr2 []A) bool {
	if len(arr1) != len(arr2) {
		return false
	}

	v1 := make([]A, len(arr1))
	v2 := make([]A, len(arr2))
	copy(v1, arr1)
	copy(v2, arr2)

	sortVerticesByID(v1)
	sortVerticesByID(v2)

	for i := range v1 {
		if v1[i].GetID() != v2[i].GetID() {
			return false
		}
	}
	return true
}

func map2slice[T1 comparable, T2 any](ori map[T1]T2) []T2 {
	res := make([]T2, 0)
	for _, v := range ori {
		res = append(res, v)
	}
	return res
}
func TestNewGraph(t *testing.T) {
	var g Graph[string, string]
	vertices := make(map[string]Vertex[string])
	edges := make(map[string]Edge[string, string])
	g = NewGraph[string, string]()
	for i := 0; i < 10; i++ {
		e := NewEdge[string, string]()
		e.SetID(fmt.Sprintf("edge%d", i))

		from := NewVertex[string]()
		from.SetID(fmt.Sprintf("vertex%d", rand.Intn(10)))
		e.SetFrom(from)

		to := NewVertex[string]()
		to.SetID(fmt.Sprintf("vertex%d", rand.Intn(10)))
		e.SetTo(to)

		vertices[from.GetID()] = from
		vertices[to.GetID()] = to
		edges[e.GetID()] = e
		g.AddEdge(e)
	}
	validate := func(graph Graph[string, string]) {
		resV := graph.GetAllVertices()
		resE := graph.GetAllEdges()
		oriV := map2slice(vertices)
		oriE := map2slice(edges)
		if !isEqualByID(resV, oriV) {
			t.Error("vertex not equal")
		}
		if !isEqualByID(resE, oriE) {
			t.Error("vertex not equal")
		}
	}
	validate(g)
	newG := Clone(g, NewGraph[string, string], NewEdge[string, string], NewVertex[string])
	validate(newG)
	Visualize(g, "./file.dot", nil, nil)
}

type GraphTestCase[T1 any, T2 any] struct {
	name     string
	setup    func(g Graph[T1, T2])      // 设置测试环境的函数
	test     func(g Graph[T1, T2])      // 实际测试执行的函数
	expected func(g Graph[T1, T2]) bool // 验证结果的函数
}

func TestGraph(t *testing.T) {
	tests := []GraphTestCase[string, string]{
		{
			name: "Add and Get Vertex",
			setup: func(g Graph[string, string]) {
				v := NewVertex[string]()
				v.SetID("v1")
				g.AddVertex(v)
			},
			test: func(g Graph[string, string]) {
				v := g.GetAllVertices()
				if len(v) != 1 {
					t.Errorf("expected 1 vertex, got %d", len(v))
				}
			},
			expected: func(g Graph[string, string]) bool {
				return len(g.GetAllVertices()) == 1
			},
		},
		{
			name: "Add and Get Edge",
			setup: func(g Graph[string, string]) {
				v1 := NewVertex[string]()
				v1.SetID("v1")
				v2 := NewVertex[string]()
				v2.SetID("v2")
				g.AddVertex(v1)
				g.AddVertex(v2)
				e := NewEdge[string, string]()
				e.SetID("e1")
				e.SetFrom(v1)
				e.SetTo(v2)
				g.AddEdge(e)
			},
			test: func(g Graph[string, string]) {
				e := g.GetAllEdges()
				if len(e) != 1 {
					t.Errorf("expected 1 edge, got %d", len(e))
				}
			},
			expected: func(g Graph[string, string]) bool {
				return len(g.GetAllEdges()) == 1
			},
		},
		{
			name: "Delete Vertex",
			setup: func(g Graph[string, string]) {
				v := NewVertex[string]()
				v.SetID("v1")
				g.AddVertex(v)
			},
			test: func(g Graph[string, string]) {
				v := NewVertex[string]()
				v.SetID("v1")
				g.DeleteVertex(v)
			},
			expected: func(g Graph[string, string]) bool {
				return len(g.GetAllVertices()) == 0
			},
		},
		{
			name: "Delete Edge",
			setup: func(g Graph[string, string]) {
				v1 := NewVertex[string]()
				v1.SetID("v1")
				v2 := NewVertex[string]()
				v2.SetID("v2")
				g.AddVertex(v1)
				g.AddVertex(v2)
				e := NewEdge[string, string]()
				e.SetID("e1")
				e.SetFrom(v1)
				e.SetTo(v2)
				g.AddEdge(e)
			},
			test: func(g Graph[string, string]) {
				e := NewEdge[string, string]()
				e.SetID("e1")
				g.DeleteEdge(e)
			},
			expected: func(g Graph[string, string]) bool {
				return len(g.GetAllEdges()) == 0
			},
		},
		{
			name: "Set and Get Vertex Property",
			setup: func(g Graph[string, string]) {
				v := NewVertex[string]()
				v.SetID("v1")
				v.SetProperty("color", "blue")
				g.AddVertex(v)
			},
			test: func(g Graph[string, string]) {
				v := g.GetVertexById("v1")
				color := v.GetProperty("color")
				if color != "blue" {
					t.Errorf("expected vertex property color To be 'blue', got '%s'", color)
				}
			},
			expected: func(g Graph[string, string]) bool {
				v := g.GetVertexById("v1")
				p := v.GetProperty("color")
				return p == "blue"
			},
		},
		{
			name: "Set and Get Edge Property",
			setup: func(g Graph[string, string]) {
				v1 := NewVertex[string]()
				v1.SetID("v1")
				v2 := NewVertex[string]()
				v2.SetID("v2")
				g.AddVertex(v1)
				g.AddVertex(v2)
				e := NewEdge[string, string]()
				e.SetID("e1")
				e.SetFrom(v1)
				e.SetTo(v2)
				e.SetProperty("weight", "10")
				g.AddEdge(e)
			},
			test: func(g Graph[string, string]) {
				e := g.GetEdgeById("e1")
				weight := e.GetProperty("weight")
				if weight != "10" {
					t.Errorf("expected edge property weight To be '10', got '%s'", weight)
				}
			},
			expected: func(g Graph[string, string]) bool {
				e := g.GetEdgeById("e1")
				p := e.GetProperty("weight")
				return p == "10"
			},
		},
		{
			name: "Set and Get Metadata",
			setup: func(g Graph[string, string]) {
				g.SetMetadata("sorted", true)
			},
			test: func(g Graph[string, string]) {
				sorted := g.GetMetadata("sorted").(bool)
				if !sorted {
					t.Errorf("expected Metadata sorted To be true, got false")
				}
			},
			expected: func(g Graph[string, string]) bool {
				return g.GetMetadata("sorted").(bool)
			},
		},
		{
			name: "Auto-delete Vertex When No Edges Remain",
			setup: func(g Graph[string, string]) {
				g.SetMetadata(CleanVertexByEdge, true)
				v := NewVertex[string]()
				v.SetID("v1")
				g.AddVertex(v)

				v2 := NewVertex[string]()
				v2.SetID("v2")
				e := NewEdge[string, string]()
				e.SetID("e1")
				e.SetFrom(v)
				e.SetTo(v2)
				g.AddEdge(e)
			},
			test: func(g Graph[string, string]) {
				e := NewEdge[string, string]()
				e.SetID("e1")
				g.DeleteEdge(e)
			},
			expected: func(g Graph[string, string]) bool {
				vertices := g.GetAllVertices()
				for _, v := range vertices {
					if v.GetID() == "v1" {
						return false
					}
				}
				return true
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := NewGraph[string, string]()
			tc.setup(g)          // 设置测试环境
			tc.test(g)           // 执行测试
			if !tc.expected(g) { // 验证结果
				t.Errorf("test %s failed", tc.name)
			}
		})
	}
}

func TestFSDump(t *testing.T) {
	gob.Register(&FSGraph[int, int]{})
	gob.Register(&FSVertexImpl[int]{})
	gob.Register(&FSEdgeImpl[int, int]{})

	graph := NewGraph[int, int](WithPersistent(true))
	v1 := NewVertex[int]()
	v1.SetID("v1")

	v2 := NewVertex[int]()
	v2.SetID("v2")

	e1 := NewEdge[int, int]()
	e1.SetID("e1")
	e1.SetFrom(v1)
	e1.SetTo(v2)

	graph.AddEdge(e1)

	err := graph.Save("/tmp/graph")
	if err != nil {
		t.Error(err)
	}

	newG, err := load[int, int]("/tmp/graph")
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(graph, newG) {
		t.Errorf("反序列化后的图与原始图不相同")
	}
}
