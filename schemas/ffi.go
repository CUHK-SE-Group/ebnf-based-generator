package schemas

import (
	"github.com/CUHK-SE-Group/generic-generator/graph"
	"github.com/CUHK-SE-Group/generic-generator/schemas/ffi"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

// This file is hardcoded with the implementation of Grammar
//go:generate protoc --go_out=ffi --go_opt=paths=source_relative ffi.proto

func marshalGrammar(g *Grammar) ([]byte, error) {
	edgeMap := make(map[string]*ffi.FSEdge)
	vertexMap := make(map[string]*ffi.FSVertex)
	metadata := make(map[string]*anypb.Any)
	for _, v := range g.internal.GetAllVertices() {
		vv := 0
		if v.GetMeta() != nil {
			vv = v.GetMeta().(int)
		}
		meta, _ := anypb.New(&ffi.IntValue{Value: uint64(vv)})
		prop := v.GetProperty(Prop)
		vertexMap[v.GetID()] = &ffi.FSVertex{
			Id: v.GetID(),
			PropertyMap: map[string]*ffi.Property{
				Prop: {
					Type:               uint64(prop.Type),
					Content:            prop.Content,
					DistanceToTerminal: int32(prop.DistanceToTerminal),
				},
			},
			Meta: meta,
		}
	}
	for _, e := range g.internal.GetAllEdges() {
		tmp := e.GetMeta().(int)
		meta, _ := anypb.New(&ffi.IntValue{Value: uint64(tmp)})
		edgeMap[e.GetID()] = &ffi.FSEdge{
			Id:   e.GetID(),
			From: e.GetFrom().GetID(),
			To:   e.GetTo().GetID(),
			Meta: meta,
		}
	}
	for k, v := range g.internal.GetAllMetadata() {
		switch v.(type) {
		case bool:
			metadata[string(k)], _ = anypb.New(&ffi.BoolValue{Value: v.(bool)})
		case string:
			metadata[string(k)], _ = anypb.New(&ffi.StringValue{Value: v.(string)})
		}
	}
	graphData := ffi.FSGraph{
		EdgeMap:   edgeMap,
		VertexMap: vertexMap,
		Metadata:  metadata,
		Dirty:     false,
	}

	return proto.Marshal(&graphData)
}
func unmarshalGrammar(data []byte) (*Grammar, error) {
	graphData := &ffi.FSGraph{}
	err := proto.Unmarshal(data, graphData)
	if err != nil {
		return nil, err
	}
	g := graph.NewGraph[string, Property](graph.WithPersistent(true))
	grammar := NewGrammar()
	vertexMap := make(map[string]graph.Vertex[Property])
	for _, v := range graphData.VertexMap {
		n := graph.NewVertex[Property]()
		n.SetID(v.Id)
		n.SetProperty(Prop, Property{
			Type:               GrammarType(v.PropertyMap[Prop].Type),
			Gram:               grammar,
			Content:            v.PropertyMap[Prop].Content,
			DistanceToTerminal: int(v.PropertyMap[Prop].DistanceToTerminal),
		})
		meta := &ffi.IntValue{}
		_ = v.Meta.UnmarshalTo(meta)
		n.SetMeta(int(meta.Value))
		g.AddVertex(n)
		vertexMap[v.Id] = n
	}
	for _, e := range graphData.EdgeMap {
		meta := &ffi.IntValue{}
		_ = e.Meta.UnmarshalTo(meta)
		n := graph.NewEdge[string, Property]()
		n.SetID(e.Id)
		n.SetFrom(vertexMap[e.From])
		n.SetTo(vertexMap[e.To])
		n.SetMeta(int(meta.Value))
		g.AddEdge(n)
	}
	for k, v := range graphData.Metadata {
		var b ffi.BoolValue
		var s ffi.StringValue
		if err := v.UnmarshalTo(&b); err == nil {
			g.SetMetadata(graph.Metadata(k), b.Value)
		}
		if err := v.UnmarshalTo(&s); err == nil {
			g.SetMetadata(graph.Metadata(k), s.Value)
		}
	}
	grammar.internal = g
	return grammar, nil
}
