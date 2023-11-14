package graph

import (
	"fmt"
	"io"
	"os"
	"text/template"
)

//const dotTemplate = `strict {{.GraphType}} {
//{{range $k, $v := .Attributes}}
//	{{$k}}="{{$v}}";
//{{end}}
//{{range $s := .Statements}}
//	"{{.Source}}" {{if .Target}}{{$.EdgeOperator}} "{{.Target}}" [ {{range $k, $v := .EdgeAttributes}}{{$k}}="{{$v}}", {{end}} weight={{.EdgeWeight}} ]{{else}}[ {{range $k, $v := .SourceAttributes}}{{$k}}="{{$v}}", {{end}} weight={{.SourceWeight}} ]{{end}};
//{{end}}
//}
//`

const dotTemplate = `strict {{.GraphType}} {
{{range $k, $v := .Attributes}}
	{{$k}}="{{$v}}";
{{end}}
{{range $s := .Statements}}
	"{{.Source}}" {{if .Target}}{{$.EdgeOperator}} "{{.Target}}" [ label="{{.EdgeLabel}}", weight={{.EdgeWeight}} ]{{else}}[ label="{{.VertexLabel}}", weight={{.SourceWeight}} ]{{end}};
{{end}}
}
`

type Metadata string

type Graph[PropertyType any] interface {
	AddEdge(edge Edge[PropertyType])
	AddVertex(vertex Vertex[PropertyType])
	DeleteEdge(edge Edge[PropertyType])
	DeleteVertex(vertex Vertex[PropertyType])
	GetOutEdges(vertex Vertex[PropertyType]) []Edge[PropertyType]
	GetInEdges(vertex Vertex[PropertyType]) []Edge[PropertyType]
	GetAllVertices() []Vertex[PropertyType]
	GetAllEdges() []Edge[PropertyType]
	SetMetadata(key Metadata, val any)
	GetMetadata(key Metadata) any
	GetAllMetadata() map[Metadata]any
	GetVertexById(id string) Vertex[PropertyType]
	GetEdgeById(id string) Edge[PropertyType]
}

type Edge[PropertyType any] interface {
	SetID(id string)
	SetFrom(vertex Vertex[PropertyType])
	SetTo(vertex Vertex[PropertyType])
	SetProperty(key string, val PropertyType)
	GetID() string
	GetFrom() Vertex[PropertyType]
	GetTo() Vertex[PropertyType]
	GetProperty(key string) PropertyType
	GetAllProperties() map[string]PropertyType
}

type Vertex[PropertyType any] interface {
	SetID(id string)
	SetProperty(key string, val PropertyType)
	GetID() string
	GetProperty(key string) PropertyType
	GetAllProperties() map[string]PropertyType
}

type GraphSafe[T any] interface {
	AddEdgeSafe(edge EdgeSafe[T])
	AddVertexSafe(vertex VertexSafe[T])
	DeleteEdgeSafe(edge EdgeSafe[T])
	DeleteVertexSafe(vertex VertexSafe[T])
	GetOutEdgesSafe(vertex VertexSafe[T]) []EdgeSafe[T]
	GetInEdgesSafe(vertex VertexSafe[T]) []EdgeSafe[T]
	GetAllVerticesSafe() []VertexSafe[T]
	GetAllEdgesSafe() []EdgeSafe[T]
}

type EdgeSafe[T any] interface {
	SetIDSafe(id string)
	SetFromSafe(vertex VertexSafe[T])
	SetToSafe(vertex VertexSafe[T])
	SetPropertySafe(key string, val T)
	GetIDSafe() string
	GetFromSafe() VertexSafe[T]
	GetToSafe() VertexSafe[T]
	GetPropertySafe(key string, val T) T
	GetAllProperties() map[string]T
}

type VertexSafe[T any] interface {
	SetIDSafe(id string)
	SetPropertySafe(key string, val T)
	GetIDSafe() string
	GetPropertySafe(key string, val T) T
	GetAllProperties() map[string]T
}

