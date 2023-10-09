package schemas

import (
	"context"
	"errors"
)

type Stack struct {
	q     []*Grammar
	trace []*Grammar
}

func (q *Stack) Push(g *Grammar) *Stack {
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
func (q *Stack) Top() *Grammar {
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
func (q *Stack) GetTrace() []*Grammar {
	return q.trace
}

func NewStack() *Stack {
	return &Stack{q: make([]*Grammar, 0),
		trace: make([]*Grammar, 0)}
}

type Context struct {
	SymCount   map[string]int
	grammarMap map[string]*Grammar
	context.Context
	HandlerIndex   int
	SymbolStack    *Stack
	ProductionRoot *Grammar
	Result         string
	finish         bool
}

func (c *Context) GetFinish() bool {
	return c.finish
}
func NewContext(grammarMap map[string]*Grammar, startSymbol string) (*Context, error) {
	_, ok := grammarMap[startSymbol]
	if !ok {
		return nil, errors.New("no such symbol in Grammar")
	}
	return &Context{
		SymCount:       map[string]int{},
		grammarMap:     grammarMap,
		SymbolStack:    NewStack().Push(grammarMap[startSymbol]),
		ProductionRoot: grammarMap[startSymbol],
	}, nil
}
