package constraint

var NodeHooks map[string][]NodeCons

func init() {
	NodeHooks = make(map[string][]NodeCons)
}

// Constraint definition of constraint
type Constraint[Tv1 any, Tv2 any] interface {
	Eval() bool
	GetVar1() Var[Tv1]
	GetVar2() Var[Tv2]
}

type NodeCons interface {
	Eval() bool
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

func RegisterConstraint[Tv1 any, Tv2 any](id string, c Constraint[Tv1, Tv2]) {
	if _, ok := NodeHooks[id]; !ok {
		NodeHooks[id] = make([]NodeCons, 0)
	}
	NodeHooks[id] = append(NodeHooks[id], c)
}
