package schemas

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

// Constraint definition of constraint
type Constraint[Tv1 any, Tv2 any] interface {
	Eval() bool
	GetVar1() Var[Tv1]
	GetVar2() Var[Tv2]
}

type constraint[Tv1 any, Tv2 any] struct {
	op Operator[Tv1, Tv2]
	v1 Var[Tv1]
	v2 Var[Tv2]
}

func (c *constraint[Tv1, Tv2]) Eval() bool {
	return c.op(c.v1, c.v2)
}

func (c *constraint[Tv1, Tv2]) GetVar1() Var[Tv1] {
	return c.v1
}
func (c *constraint[Tv1, Tv2]) GetVar2() Var[Tv2] {
	return c.v2
}

func BuildConstraint[Tv1 any, Tv2 any](op Operator[Tv1, Tv2], v1 Var[Tv1], v2 Var[Tv2]) Constraint[Tv1, Tv2] {
	return &constraint[Tv1, Tv2]{
		op: op,
		v1: v1,
		v2: v2,
	}
}

/*
List of Variables
*/

type VarFactory[T any] interface {
	Build() func() Var[T]
}

// BuildFromNode Example:
// node1:=&Node{}
// node2:=&Node{}
// BuildConstraint[*Node, *Node](DefinedBefore, BuildFromNode{n:node1}.Build()(), BuildFromNode{n:node2}.Build()())
type BuildFromNode struct {
	n *Node
}

func (b *BuildFromNode) Build() func() Var[*Node] {
	return func() Var[*Node] {
		return NewVar[*Node]("node", b.n)
	}
}

/*
list of operators
*/

type Operator[Tv1 any, Tv2 any] func(v1 Var[Tv1], v2 Var[Tv2]) bool

func DefinedBefore(v1 Var[string], v2 Var[string]) bool {
	return true
}
