package schemas

type ResponseCallBack func(*Result)

type Result struct {
	path   []*Grammar
	output []string
	ctx    *Context
}

func NewResult(ctx *Context) *Result {
	return &Result{
		path: []*Grammar{},
		ctx:  ctx,
	}
}

func (r *Result) GetCtx() *Context {
	return r.ctx
}
func (r *Result) AddNode(n *Grammar) *Result {
	r.path = append(r.path, n)
	return r
}

func (r *Result) AddOutput(s string) *Result {
	r.output = append(r.output, s)
	return r
}

func (r *Result) GetPath() []*Grammar {
	return r.path
}

func (r *Result) GetOutput() []string {
	return r.output
}