func Clone[T any](graph Graph[T], newGraph func() Graph[T], newEdge func() Edge[T], newVertex func() Vertex[T]) Graph[T] {
	// Use the provided factory function to create a new graph instance
	clonedGraph := newGraph()
	for k, v := range graph.GetAllMetadata() {
		clonedGraph.SetMetadata(k, v)
	}

	// Create a map to track the mapping from original vertices to cloned vertices
	vertexMap := make(map[string]Vertex[T])

	// Clone all vertices
	for _, v := range graph.GetAllVertices() {
		clonedVertex := newVertex() // Use the factory function to create a new vertex instance
		clonedVertex.SetID(v.GetID())
		// Retrieve and set all properties
		for key, val := range v.GetAllProperties() {
			clonedVertex.SetProperty(key, val)
		}
		// Add to the new graph and update the map
		clonedGraph.AddVertex(clonedVertex)
		vertexMap[v.GetID()] = clonedVertex
	}

	// Clone all edges
	for _, e := range graph.GetAllEdges() {
		clonedEdge := newEdge() // Use the factory function to create a new edge instance
		clonedEdge.SetID(e.GetID())
		// Set the start and end points, using the map to find the corresponding cloned vertices
		clonedEdge.SetFrom(vertexMap[e.GetFrom().GetID()])
		clonedEdge.SetTo(vertexMap[e.GetTo().GetID()])
		// Retrieve and set all properties
		for key, val := range e.GetAllProperties() {
			clonedEdge.SetProperty(key, val)
		}
		if clonedEdge.GetFrom() == nil || clonedEdge.GetTo() == nil {
			fmt.Println("nil")
		}
		// Add to the new graph
		clonedGraph.AddEdge(clonedEdge)
	}

	// Return the cloned graph
	return clonedGraph
}

func Visualize[PropertyType any](graph Graph[PropertyType], filename string, f func(vertex Vertex[PropertyType]) string) error {
	desc, err := generateDOT(graph, f)
	if err != nil {
		return fmt.Errorf("failed to generate DOT description: %w", err)
	}
	w, _ := os.Create(filename)
	return renderDOT(w, desc)
}

type description[PropertyType any] struct {
	GraphType    string
	Attributes   map[string]string
	EdgeOperator string
	Statements   []statement[PropertyType]
}

type statement[PropertyType any] struct {
	Source       interface{}
	Target       interface{}
	SourceWeight int
	EdgeLabel    string
	EdgeWeight   int
	VertexLabel  string
}

func generateDOT[PropertyType any](g Graph[PropertyType], f func(node Vertex[PropertyType]) string) (description[PropertyType], error) {
	desc := description[PropertyType]{
		GraphType:    "graph",
		Attributes:   make(map[string]string),
		EdgeOperator: "--",
		Statements:   make([]statement[PropertyType], 0),
	}
	if f == nil {
		f = func(node Vertex[PropertyType]) string {
			return node.GetID()
		}
	}

	desc.GraphType = "digraph"
	desc.EdgeOperator = "->"

	for _, vertex := range g.GetAllVertices() {
		stmt := statement[PropertyType]{
			Source:       vertex.GetID(),
			SourceWeight: 1,
			VertexLabel:  f(vertex),
		}
		desc.Statements = append(desc.Statements, stmt)

		for _, edge := range g.GetOutEdges(vertex) {
			stmt1 := statement[PropertyType]{
				Source:     vertex.GetID(),
				Target:     edge.GetTo().GetID(),
				EdgeWeight: 1,
				//EdgeLabel:  f(edge),
			}
			desc.Statements = append(desc.Statements, stmt1)
		}
	}

	return desc, nil
}

func renderDOT[PropertyType any](w io.Writer, d description[PropertyType]) error {
	tpl, err := template.New("dotTemplate").Parse(dotTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	return tpl.Execute(w, d)
}
