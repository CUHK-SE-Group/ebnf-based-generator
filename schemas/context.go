package schemas

import (
	"context"
	"errors"
)

type Stack struct {
	q     []*Node
	trace []*Node
}

func (q *Stack) Push(g *Node) *Stack {
	if g == nil {
		panic(g)
	}
	q.q = append(q.q, g)
	return q
}

func (q *Stack) Pop() *Stack {
	q.trace = append(q.trace, q.q[len(q.q)-1])
	q.q = q.q[:len(q.q)-1]
	return q
}
func (q *Stack) Top() *Node {
	if len(q.q) > 0 {
		return q.q[len(q.q)-1]
	}
	return nil
}
func (q *Stack) Empty() bool {
	if len(q.q) == 0 {
		return true
	}
	for _, v := range q.q {
		if v != nil {
			return false
		}
	}
	return true
}
func (q *Stack) GetTrace() []*Node {
	return q.trace
}

func NewStack() *Stack {
	return &Stack{q: make([]*Node, 0),
		trace: make([]*Node, 0)}
}

type Context struct {
	SymCount   map[string]int
	grammarMap *Grammar
	context.Context
	HandlerIndex   int
	SymbolStack    *Stack
	ProductionRoot *Node
	Result         string
	finish         bool
}

func (c *Context) GetFinish() bool {
	return c.finish
}
func NewContext(grammarMap *Grammar, startSymbol string) (*Context, error) {
	node := grammarMap.GetNode(startSymbol)
	if node == nil {
		return nil, errors.New("no such symbol")
	}
	return &Context{
		SymCount:       map[string]int{},
		grammarMap:     grammarMap,
		SymbolStack:    NewStack().Push(node),
		ProductionRoot: node,
	}, nil
}
