package ctxlog_test

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/cyrusaf/ctxlog"
	"golang.org/x/exp/slog"
)

func TestLogger(t *testing.T) {
	tt := []struct {
		fn    func(context.Context, string, ...any)
		level slog.Level
	}{
		{ctxlog.InfoCtx, slog.LevelInfo},
		{ctxlog.WarnCtx, slog.LevelWarn},
		{ctxlog.ErrorCtx, slog.LevelError},
		{ctxlog.DebugCtx, slog.LevelDebug},
	}

	for _, tc := range tt {
		t.Run(tc.level.String(), func(t *testing.T) {
			ctx := context.Background()
			b := bytes.Buffer{}

			jsonHandler := slog.HandlerOptions{
				Level: slog.LevelDebug,
			}.NewJSONHandler(&b)
			ctxlog.Logger = slog.New(ctxlog.TagHandler{jsonHandler})

			ctx = ctxlog.WithTag(ctx, "hello", "world")
			tc.fn(ctx, "test", "foo", "bar")

			l := make(map[string]interface{})
			json.Unmarshal(b.Bytes(), &l)
			if len(l) != 5 {
				t.Fatalf("expected 5 keys in log line, but got %d instead", len(l))
			}
			if val, ok := l["level"].(string); !ok || val != tc.level.String() {
				t.Fatalf(`expect level to be %q but got "%v instead"`, tc.level.String(), l["level"])
			}
			if val, ok := l["msg"].(string); !ok || val != "test" {
				t.Fatalf(`expect msg to be "test" but got "%v" instead"`, l["msg"])
			}
			if val, ok := l["hello"].(string); !ok || val != "world" {
				t.Fatalf(`expect hello tag to be "world" but got "%v" instead"`, l["hello"])
			}
			if val, ok := l["foo"].(string); !ok || val != "bar" {
				t.Fatalf(`expect foo tag to be "bar" but got "%v" instead"`, l["foo"])
			}
		})
	}
}
