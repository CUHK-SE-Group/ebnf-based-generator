package schemas

type Chain struct {
	ServiceType string
	Name        string
	Handlers    []Handler
}

func (c *Chain) Clone() Chain {
	var clone = Chain{
		ServiceType: c.ServiceType,
		Name:        c.Name,
		Handlers:    make([]Handler, len(c.Handlers)),
	}
	copy(clone.Handlers, c.Handlers)
	return clone
}

// AddHandler chain can add a handler
func (c *Chain) AddHandler(h Handler) {
	c.Handlers = append(c.Handlers, h)
}

// Next is for to handle next handler in the chain
func (c *Chain) Next(i *Context, f ResponseCallBack) {
	index := i.HandlerIndex
	if index >= len(c.Handlers) {
		r := NewResult()
		f(r)
		return
	}
	i.HandlerIndex++
	c.Handlers[index].Handle(c, i, f)
}
