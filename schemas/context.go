package schemas

import (
	"context"
	"errors"
)

type Stack struct {
	q []*Grammar
}

func (q *Stack) Push(g *Grammar) *Stack {
	q.q = append(q.q, g)
	return q
}

func (q *Stack) Pop() *Stack {
	q.q = q.q[:len(q.q)-1]
	return q
}
func (q *Stack) Top() *Grammar {
	if len(q.q) > 0 {
		return q.q[len(q.q)-1]
	}
	return nil
}

func NewQueue() *Stack {
	return &Stack{make([]*Grammar, 0)}
}

type Context struct {
	SymCount   map[string]int
	grammarMap map[string]*Grammar
	context.Context
	HandlerIndex   int
	SymbolStack    *Stack
	ProductionRoot *Grammar
}

func NewContext(grammarMap map[string]*Grammar, startSymbol string) (*Context, error) {
	_, ok := grammarMap[startSymbol]
	if !ok {
		return nil, errors.New("no such symbol in Grammar")
	}
	return &Context{
		SymCount:       map[string]int{},
		grammarMap:     grammarMap,
		SymbolStack:    NewQueue().Push(grammarMap[startSymbol]),
		ProductionRoot: grammarMap[startSymbol],
	}, nil
}
