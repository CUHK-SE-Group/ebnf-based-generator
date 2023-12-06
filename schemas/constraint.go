package schemas

import (
	"regexp"
)

type ConsAction int

const (
	FUNC  ConsAction = iota + 1
	REGEX ConsAction = iota + 1
)

/*
How to define a constraint?
1. There will always be a pre-condition of a constraint.
    if a symbol(or any other information is found in FirstPlaceNode, it needs to be recorded in schemas.Context)
2. Then, in the next step, if a symbol found that it is in SecondPlaceNode, it will read the information that stored in schemas.Context, then take corresponding action
*/

var FirstPlaceNode map[string]func(ctx *Context) *Context // key's type is path pattern
var SecondPlaceNode map[string]Action

func init() {
	FirstPlaceNode = make(map[string]func(ctx *Context) *Context)
	SecondPlaceNode = make(map[string]Action)
}

type Action struct {
	Type   ConsAction
	Regexp *regexp.Regexp
	Func   func(ctx *Context) *Context
}

type Constraint struct {
	FirstNode  string
	SecondNode string
	FirstOp    func(ctx *Context) *Context
	SecondOp   Action
}

func Register(cons ...Constraint) {
	for _, c := range cons {
		FirstPlaceNode[c.FirstNode] = c.FirstOp
		SecondPlaceNode[c.SecondNode] = c.SecondOp
	}
}
