package graph

import (
	"log/slog"
)

const (
	CleanVertexByEdge Metadata = "cleanVertexWhenNoEdge"
)

type MemGraph struct {
	edgeMap   map[string]Edge[string]
	vertexMap map[string]Vertex[string]
	metadata  map[Metadata]any

	// index
	vertex2OutEdges map[string][]Edge[string]
	vertex2InEdges  map[string][]Edge[string]
	dirty           bool
}

type VertexImpl struct {
	id          string
	propertyMap map[string]string
}

type EdgeImpl struct {
	id          string
	from        Vertex[string]
	to          Vertex[string]
	propertyMap map[string]string
}

func NewGraph() Graph[string] {
	return &MemGraph{
		edgeMap:         make(map[string]Edge[string]),
		vertexMap:       make(map[string]Vertex[string]),
		vertex2OutEdges: make(map[string][]Edge[string]),
		vertex2InEdges:  make(map[string][]Edge[string]),
		dirty:           false,
		metadata:        make(map[Metadata]any),
	}
}

func NewVertex() Vertex[string] {
	return &VertexImpl{
		id:          "",
		propertyMap: make(map[string]string),
	}
}

func NewEdge() Edge[string] {
	return &EdgeImpl{
		id:          "",
		from:        nil,
		to:          nil,
		propertyMap: make(map[string]string),
	}
}

func (g *MemGraph) updateIndex() {
	in := func(m map[string][]Edge[string], key string) bool {
		_, ok := m[key]
		return ok
	}
	if g.dirty {
		if g.metadata[CleanVertexByEdge].(bool) && len(g.edgeMap) == 0 {
			g.vertexMap = make(map[string]Vertex[string])
		}
		g.vertex2OutEdges = make(map[string][]Edge[string])
		g.vertex2InEdges = make(map[string][]Edge[string])
		for _, e := range g.edgeMap {
			if _, ok := g.vertex2InEdges[e.GetTo().GetID()]; !ok {
				g.vertex2InEdges[e.GetTo().GetID()] = make([]Edge[string], 0)
			}
			if _, ok := g.vertex2OutEdges[e.GetFrom().GetID()]; !ok {
				g.vertex2OutEdges[e.GetFrom().GetID()] = make([]Edge[string], 0)
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

func (g *MemGraph) SetMetadata(key Metadata, val any) {
	g.metadata[key] = val
}
func (g *MemGraph) GetMetadata(key Metadata) any {
	return g.metadata[key]
}

func (g *MemGraph) GetAllMetadata() map[Metadata]any {
	return g.metadata
}

func (g *MemGraph) GetVertexById(id string) Vertex[string] {
	return g.vertexMap[id]
}
func (g *MemGraph) GetEdgeById(id string) Edge[string] {
	return g.edgeMap[id]
}
func (g *MemGraph) AddEdge(edge Edge[string]) {
	if _, ok := g.edgeMap[edge.GetID()]; ok {
		slog.Warn("edge already exists", "id", edge.GetID())
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

func (g *MemGraph) AddVertex(vertex Vertex[string]) {
	if _, ok := g.vertexMap[vertex.GetID()]; ok {
		slog.Warn("vertex already exists", "id", vertex.GetID())
	}
	g.vertexMap[vertex.GetID()] = vertex
	g.dirty = true
}

func (g *MemGraph) DeleteEdge(edge Edge[string]) {
	if _, ok := g.edgeMap[edge.GetID()]; ok {
		delete(g.edgeMap, edge.GetID())
	}
	slog.Warn("edge does not exist", "id", edge.GetID())
	g.dirty = true
}

func (g *MemGraph) DeleteVertex(vertex Vertex[string]) {
	if _, ok := g.vertexMap[vertex.GetID()]; ok {
		delete(g.vertexMap, vertex.GetID())
	}
	slog.Warn("vertex does not exist", "id", vertex.GetID())
	g.dirty = true
}

func (g *MemGraph) GetOutEdges(vertex Vertex[string]) []Edge[string] {
	g.updateIndex()
	return g.vertex2OutEdges[vertex.GetID()]
}

func (g *MemGraph) GetInEdges(vertex Vertex[string]) []Edge[string] {
	g.updateIndex()
	return g.vertex2InEdges[vertex.GetID()]
}

func (g *MemGraph) GetAllVertices() []Vertex[string] {
	g.updateIndex()
	var all []Vertex[string]
	for _, v := range g.vertexMap {
		all = append(all, v)
	}
	return all
}

func (g *MemGraph) GetAllEdges() []Edge[string] {
	g.updateIndex()
	var all []Edge[string]
	for _, e := range g.edgeMap {
		all = append(all, e)
	}
	return all
}

func (n *VertexImpl) SetID(id string) {
	n.id = id
}

func (n *VertexImpl) SetProperty(key string, val string) {
	n.propertyMap[key] = val
}

func (n *VertexImpl) GetID() string {
	return n.id
}

func (n *VertexImpl) GetProperty(key string) string {
	if _, ok := n.propertyMap[key]; ok {
		return n.propertyMap[key]
	}
	return ""
}

func (n *VertexImpl) GetAllProperties() map[string]string {
	return n.propertyMap
}

func (e *EdgeImpl) SetID(id string) {
	e.id = id
}

func (e *EdgeImpl) SetFrom(vertex Vertex[string]) {
	e.from = vertex
}

func (e *EdgeImpl) SetTo(vertex Vertex[string]) {
	e.to = vertex
}

func (e *EdgeImpl) SetProperty(key string, val string) {
	e.propertyMap[key] = val
}

func (e *EdgeImpl) GetID() string {
	return e.id
}

func (e *EdgeImpl) GetFrom() Vertex[string] {
	return e.from
}

func (e *EdgeImpl) GetTo() Vertex[string] {
	return e.to
}

func (e *EdgeImpl) GetProperty(key string) string {
	if _, ok := e.propertyMap[key]; ok {
		return e.propertyMap[key]
	}
	return ""
}

func (e *EdgeImpl) GetAllProperties() map[string]string {
	return e.propertyMap
}
