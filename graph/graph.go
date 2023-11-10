package graph

import "fmt"

type Metadata string

type Graph[T any] interface {
	AddEdge(edge Edge[T])
	AddVertex(vertex Vertex[T])
	DeleteEdge(edge Edge[T])
	DeleteVertex(vertex Vertex[T])
	GetOutEdges(vertex Vertex[T]) []Edge[T]
	GetInEdges(vertex Vertex[T]) []Edge[T]
	GetAllVertices() []Vertex[T]
	GetAllEdges() []Edge[T]
	SetMetadata(key Metadata, val bool)
	GetMetadata(key Metadata) bool
	GetAllMetadata() map[Metadata]bool
	GetVertexById(id string) Vertex[T]
	GetEdgeById(id string) Edge[T]
}

type Edge[T any] interface {
	SetID(id string)
	SetFrom(vertex Vertex[T])
	SetTo(vertex Vertex[T])
	SetProperty(key string, val T)
	GetID() string
	GetFrom() Vertex[T]
	GetTo() Vertex[T]
	GetProperty(key string) T
	GetAllProperties() map[string]T
}

type Vertex[T any] interface {
	SetID(id string)
	SetProperty(key string, val T)
	GetID() string
	GetProperty(key string) T
	GetAllProperties() map[string]T
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
