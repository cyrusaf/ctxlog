package ctxlog

import (
	"context"
	"os"

	"golang.org/x/exp/slog"
)

var (
	// Logger can be overwritten by the consumer if they want to change the
	// settings. It can also be used directly in case the global log methods
	// don't expose all functionality needed.
	Logger *slog.Logger = slog.New(TagHandler{slog.NewJSONHandler(os.Stdout)})
)

func InfoCtx(ctx context.Context, msg string, args ...any) {
	Logger.InfoCtx(ctx, msg, args...)
}

func DebugCtx(ctx context.Context, msg string, args ...any) {
	Logger.DebugCtx(ctx, msg, args...)
}

func WarnCtx(ctx context.Context, msg string, args ...any) {
	Logger.WarnCtx(ctx, msg, args...)
}

func ErrorCtx(ctx context.Context, msg string, args ...any) {
	Logger.ErrorCtx(ctx, msg, args...)
}

func LogCtx(ctx context.Context, level slog.Level, msg string, args ...any) {
	Logger.Log(ctx, level, msg, args...)
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
