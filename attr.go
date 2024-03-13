package ctxlog

import (
	"context"
	"log/slog"
)

type ctxkey struct {}

func WithAttrs(ctx context.Context, newAttrs ...slog.Attr) context.Context {
	// Get fields if set
	attrs, _ := ctx.Value(ctxkey{}).([]slog.Attr)
	attrs = append(attrs, newAttrs...)
	ctx = context.WithValue(ctx, ctxkey{}, attrs)
	return ctx
}

func GetAttrs(ctx context.Context) []slog.Attr {
	attrs, _ := ctx.Value(ctxkey{}).([]slog.Attr)
	return attrs
}
