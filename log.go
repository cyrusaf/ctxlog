package ctxlog

import (
	"context"
	"log/slog"
	"os"
)

func NewHandler(baseHandler slog.Handler) Handler {
	if baseHandler == nil {
		baseHandler = slog.NewJSONHandler(os.Stdout, nil)
	}
	return Handler{baseHandler}
}

type Handler struct {
	slog.Handler
}

func (t Handler) Handle(ctx context.Context, r slog.Record) error {
	r.AddAttrs(GetAttrs(ctx)...)
	return t.Handler.Handle(ctx, r)
}
