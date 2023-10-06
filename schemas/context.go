package schemas

import "context"

type Context struct {
	SymCount map[string]int
	context.Context
	HandlerIndex int
}

func NewContext() *Context {
	return &Context{
		SymCount: map[string]int{},
	}
}

type Config struct {
	SymbolLimit map[string]int
}
