package schemas

import (
	"github.com/CUHK-SE-Group/generic-generator/graph"
	"regexp"
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

//var FirstPlaceNode map[string]Action // key's type is path pattern
//var SecondPlaceNode map[string]Action

func init() {
	//FirstPlaceNode = make(map[string]func(ctx *Context) (*Context, error))
	//SecondPlaceNode = make(map[string]Action)
}

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
}

func Register(cons ...Constraint) {
	//for _, c := range cons {
	//	FirstPlaceNode[c.FirstNode] = c.FirstOp
	//	SecondPlaceNode[c.SecondNode] = c.SecondOp
	//}
}

type ConstraintGraph struct {
	internal graph.Graph[Constraint, string]
}

func (c *ConstraintGraph) AddConstraint(cons Constraint) {
	from := graph.NewVertex[string]()
	to := graph.NewVertex[string]()
	from.SetID(cons.FirstNode)
	to.SetID(cons.SecondNode)

	if cons.FirstNode == "" || cons.SecondNode == "" {
		if cons.FirstNode != "" {

		}
	}

	edge := graph.NewEdge[Constraint, string]()
	edge.SetFrom(from)
	edge.SetTo(to)
	edge.SetProperty(ConsProp, cons)
	c.internal.AddEdge(edge)

}
