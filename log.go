package ctxlog

import (
	"context"
	"os"

	"golang.org/x/exp/slog"
)

func NewTagHandler(baseHandler slog.Handler) TagHandler {
	if baseHandler == nil {
		baseHandler = slog.NewJSONHandler(os.Stdout)
	}
	return TagHandler{baseHandler}
}

type TagHandler struct {
	slog.Handler
}

func (t TagHandler) Handle(ctx context.Context, r slog.Record) error {
	r.AddAttrs(tagAttrSlice(GetTags(ctx))...)
	return t.Handler.Handle(ctx, r)
}

func tagAttrSlice(tags map[string]interface{}) []slog.Attr {
	s := make([]slog.Attr, 0, len(tags))
	for key, value := range tags {
		s = append(s, slog.Any(key, value))
	}
	return s
}
