package schemas

import (
	"log/slog"
)

type Chain struct {
	Name     string
	Handlers []Handler
}

func (c *Chain) Clone() Chain {
	var clone = Chain{
		Name:     c.Name,
		Handlers: make([]Handler, len(c.Handlers)),
	}
	copy(clone.Handlers, c.Handlers)
	return clone
}

// AddHandler chain can add a handler
func (c *Chain) AddHandler(h Handler) {
	c.Handlers = append(c.Handlers, h)
}

// Next is for to handle next handler in the chain
func (c *Chain) Next(ctx *Context, f ResponseCallBack) {
	for index := ctx.HandlerIndex; index < len(c.Handlers); index++ {
		ctx.HandlerIndex++
		if ctx.SymbolStack.Top() == nil || ctx.SymbolStack.Empty() {
			if !ctx.SymbolStack.Empty() {
				panic(ctx)
			}
			if ctx.finish {
				slog.Error("Warning: Symbol queue should not be empty")
			}
			ctx.finish = true
			r := NewResult(ctx)
			f(r)
			return
		}

		// 如果类型符合
		if ctx.SymbolStack.Top().GetType()&c.Handlers[index].Type() != 0 && satisfy(ctx, c.Handlers[index]) {
			c.Handlers[index].Handle(c, ctx, f)
		}
	}
	r := NewResult(ctx)
	f(r)
}

func satisfy(ctx *Context, handler Handler) bool {
	return true
}
func CreateChain(chainName string, handlers ...Handler) (*Chain, error) {
	c := &Chain{
		Name: chainName,
	}
	for _, h := range handlers {
		c.AddHandler(h)
	}

	return c, nil
}
