package schemas

import (
	"fmt"
	"regexp"
	"sort"

	"github.com/CUHK-SE-Group/generic-generator/graph"
)

type ConsAction int
type ConstraintType int

const (
	ConsProp            = "cons"
	FUNC     ConsAction = iota + 1
	REGEX    ConsAction = iota + 1
)

/*
How to define a constraint?
1. There will always be a pre-condition of a constraint.
    if a symbol(or any other information is found in FirstPlaceNode, it needs to be recorded in schemas.Context)
2. Then, in the next step, if a symbol found that it is in SecondPlaceNode, it will read the information that stored in schemas.Context, then take corresponding action
*/

type Action struct {
	Type   ConsAction
	Regexp *regexp.Regexp
	Func   func(ctx *Context) (*Context, error)
}

type Constraint struct {
	ConstraintType ConstraintType
	FirstNode      string
	SecondNode     string
	FirstOp        Action
	SecondOp       Action
	Weight         int
}

type ConstraintGraph struct {
	internal graph.Graph[Constraint, string]
}

func (c *ConstraintGraph) AddBinaryConstraint(cons Constraint) {
	from := graph.NewVertex[string]()
	to := graph.NewVertex[string]()
	from.SetID(cons.FirstNode)
	to.SetID(cons.SecondNode)

	edge := graph.NewEdge[Constraint, string]()
	edge.SetFrom(from)
	edge.SetTo(to)
	edge.SetProperty(ConsProp, cons)
	edge.SetID(fmt.Sprintf("%s->%s", from, to))
	c.internal.AddEdge(edge)
}

func (c *ConstraintGraph) GetConstraints() []Constraint {
	edges := c.internal.GetAllEdges()
	res := make([]Constraint, 0)
	for _, v := range edges {
		res = append(res, v.GetProperty(ConsProp))
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].Weight > res[j].Weight
	})
	return res
}

func NewConstraintGraph() *ConstraintGraph {
	return &ConstraintGraph{internal: graph.NewGraph[Constraint, string]()}
}
