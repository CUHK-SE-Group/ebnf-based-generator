package graph

import (
	"log/slog"
)

const (
	CleanVertexByEdge Metadata = "cleanVertexWhenNoEdge"
	GraphFile         Metadata = "graphFile"
)

type MemGraph[EdgePropertyType any, VertexPropertyType any] struct {
	edgeMap   map[string]Edge[EdgePropertyType, VertexPropertyType]
	vertexMap map[string]Vertex[VertexPropertyType]
	metadata  map[Metadata]any

	// index
	vertex2OutEdges map[string][]Edge[EdgePropertyType, VertexPropertyType]
	vertex2InEdges  map[string][]Edge[EdgePropertyType, VertexPropertyType]
	dirty           bool
}

func (g *MemGraph[EdgePropertyType, VertexPropertyType]) updateIndex() {
	in := func(m map[string][]Edge[EdgePropertyType, VertexPropertyType], key string) bool {
		_, ok := m[key]
		return ok
	}
	if g.dirty {
		if g.metadata[CleanVertexByEdge].(bool) && len(g.edgeMap) == 0 {
			g.vertexMap = make(map[string]Vertex[VertexPropertyType])
		}
		g.vertex2OutEdges = make(map[string][]Edge[EdgePropertyType, VertexPropertyType])
		g.vertex2InEdges = make(map[string][]Edge[EdgePropertyType, VertexPropertyType])
		for _, e := range g.edgeMap {
			if _, ok := g.vertex2InEdges[e.GetTo().GetID()]; !ok {
				g.vertex2InEdges[e.GetTo().GetID()] = make([]Edge[EdgePropertyType, VertexPropertyType], 0)
			}
			if _, ok := g.vertex2OutEdges[e.GetFrom().GetID()]; !ok {
				g.vertex2OutEdges[e.GetFrom().GetID()] = make([]Edge[EdgePropertyType, VertexPropertyType], 0)
			}
			g.vertex2InEdges[e.GetTo().GetID()] = append(g.vertex2InEdges[e.GetTo().GetID()], e)
			g.vertex2OutEdges[e.GetFrom().GetID()] = append(g.vertex2OutEdges[e.GetFrom().GetID()], e)

			if g.metadata[CleanVertexByEdge].(bool) {
				// clean single vertex
				if !in(g.vertex2InEdges, e.GetFrom().GetID()) && !in(g.vertex2OutEdges, e.GetFrom().GetID()) {
					delete(g.vertexMap, e.GetFrom().GetID())
				}
				if !in(g.vertex2InEdges, e.GetTo().GetID()) && !in(g.vertex2OutEdges, e.GetTo().GetID()) {
					delete(g.vertexMap, e.GetTo().GetID())
				}
			}
		}

		g.dirty = false
	}
}

func (g *MemGraph[EdgePropertyType, VertexPropertyType]) SetMetadata(key Metadata, val any) {
	g.metadata[key] = val
}
func (g *MemGraph[EdgePropertyType, VertexPropertyType]) GetMetadata(key Metadata) any {
	return g.metadata[key]
}

func (g *MemGraph[EdgePropertyType, VertexPropertyType]) GetAllMetadata() map[Metadata]any {
	return g.metadata
}

func (g *MemGraph[EdgePropertyType, VertexPropertyType]) GetVertexById(id string) Vertex[VertexPropertyType] {
	return g.vertexMap[id]
}
func (g *MemGraph[EdgePropertyType, VertexPropertyType]) GetEdgeById(id string) Edge[EdgePropertyType, VertexPropertyType] {
	return g.edgeMap[id]
}
func (g *MemGraph[EdgePropertyType, VertexPropertyType]) AddEdge(edge Edge[EdgePropertyType, VertexPropertyType]) {
	if _, ok := g.edgeMap[edge.GetID()]; ok {
		slog.Debug("edge already exists", "Id", edge.GetID(), "From", edge.GetFrom().GetID(), "To", edge.GetTo().GetID())
	}
	g.edgeMap[edge.GetID()] = edge
	if _, ok := g.vertexMap[edge.GetFrom().GetID()]; !ok {
		g.AddVertex(edge.GetFrom())
	}
	if _, ok := g.vertexMap[edge.GetTo().GetID()]; !ok {
		g.AddVertex(edge.GetTo())
	}
	g.dirty = true
}

func (g *MemGraph[EdgePropertyType, VertexPropertyType]) AddVertex(vertex Vertex[VertexPropertyType]) {
	if _, ok := g.vertexMap[vertex.GetID()]; ok {
		slog.Warn("vertex already exists", "Id", vertex.GetID())
	}
	g.vertexMap[vertex.GetID()] = vertex
	g.dirty = true
}

func (g *MemGraph[EdgePropertyType, VertexPropertyType]) DeleteEdge(edge Edge[EdgePropertyType, VertexPropertyType]) {
	if _, ok := g.edgeMap[edge.GetID()]; ok {
		delete(g.edgeMap, edge.GetID())
	}
	slog.Warn("edge does not exist", "Id", edge.GetID())
	g.dirty = true
}

func (g *MemGraph[EdgePropertyType, VertexPropertyType]) DeleteVertex(vertex Vertex[VertexPropertyType]) {
	if _, ok := g.vertexMap[vertex.GetID()]; ok {
		delete(g.vertexMap, vertex.GetID())
	}
	slog.Warn("vertex does not exist", "Id", vertex.GetID())
	g.dirty = true
}

func (g *MemGraph[EdgePropertyType, VertexPropertyType]) GetOutEdges(vertex Vertex[VertexPropertyType]) []Edge[EdgePropertyType, VertexPropertyType] {
	g.updateIndex()
	return g.vertex2OutEdges[vertex.GetID()]
}

func (g *MemGraph[EdgePropertyType, VertexPropertyType]) GetInEdges(vertex Vertex[VertexPropertyType]) []Edge[EdgePropertyType, VertexPropertyType] {
	g.updateIndex()
	return g.vertex2InEdges[vertex.GetID()]
}

func (g *MemGraph[EdgePropertyType, VertexPropertyType]) GetAllVertices() []Vertex[VertexPropertyType] {
	g.updateIndex()
	var all []Vertex[VertexPropertyType]
	for _, v := range g.vertexMap {
		all = append(all, v)
	}
	return all
}

func (g *MemGraph[EdgePropertyType, VertexPropertyType]) GetAllEdges() []Edge[EdgePropertyType, VertexPropertyType] {
	g.updateIndex()
	var all []Edge[EdgePropertyType, VertexPropertyType]
	for _, e := range g.edgeMap {
		all = append(all, e)
	}
	return all
}

func (g *MemGraph[EdgePropertyType, VertexPropertyType]) Save(file string) error {
	panic("unimplemented")
}
func (g *MemGraph[EdgePropertyType, VertexPropertyType]) Load(file string) error {
	panic("unimplemented")
}
