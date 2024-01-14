package graph

import (
	"bytes"
	"encoding/gob"
	"log/slog"
	"os"
)

type FSGraph[EdgePropertyType any, VertexPropertyType any] struct {
	EdgeMap   map[string]Edge[EdgePropertyType, VertexPropertyType]
	VertexMap map[string]Vertex[VertexPropertyType]
	Metadata  map[Metadata]any

	// index
	Vertex2OutEdges map[string][]Edge[EdgePropertyType, VertexPropertyType]
	Vertex2InEdges  map[string][]Edge[EdgePropertyType, VertexPropertyType]
	Dirty           bool
}

type FSVertexImpl[PropertyType any] struct {
	Id          string
	PropertyMap map[string]PropertyType
	Meta        any
}

type FSEdgeImpl[EdgePropertyType any, VertexPropertyType any] struct {
	Id          string
	From        Vertex[VertexPropertyType]
	To          Vertex[VertexPropertyType]
	PropertyMap map[string]EdgePropertyType
	Meta        any
}

func (g *FSGraph[EdgePropertyType, VertexPropertyType]) updateIndex() {
	in := func(m map[string][]Edge[EdgePropertyType, VertexPropertyType], key string) bool {
		_, ok := m[key]
		return ok
	}
	if g.Dirty {
		if g.Metadata[CleanVertexByEdge].(bool) && len(g.EdgeMap) == 0 {
			g.VertexMap = make(map[string]Vertex[VertexPropertyType])
		}
		g.Vertex2OutEdges = make(map[string][]Edge[EdgePropertyType, VertexPropertyType])
		g.Vertex2InEdges = make(map[string][]Edge[EdgePropertyType, VertexPropertyType])
		for _, e := range g.EdgeMap {
			if _, ok := g.Vertex2InEdges[e.GetTo().GetID()]; !ok {
				g.Vertex2InEdges[e.GetTo().GetID()] = make([]Edge[EdgePropertyType, VertexPropertyType], 0)
			}
			if _, ok := g.Vertex2OutEdges[e.GetFrom().GetID()]; !ok {
				g.Vertex2OutEdges[e.GetFrom().GetID()] = make([]Edge[EdgePropertyType, VertexPropertyType], 0)
			}
			g.Vertex2InEdges[e.GetTo().GetID()] = append(g.Vertex2InEdges[e.GetTo().GetID()], e)
			g.Vertex2OutEdges[e.GetFrom().GetID()] = append(g.Vertex2OutEdges[e.GetFrom().GetID()], e)

			if g.Metadata[CleanVertexByEdge].(bool) {
				// clean single vertex
				if !in(g.Vertex2InEdges, e.GetFrom().GetID()) && !in(g.Vertex2OutEdges, e.GetFrom().GetID()) {
					delete(g.VertexMap, e.GetFrom().GetID())
				}
				if !in(g.Vertex2InEdges, e.GetTo().GetID()) && !in(g.Vertex2OutEdges, e.GetTo().GetID()) {
					delete(g.VertexMap, e.GetTo().GetID())
				}
			}
		}

		g.Dirty = false
	}
}

func (g *FSGraph[EdgePropertyType, VertexPropertyType]) SetMetadata(key Metadata, val any) {
	g.Metadata[key] = val
}
func (g *FSGraph[EdgePropertyType, VertexPropertyType]) GetMetadata(key Metadata) any {
	return g.Metadata[key]
}

func (g *FSGraph[EdgePropertyType, VertexPropertyType]) GetAllMetadata() map[Metadata]any {
	return g.Metadata
}

func (g *FSGraph[EdgePropertyType, VertexPropertyType]) GetVertexById(id string) Vertex[VertexPropertyType] {
	return g.VertexMap[id]
}
func (g *FSGraph[EdgePropertyType, VertexPropertyType]) GetEdgeById(id string) Edge[EdgePropertyType, VertexPropertyType] {
	return g.EdgeMap[id]
}
func (g *FSGraph[EdgePropertyType, VertexPropertyType]) AddEdge(edge Edge[EdgePropertyType, VertexPropertyType]) {
	if _, ok := g.EdgeMap[edge.GetID()]; ok {
		slog.Debug("edge already exists", "Id", edge.GetID(), "From", edge.GetFrom().GetID(), "To", edge.GetTo().GetID())
	}
	g.EdgeMap[edge.GetID()] = edge
	if _, ok := g.VertexMap[edge.GetFrom().GetID()]; !ok {
		g.AddVertex(edge.GetFrom())
	}
	if _, ok := g.VertexMap[edge.GetTo().GetID()]; !ok {
		g.AddVertex(edge.GetTo())
	}
	g.Dirty = true
}

func (g *FSGraph[EdgePropertyType, VertexPropertyType]) AddVertex(vertex Vertex[VertexPropertyType]) {
	if _, ok := g.VertexMap[vertex.GetID()]; ok {
		slog.Warn("vertex already exists", "Id", vertex.GetID())
	}
	g.VertexMap[vertex.GetID()] = vertex
	g.Dirty = true
}

func (g *FSGraph[EdgePropertyType, VertexPropertyType]) DeleteEdge(edge Edge[EdgePropertyType, VertexPropertyType]) {
	if _, ok := g.EdgeMap[edge.GetID()]; ok {
		delete(g.EdgeMap, edge.GetID())
	}
	slog.Warn("edge does not exist", "Id", edge.GetID())
	g.Dirty = true
}

