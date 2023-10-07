package schemas

import (
	"fmt"
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
		if ctx.SymbolStack.Top() == nil {
			slog.Error("Warning: Symbol queue should not be empty")
			break
		}
		// 如果类型符合
		if ctx.SymbolStack.Top().gtype&c.Handlers[index].Type() != 0 {
			fmt.Println("passing", c.Handlers[index].Name())
			fmt.Println("cur node: ", ctx.SymbolStack.Top().content)
			c.Handlers[index].Handle(c, ctx, f)
			return
		}
	}
	// 最后一个handler已执行结束，则回调
	r := NewResult(ctx)
	f(r)
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
