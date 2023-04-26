package ctxlog

import (
	"context"

	"golang.org/x/exp/slog"
)

type ctxkey string

const (
	attrsKey ctxkey = "ctxlogattrs"
)

func WithAttrs(ctx context.Context, newAttrs ...slog.Attr) context.Context {
	// Get fields if set
	attrs, _ := ctx.Value(attrsKey).([]slog.Attr)
	attrs = append(attrs, newAttrs...)
	ctx = context.WithValue(ctx, attrsKey, attrs)
	return ctx
}

func GetAttrs(ctx context.Context) []slog.Attr {
	attrs, _ := ctx.Value(attrsKey).([]slog.Attr)
	return attrs
}
