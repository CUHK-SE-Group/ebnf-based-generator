package schemas

import (
	"context"
	"errors"
	"strings"

	"github.com/hashicorp/go-memdb"
)

type Mode int

const (
	ShrinkMode Mode = iota + 1
)

type Stack struct {
	q               []*Node
	trace           []*Node
	ProductionTrace []string
}

func (q *Stack) Push(g ...*Node) *Stack {
	if g == nil {
		panic(g)
	}
	q.q = append(q.q, g...)
	return q
}

func (q *Stack) Pop() *Stack {
	curSym := q.q[len(q.q)-1]
	q.trace = append(q.trace, curSym)
	lastSym := ""
	if len(q.ProductionTrace) > 0 {
		lastSym = q.ProductionTrace[len(q.ProductionTrace)-1]
	}

	if len(q.ProductionTrace) == 0 || lastSym == curSym.GetID() {
		q.ProductionTrace = append(q.ProductionTrace, strings.TrimSpace(curSym.GetID()))
	}
	if lastSym != curSym.GetID() && !strings.HasPrefix(curSym.GetID(), lastSym) {
		q.ProductionTrace = append(q.ProductionTrace, strings.Split(strings.TrimSpace(curSym.GetID()), "#")[0])
	}
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
func (q *Stack) GetStack() []*Node {
	return q.q
}
func NewStack() *Stack {
	return &Stack{
		q:               make([]*Node, 0),
		trace:           make([]*Node, 0),
		ProductionTrace: make([]string, 0),
	}
}

type Context struct {
	Grammar *Grammar
	context.Context
	HandlerIndex int
	SymbolStack  *Stack
	Result       *Derivation
	finish       bool
	Storage      *memdb.MemDB

	VisitedEdge    map[string]int
	Mode           Mode
	Constraint     *ConstraintGraph
	MemoryExchange map[string]int
}

type NodeRuntimeInfo struct {
	Count        int
	ID           string
	SampledValue map[string]int
}

func (c *Context) GetFinish() bool {
	return c.finish
}
func NewContext(grammarMap *Grammar, startSymbol string, ctx context.Context, cons *ConstraintGraph, gendb func() *memdb.MemDB) (*Context, error) {
	node := grammarMap.GetNode(startSymbol)
	if node == nil {
		return nil, errors.New("no such symbol")
	}
	var db *memdb.MemDB
	var err error
	if gendb == nil {
		schema := &memdb.DBSchema{
			Tables: map[string]*memdb.TableSchema{
				"nodeRuntimeInfo": &memdb.TableSchema{
					Name: "nodeRuntimeInfo",
					Indexes: map[string]*memdb.IndexSchema{
						"id": &memdb.IndexSchema{
							Name:    "id",
							Unique:  true,
							Indexer: &memdb.StringFieldIndex{Field: "ID"},
						},
						"count": &memdb.IndexSchema{
							Name:    "count",
							Unique:  false,
							Indexer: &memdb.IntFieldIndex{Field: "Count"},
						},
					},
				},
			},
		}
		db, err = memdb.NewMemDB(schema)
		if err != nil {
			panic(err)
		}
	} else {
		db = gendb()
	}

	return &Context{
		Grammar:     grammarMap,
		Context:     ctx, // 使用带有超时的context
		SymbolStack: NewStack().Push(node),
		Storage:     db,
		VisitedEdge: map[string]int{},
		Result: &Derivation{
			Grammar:     NewGrammar(WithStartSym(startSymbol)),
			EdgeHistory: make([]string, 0),
			SymbolCnt:   make(map[string]int),
		},
		Constraint:     cons,
		MemoryExchange: map[string]int{},
	}, nil
}
