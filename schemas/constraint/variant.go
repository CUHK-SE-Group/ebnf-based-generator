package constraint

import "github.com/CUHK-SE-Group/generic-generator/schemas"

// Var definition of variable of the metadata
type Var[T any] interface {
	Get() T
	Type() string
}
type vvar[T any] struct {
	ttype   string
	content T
}

func (v *vvar[T]) Get() T {
	return v.content
}

func (v *vvar[T]) Type() string {
	return v.ttype
}
func NewVar[T any](t string, content T) Var[T] {
	return &vvar[T]{content: content, ttype: t}
}

/*
List of Variables
*/

type VarFactory[T any] interface {
	Build() func() Var[T]
}

// NodeInfo node information
type NodeInfo struct {
	ID               string
	OccurredTimes    int
	OccurredPosition string
}

// BuildFromNode Example:
// node1:=&Node{}
// node2:=&Node{}
// BuildConstraint[*Node, *Node](DefinedBefore, BuildFromNode{n:node1}.Build()(), BuildFromNode{n:node2}.Build()())
type BuildFromNode struct {
	n *schemas.Node
}

func (b *BuildFromNode) Build() func() Var[NodeInfo] {
	return func() Var[NodeInfo] {
		return NewVar[NodeInfo]("node", NodeInfo{
			ID:               "test",
			OccurredTimes:    1,
			OccurredPosition: "34.5.2.4",
		})
	}
}
