package schemas

import "errors"

var errViolateBuildIn = errors.New("can not replace build-in handler func")
var buildIn = []string{}
var ErrDuplicatedHandler = errors.New("duplicated handler registration")
var funcMap = make(map[GrammarType][]func() Handler)

type Handler interface {
	Handle(*Chain, *Context, ResponseCallBack)
	HookRoute() []string
	Name() string
}

func RegisterHandler(name GrammarType, f func() Handler) error {
	_, ok := funcMap[name]
	if ok {
		return ErrDuplicatedHandler
	}
	funcMap[name] = append(funcMap[name], f)
	return nil
}
