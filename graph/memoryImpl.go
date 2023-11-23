package graph

import (
	"log/slog"
)

const (
	CleanVertexByEdge Metadata = "cleanVertexWhenNoEdge"
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

type VertexImpl[PropertyType any] struct {
	id          string
	propertyMap map[string]PropertyType
	meta        any
}

type EdgeImpl[EdgePropertyType any, VertexPropertyType any] struct {
	id          string
	from        Vertex[VertexPropertyType]
	to          Vertex[VertexPropertyType]
	propertyMap map[string]EdgePropertyType
	meta        any
}

func NewGraph[EdgePropertyType any, VertexPropertyType any]() Graph[EdgePropertyType, VertexPropertyType] {
	m := &MemGraph[EdgePropertyType, VertexPropertyType]{
		edgeMap:         make(map[string]Edge[EdgePropertyType, VertexPropertyType]),
		vertexMap:       make(map[string]Vertex[VertexPropertyType]),
		vertex2OutEdges: make(map[string][]Edge[EdgePropertyType, VertexPropertyType]),
		vertex2InEdges:  make(map[string][]Edge[EdgePropertyType, VertexPropertyType]),
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

func NewEdge[EdgePropertyType any, VertexPropertyType any]() Edge[EdgePropertyType, VertexPropertyType] {
	return &EdgeImpl[EdgePropertyType, VertexPropertyType]{
		id:          "",
		from:        nil,
		to:          nil,
		propertyMap: make(map[string]EdgePropertyType),
	}
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
		slog.Warn("edge already exists", "id", edge.GetID(), "from", edge.GetFrom().GetID(), "to", edge.GetTo().GetID())
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
		slog.Warn("vertex already exists", "id", vertex.GetID())
	}
	g.vertexMap[vertex.GetID()] = vertex
	g.dirty = true
}

func (g *MemGraph[EdgePropertyType, VertexPropertyType]) DeleteEdge(edge Edge[EdgePropertyType, VertexPropertyType]) {
	if _, ok := g.edgeMap[edge.GetID()]; ok {
		delete(g.edgeMap, edge.GetID())
	}
	slog.Warn("edge does not exist", "id", edge.GetID())
	g.dirty = true
}

func (g *MemGraph[EdgePropertyType, VertexPropertyType]) DeleteVertex(vertex Vertex[VertexPropertyType]) {
	if _, ok := g.vertexMap[vertex.GetID()]; ok {
		delete(g.vertexMap, vertex.GetID())
	}
	slog.Warn("vertex does not exist", "id", vertex.GetID())
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

func (n *VertexImpl[PropertyType]) SetID(id string) {
	n.id = id
}
func (n *VertexImpl[PropertyType]) SetMeta(m any) {
	n.meta = m
}
func (n *VertexImpl[PropertyType]) GetMeta() any {
	return n.meta
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

func (e *EdgeImpl[EdgePropertyType, VertexPropertyType]) SetID(id string) {
	e.id = id
}

func (e *EdgeImpl[EdgePropertyType, VertexPropertyType]) SetFrom(vertex Vertex[VertexPropertyType]) {
	e.from = vertex
}

func (e *EdgeImpl[EdgePropertyType, VertexPropertyType]) SetTo(vertex Vertex[VertexPropertyType]) {
	e.to = vertex
}

func (e *EdgeImpl[EdgePropertyType, VertexPropertyType]) SetProperty(key string, val EdgePropertyType) {
	e.propertyMap[key] = val
}

func (e *EdgeImpl[EdgePropertyType, VertexPropertyType]) GetID() string {
	return e.id
}

func (e *EdgeImpl[EdgePropertyType, VertexPropertyType]) GetFrom() Vertex[VertexPropertyType] {
	return e.from
}

func (e *EdgeImpl[EdgePropertyType, VertexPropertyType]) GetTo() Vertex[VertexPropertyType] {
	return e.to
}

func (e *EdgeImpl[EdgePropertyType, VertexPropertyType]) GetProperty(key string) EdgePropertyType {
	val, _ := e.propertyMap[key]
	return val
}

func (e *EdgeImpl[EdgePropertyType, VertexPropertyType]) GetAllProperties() map[string]EdgePropertyType {
	return e.propertyMap
}
func (e *EdgeImpl[EdgePropertyType, VertexPropertyType]) SetMeta(m any) {
	e.meta = m
}
func (e *EdgeImpl[EdgePropertyType, VertexPropertyType]) GetMeta() any {
	return e.meta
}
