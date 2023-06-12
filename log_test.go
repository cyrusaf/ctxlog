package ctxlog_test

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/cyrusaf/ctxlog"
	"golang.org/x/exp/slog"
)

func ExampleHandler() {
	ctx := context.Background()

	// Create a ctxlog and json based logger and set it as the default logger
	handlerOpts := slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Remove time from the output for predictable test output.
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	}

	logger := slog.New(ctxlog.NewHandler(slog.NewJSONHandler(os.Stdout, &handlerOpts)))
	slog.SetDefault(logger)

	// Can set attrs on the context using ctxlog.WithAttrs(ctx, ...slog.Attr)
	ctx = ctxlog.WithAttrs(ctx, slog.String("hello", "world"))

	// Use slog methods such as InfoCtx and the ctxlog handler will automatically
	// attach attrs from the context to the structured logs.
	slog.InfoCtx(ctx, "test")
	// Output:{"level":"INFO","msg":"test","hello":"world"}
}

func TestLogger(t *testing.T) {
	tt := []struct {
		fn    func(context.Context, string, ...any)
		level slog.Level
	}{
		{slog.InfoCtx, slog.LevelInfo},
		{slog.WarnCtx, slog.LevelWarn},
		{slog.ErrorCtx, slog.LevelError},
		{slog.DebugCtx, slog.LevelDebug},
	}

	for _, tc := range tt {
		t.Run(tc.level.String(), func(t *testing.T) {
			ctx := context.Background()
			b := bytes.Buffer{}

			tagHandler := ctxlog.NewHandler(slog.NewJSONHandler(&b, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			}))
			slog.SetDefault(slog.New(tagHandler))

			ctx = ctxlog.WithAttrs(ctx, slog.String("hello", "world"))
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