func (g *FSGraph[EdgePropertyType, VertexPropertyType]) DeleteVertex(vertex Vertex[VertexPropertyType]) {
	if _, ok := g.VertexMap[vertex.GetID()]; ok {
		delete(g.VertexMap, vertex.GetID())
	}
	slog.Warn("vertex does not exist", "Id", vertex.GetID())
	g.Dirty = true
}

func (g *FSGraph[EdgePropertyType, VertexPropertyType]) GetOutEdges(vertex Vertex[VertexPropertyType]) []Edge[EdgePropertyType, VertexPropertyType] {
	g.updateIndex()
	return g.Vertex2OutEdges[vertex.GetID()]
}

func (g *FSGraph[EdgePropertyType, VertexPropertyType]) GetInEdges(vertex Vertex[VertexPropertyType]) []Edge[EdgePropertyType, VertexPropertyType] {
	g.updateIndex()
	return g.Vertex2InEdges[vertex.GetID()]
}

func (g *FSGraph[EdgePropertyType, VertexPropertyType]) GetAllVertices() []Vertex[VertexPropertyType] {
	g.updateIndex()
	var all []Vertex[VertexPropertyType]
	for _, v := range g.VertexMap {
		all = append(all, v)
	}
	return all
}

func (g *FSGraph[EdgePropertyType, VertexPropertyType]) GetAllEdges() []Edge[EdgePropertyType, VertexPropertyType] {
	g.updateIndex()
	var all []Edge[EdgePropertyType, VertexPropertyType]
	for _, e := range g.EdgeMap {
		all = append(all, e)
	}
	return all
}
func (g *FSGraph[EdgePropertyType, VertexPropertyType]) Save(file string) error {
	graph, err := serializeGraph(g)
	if err != nil {
		return err
	}
	err = os.WriteFile(file, graph, 0644)
	return err
}

func load[EdgePropertyType any, VertexPropertyType any](file string) (Graph[EdgePropertyType, VertexPropertyType], error) {
	fileData, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	deserializedGraph, err := deserializeGraph[EdgePropertyType, VertexPropertyType](fileData)
	if err != nil {
		return nil, err
	}
	return deserializedGraph, nil
}
func (n *FSVertexImpl[PropertyType]) SetID(id string) {
	n.Id = id
}
func (n *FSVertexImpl[PropertyType]) SetMeta(m any) {
	n.Meta = m
}
func (n *FSVertexImpl[PropertyType]) GetMeta() any {
	return n.Meta
}
func (n *FSVertexImpl[PropertyType]) SetProperty(key string, val PropertyType) {
	n.PropertyMap[key] = val
}

func (n *FSVertexImpl[PropertyType]) GetID() string {
	return n.Id
}

func (n *FSVertexImpl[PropertyType]) GetProperty(key string) PropertyType {
	val, _ := n.PropertyMap[key]
	return val
}

func (n *FSVertexImpl[PropertyType]) GetAllProperties() map[string]PropertyType {
	return n.PropertyMap
}

func (e *FSEdgeImpl[EdgePropertyType, VertexPropertyType]) SetID(id string) {
	e.Id = id
}

func (e *FSEdgeImpl[EdgePropertyType, VertexPropertyType]) SetFrom(vertex Vertex[VertexPropertyType]) {
	e.From = vertex
}

func (e *FSEdgeImpl[EdgePropertyType, VertexPropertyType]) SetTo(vertex Vertex[VertexPropertyType]) {
	e.To = vertex
}

func (e *FSEdgeImpl[EdgePropertyType, VertexPropertyType]) SetProperty(key string, val EdgePropertyType) {
	e.PropertyMap[key] = val
}

func (e *FSEdgeImpl[EdgePropertyType, VertexPropertyType]) GetID() string {
	return e.Id
}

func (e *FSEdgeImpl[EdgePropertyType, VertexPropertyType]) GetFrom() Vertex[VertexPropertyType] {
	return e.From
}

func (e *FSEdgeImpl[EdgePropertyType, VertexPropertyType]) GetTo() Vertex[VertexPropertyType] {
	return e.To
}

func (e *FSEdgeImpl[EdgePropertyType, VertexPropertyType]) GetProperty(key string) EdgePropertyType {
	val, _ := e.PropertyMap[key]
	return val
}

func (e *FSEdgeImpl[EdgePropertyType, VertexPropertyType]) GetAllProperties() map[string]EdgePropertyType {
	return e.PropertyMap
}
func (e *FSEdgeImpl[EdgePropertyType, VertexPropertyType]) SetMeta(m any) {
	e.Meta = m
}
func (e *FSEdgeImpl[EdgePropertyType, VertexPropertyType]) GetMeta() any {
	return e.Meta
}

func serializeGraph[EdgePropertyType any, VertexPropertyType any](graph *FSGraph[EdgePropertyType, VertexPropertyType]) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	// Encoding the graph
	if err := encoder.Encode(graph); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
func deserializeGraph[EdgePropertyType any, VertexPropertyType any](data []byte) (*FSGraph[EdgePropertyType, VertexPropertyType], error) {
	var graph FSGraph[EdgePropertyType, VertexPropertyType]
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)

	// Decoding the graph
	if err := decoder.Decode(&graph); err != nil {
		return nil, err
	}

	return &graph, nil
}
