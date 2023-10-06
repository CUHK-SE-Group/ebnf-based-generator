package log

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"strings"
)

type CallerInfoHandler struct {
	innerHandler slog.Handler
}

func (h *CallerInfoHandler) Handle(ctx context.Context, r slog.Record) error {
	pc, file, _, ok := runtime.Caller(3) // Adjust the skip value as needed
	if ok {
		shortFile := file[strings.LastIndex(file, "/")+1:]
		funcName := runtime.FuncForPC(pc).Name()
		shortFuncName := funcName[strings.LastIndex(funcName, ".")+1:]
		r.Message = fmt.Sprintf("%s:%s: %s", shortFile, shortFuncName, r.Message)
	}
	return h.innerHandler.Handle(ctx, r)
}

func (h *CallerInfoHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.innerHandler.Enabled(ctx, level)
}

func (h *CallerInfoHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &CallerInfoHandler{innerHandler: h.innerHandler.WithAttrs(attrs)}
}

func (h *CallerInfoHandler) WithGroup(name string) slog.Handler {
	return &CallerInfoHandler{innerHandler: h.innerHandler.WithGroup(name)}
}

func NewCallerInfoHandler(innerHandler slog.Handler) *CallerInfoHandler {
	return &CallerInfoHandler{innerHandler: innerHandler}
}
