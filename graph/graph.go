package graph

import (
	"fmt"
	"io"
	"os"
	"text/template"
)

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

type Graph[EdgePropertyType any, VertexPropertyType any] interface {
	AddEdge(edge Edge[EdgePropertyType, VertexPropertyType])
	AddVertex(vertex Vertex[VertexPropertyType])
	DeleteEdge(edge Edge[EdgePropertyType, VertexPropertyType])
	DeleteVertex(vertex Vertex[VertexPropertyType])
	GetOutEdges(vertex Vertex[VertexPropertyType]) []Edge[EdgePropertyType, VertexPropertyType]
	GetInEdges(vertex Vertex[VertexPropertyType]) []Edge[EdgePropertyType, VertexPropertyType]
	GetAllVertices() []Vertex[VertexPropertyType]
	GetAllEdges() []Edge[EdgePropertyType, VertexPropertyType]
	SetMetadata(key Metadata, val any)
	GetMetadata(key Metadata) any
	GetAllMetadata() map[Metadata]any
	GetVertexById(id string) Vertex[VertexPropertyType]
	GetEdgeById(id string) Edge[EdgePropertyType, VertexPropertyType]
}

type Edge[EdgePropertyType any, VertexPropertyType any] interface {
	SetID(id string)
	SetFrom(vertex Vertex[VertexPropertyType])
	SetTo(vertex Vertex[VertexPropertyType])
	SetProperty(key string, val EdgePropertyType)
	GetID() string
	GetFrom() Vertex[VertexPropertyType]
	GetTo() Vertex[VertexPropertyType]
	GetProperty(key string) EdgePropertyType
	GetAllProperties() map[string]EdgePropertyType
	SetMeta(m any)
	GetMeta() any
}

type Vertex[VertexPropertyType any] interface {
	SetID(id string)
	SetProperty(key string, val VertexPropertyType)
	GetID() string
	GetProperty(key string) VertexPropertyType
	GetAllProperties() map[string]VertexPropertyType
	SetMeta(m any)
	GetMeta() any
}

func CloneVertex[Vpt any](v Vertex[Vpt], newVertex func() Vertex[Vpt]) Vertex[Vpt] {
	clonedVertex := newVertex() // Use the factory function To create a new vertex instance
	clonedVertex.SetID(v.GetID())
	// Retrieve and set all properties
	for key, val := range v.GetAllProperties() {
		clonedVertex.SetProperty(key, val)
	}
	clonedVertex.SetMeta(v.GetMeta())
	return clonedVertex
}

// Clone Ept: EdgePropertyType, Vpt: VertexPropertyType
func Clone[Ept any, Vpt any](graph Graph[Ept, Vpt], newGraph func(...Option) Graph[Ept, Vpt], newEdge func() Edge[Ept, Vpt], newVertex func() Vertex[Vpt]) Graph[Ept, Vpt] {
	// Use the provided factory function To create a new graph instance
	clonedGraph := newGraph()
	for k, v := range graph.GetAllMetadata() {
		clonedGraph.SetMetadata(k, v)
	}

	// Create a map To track the mapping From original vertices To cloned vertices
	vertexMap := make(map[string]Vertex[Vpt])

	// Clone all vertices
	for _, v := range graph.GetAllVertices() {
		clonedVertex := newVertex() // Use the factory function To create a new vertex instance
		clonedVertex.SetID(v.GetID())
		// Retrieve and set all properties
		for key, val := range v.GetAllProperties() {
			clonedVertex.SetProperty(key, val)
		}
		// Add To the new graph and update the map
		clonedGraph.AddVertex(clonedVertex)
		vertexMap[v.GetID()] = clonedVertex
	}

	// Clone all edges
	for _, e := range graph.GetAllEdges() {
		clonedEdge := newEdge() // Use the factory function To create a new edge instance
		clonedEdge.SetID(e.GetID())
		// Set the start and end points, using the map To find the corresponding cloned vertices
		clonedEdge.SetFrom(vertexMap[e.GetFrom().GetID()])
		clonedEdge.SetTo(vertexMap[e.GetTo().GetID()])
		// Retrieve and set all properties
		for key, val := range e.GetAllProperties() {
			clonedEdge.SetProperty(key, val)
		}
		if clonedEdge.GetFrom() == nil || clonedEdge.GetTo() == nil {
			fmt.Println("nil")
		}
		// Add To the new graph
		clonedGraph.AddEdge(clonedEdge)
	}

	// Return the cloned graph
	return clonedGraph
}

