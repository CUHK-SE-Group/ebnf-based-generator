package graph

import (
	"log/slog"
)

const (
	CleanVertexByEdge Metadata = "cleanVertexWhenNoEdge"
)

type MemGraph[PropertyType any] struct {
	edgeMap   map[string]Edge[PropertyType]
	vertexMap map[string]Vertex[PropertyType]
	metadata  map[Metadata]any

	// index
	vertex2OutEdges map[string][]Edge[PropertyType]
	vertex2InEdges  map[string][]Edge[PropertyType]
	dirty           bool
}

type VertexImpl[PropertyType any] struct {
	id          string
	propertyMap map[string]PropertyType
}

type EdgeImpl[PropertyType any] struct {
	id          string
	from        Vertex[PropertyType]
	to          Vertex[PropertyType]
	propertyMap map[string]PropertyType
}

func NewGraph[PropertyType any]() Graph[PropertyType] {
	m := &MemGraph[PropertyType]{
		edgeMap:         make(map[string]Edge[PropertyType]),
		vertexMap:       make(map[string]Vertex[PropertyType]),
		vertex2OutEdges: make(map[string][]Edge[PropertyType]),
		vertex2InEdges:  make(map[string][]Edge[PropertyType]),
		dirty:           false,
		metadata:        make(map[Metadata]any),
	}
	m.SetMetadata(CleanVertexByEdge, false)
	return m
}

func NewVertex[PropertyType any]() Vertex[PropertyType] {
	return &VertexImpl[PropertyType]{
		id:          "",
		propertyMap: make(map[string]PropertyType),
	}
}

func NewEdge[PropertyType any]() Edge[PropertyType] {
	return &EdgeImpl[PropertyType]{
		id:          "",
		from:        nil,
		to:          nil,
		propertyMap: make(map[string]PropertyType),
	}
}

func (g *MemGraph[PropertyType]) updateIndex() {
	in := func(m map[string][]Edge[PropertyType], key string) bool {
		_, ok := m[key]
		return ok
	}
	if g.dirty {
		if g.metadata[CleanVertexByEdge].(bool) && len(g.edgeMap) == 0 {
			g.vertexMap = make(map[string]Vertex[PropertyType])
		}
		g.vertex2OutEdges = make(map[string][]Edge[PropertyType])
		g.vertex2InEdges = make(map[string][]Edge[PropertyType])
		for _, e := range g.edgeMap {
			if _, ok := g.vertex2InEdges[e.GetTo().GetID()]; !ok {
				g.vertex2InEdges[e.GetTo().GetID()] = make([]Edge[PropertyType], 0)
			}
			if _, ok := g.vertex2OutEdges[e.GetFrom().GetID()]; !ok {
				g.vertex2OutEdges[e.GetFrom().GetID()] = make([]Edge[PropertyType], 0)
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

func (g *MemGraph[PropertyType]) SetMetadata(key Metadata, val any) {
	g.metadata[key] = val
}
func (g *MemGraph[PropertyType]) GetMetadata(key Metadata) any {
	return g.metadata[key]
}

func (g *MemGraph[PropertyType]) GetAllMetadata() map[Metadata]any {
	return g.metadata
}

func (g *MemGraph[PropertyType]) GetVertexById(id string) Vertex[PropertyType] {
	return g.vertexMap[id]
}
func (g *MemGraph[PropertyType]) GetEdgeById(id string) Edge[PropertyType] {
	return g.edgeMap[id]
}
func (g *MemGraph[PropertyType]) AddEdge(edge Edge[PropertyType]) {
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

func (g *MemGraph[PropertyType]) AddVertex(vertex Vertex[PropertyType]) {
	if _, ok := g.vertexMap[vertex.GetID()]; ok {
		slog.Warn("vertex already exists", "id", vertex.GetID())
	}
	g.vertexMap[vertex.GetID()] = vertex
	g.dirty = true
}

func (g *MemGraph[PropertyType]) DeleteEdge(edge Edge[PropertyType]) {
	if _, ok := g.edgeMap[edge.GetID()]; ok {
		delete(g.edgeMap, edge.GetID())
	}
	slog.Warn("edge does not exist", "id", edge.GetID())
	g.dirty = true
}

func (g *MemGraph[PropertyType]) DeleteVertex(vertex Vertex[PropertyType]) {
	if _, ok := g.vertexMap[vertex.GetID()]; ok {
		delete(g.vertexMap, vertex.GetID())
	}
	slog.Warn("vertex does not exist", "id", vertex.GetID())
	g.dirty = true
}

func (g *MemGraph[PropertyType]) GetOutEdges(vertex Vertex[PropertyType]) []Edge[PropertyType] {
	g.updateIndex()
	return g.vertex2OutEdges[vertex.GetID()]
}

func (g *MemGraph[PropertyType]) GetInEdges(vertex Vertex[PropertyType]) []Edge[PropertyType] {
	g.updateIndex()
	return g.vertex2InEdges[vertex.GetID()]
}

func (g *MemGraph[PropertyType]) GetAllVertices() []Vertex[PropertyType] {
	g.updateIndex()
	var all []Vertex[PropertyType]
	for _, v := range g.vertexMap {
		all = append(all, v)
	}
	return all
}

func (g *MemGraph[PropertyType]) GetAllEdges() []Edge[PropertyType] {
	g.updateIndex()
	var all []Edge[PropertyType]
	for _, e := range g.edgeMap {
		all = append(all, e)
	}
	return all
}

func (n *VertexImpl[PropertyType]) SetID(id string) {
	n.id = id
}

func (n *VertexImpl[PropertyType]) SetProperty(key string, val PropertyType) {
	n.propertyMap[key] = val
}

func (n *VertexImpl[PropertyType]) GetID() string {
	return n.id
}

func (n *VertexImpl[PropertyType]) GetProperty(key string) PropertyType {
	val, _ := n.propertyMap[key]
	return val
}

func (n *VertexImpl[PropertyType]) GetAllProperties() map[string]PropertyType {
	return n.propertyMap
}

func (e *EdgeImpl[PropertyType]) SetID(id string) {
	e.id = id
}

func (e *EdgeImpl[PropertyType]) SetFrom(vertex Vertex[PropertyType]) {
	e.from = vertex
}

func (e *EdgeImpl[PropertyType]) SetTo(vertex Vertex[PropertyType]) {
	e.to = vertex
}

func (e *EdgeImpl[PropertyType]) SetProperty(key string, val PropertyType) {
	e.propertyMap[key] = val
}

func (e *EdgeImpl[PropertyType]) GetID() string {
	return e.id
}

func (e *EdgeImpl[PropertyType]) GetFrom() Vertex[PropertyType] {
	return e.from
}

func (e *EdgeImpl[PropertyType]) GetTo() Vertex[PropertyType] {
	return e.to
}

func (e *EdgeImpl[PropertyType]) GetProperty(key string) PropertyType {
	val, _ := e.propertyMap[key]
	return val
}

func (e *EdgeImpl[PropertyType]) GetAllProperties() map[string]PropertyType {
	return e.propertyMap
}
