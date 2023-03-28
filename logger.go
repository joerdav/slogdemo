package main

import (
	"context"
	"log/slog"
)

type contextLogger struct {
	handler slog.Handler
}

func (cl *contextLogger) Enabled(ctx context.Context, l slog.Level) bool {
	return cl.handler.Enabled(ctx, l)
}

func (cl *contextLogger) Handle(ctx context.Context, r slog.Record) error {
	rc := r.Clone()
	c, ok := ctx.Value(logValuesKey).([]slog.Attr)
	if ok {
		rc.AddAttrs(c...)
	}
	return cl.handler.Handle(ctx, rc)
}

func (cl *contextLogger) WithGroup(name string) slog.Handler {
	return &contextLogger{cl.handler.WithGroup(name)}
}

func (cl *contextLogger) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &contextLogger{cl.handler.WithAttrs(attrs)}
}

type contextKeyType int

const logValuesKey = contextKeyType(0)

func appendAttrs(ctx context.Context, args ...any) context.Context {
	var (
		attr  slog.Attr
		attrs []slog.Attr
	)
	for len(args) > 0 {
		attr, args = argsToAttr(args)
		attrs = append(attrs, attr)
	}
	c, ok := ctx.Value(logValuesKey).([]slog.Attr)
	if !ok {
		c = []slog.Attr{}
	}
	return context.WithValue(ctx, logValuesKey, append(c, attrs...))
}

func clearAttrs(ctx context.Context) context.Context {
	return context.WithValue(ctx, logValuesKey, nil)
}

const badKey = "!BADKEY"

func argsToAttr(args []any) (slog.Attr, []any) {
	switch x := args[0].(type) {
	case string:
		if len(args) == 1 {
			return slog.String(badKey, x), nil
		}
		a := slog.Any(x, args[1])
		a.Value = a.Value.Resolve()
		return a, args[2:]

	case slog.Attr:
		x.Value = x.Value.Resolve()
		return x, args[1:]

	default:
		return slog.Any(badKey, x), args[1:]
	}
}