func Visualize[EdgePropertyType any, VertexPropertyType any](graph Graph[EdgePropertyType, VertexPropertyType], filename string, labelFunc func(vertex Vertex[VertexPropertyType]) string, parseID func(vertex Vertex[VertexPropertyType]) string) error {
	desc, err := generateDOT(graph, labelFunc, parseID)
	if err != nil {
		return fmt.Errorf("failed To generate DOT description: %w", err)
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

// design flaw: only vertex property can be shown
func generateDOT[EdgePropertyType any, VertexPropertyType any](g Graph[EdgePropertyType, VertexPropertyType], f func(node Vertex[VertexPropertyType]) string, parseID func(node Vertex[VertexPropertyType]) string) (description[VertexPropertyType], error) {
	desc := description[VertexPropertyType]{
		GraphType:    "graph",
		Attributes:   make(map[string]string),
		EdgeOperator: "--",
		Statements:   make([]statement[VertexPropertyType], 0),
	}
	if f == nil {
		f = func(node Vertex[VertexPropertyType]) string {
			return node.GetID()
		}
	}
	if parseID == nil {
		parseID = func(node Vertex[VertexPropertyType]) string {
			return node.GetID()
		}
	}

	desc.GraphType = "digraph"
	desc.EdgeOperator = "->"

	for _, vertex := range g.GetAllVertices() {
		stmt := statement[VertexPropertyType]{
			Source:       parseID(vertex),
			SourceWeight: 1,
			VertexLabel:  f(vertex),
		}
		desc.Statements = append(desc.Statements, stmt)

		for _, edge := range g.GetOutEdges(vertex) {
			stmt1 := statement[VertexPropertyType]{
				Source:     parseID(vertex),
				Target:     parseID(edge.GetTo()),
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
		return fmt.Errorf("failed To parse template: %w", err)
	}

	return tpl.Execute(w, d)
}

func GenerateCypher[EdgePropertyType any, VertexPropertyType any](graph Graph[EdgePropertyType, VertexPropertyType], f func(n Vertex[VertexPropertyType]) string) ([]string, []string) {
	// Initialize strings To hold Cypher statements
	var nodesCyphers []string
	var edgesCyphers []string

	// Generate Cypher for nodes
	for _, node := range graph.GetAllVertices() {
		// Add a Cypher statement for the current node
		nodesCypher := f(node)
		nodesCyphers = append(nodesCyphers, nodesCypher)
	}

	// Generate Cypher for edges
	for _, edge := range graph.GetAllEdges() {
		// Add a Cypher statement for the current edge
		edgesCypher := ""
		edgesCypher += fmt.Sprintf("MATCH (a {ID: '%s'}), (b {ID: '%s'}) CREATE (a)-[:child]->(b)", edge.GetFrom().GetID(), edge.GetTo().GetID())
		edgesCyphers = append(edgesCyphers, edgesCypher)
	}
	return nodesCyphers, edgesCyphers
}

type Options struct {
	FSEnabled    bool
	ReadFileName string
}
type Option func(*Options)

func WithPersistent(persistent bool) Option {
	return func(o *Options) {
		o.FSEnabled = persistent
	}
}

func NewGraph[EdgePropertyType any, VertexPropertyType any](opts ...Option) Graph[EdgePropertyType, VertexPropertyType] {
	options := Options{}
	for _, opt := range opts {
		opt(&options)
	}

	if options.FSEnabled {
		m := &FSGraph[EdgePropertyType, VertexPropertyType]{
			EdgeMap:         make(map[string]Edge[EdgePropertyType, VertexPropertyType]),
			VertexMap:       make(map[string]Vertex[VertexPropertyType]),
			Vertex2OutEdges: make(map[string][]Edge[EdgePropertyType, VertexPropertyType]),
			Vertex2InEdges:  make(map[string][]Edge[EdgePropertyType, VertexPropertyType]),
			Dirty:           false,
			Metadata:        make(map[Metadata]any),
		}
		m.SetMetadata(CleanVertexByEdge, false)
		return m
	}

	if options.ReadFileName != "" {
		m, err := load[EdgePropertyType, VertexPropertyType](options.ReadFileName)
		if err != nil {
			panic(err)
		}
		return m
	}

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
	return &FSVertexImpl[PropertyType]{
		Id:          "",
		PropertyMap: make(map[string]PropertyType),
	}
}

func NewEdge[EdgePropertyType any, VertexPropertyType any]() Edge[EdgePropertyType, VertexPropertyType] {
	return &FSEdgeImpl[EdgePropertyType, VertexPropertyType]{
		Id:          "",
		From:        nil,
		To:          nil,
		PropertyMap: make(map[string]EdgePropertyType),
	}
}
