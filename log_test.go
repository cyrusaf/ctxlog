package ctxlog_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"testing"

	"github.com/cyrusaf/ctxlog"
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

	// Use slog methods such as InfoContext and the ctxlog handler will automatically
	// attach attrs from the context to the structured logs.
	slog.InfoContext(ctx, "test")
	// Output:{"level":"INFO","msg":"test","hello":"world"}
}

func handleRequest(ctx context.Context, requestID string) {
	_ = ctxlog.WithGlobalAttrs(ctx, slog.String("request_id", requestID))
}

func ExampleGlobalAttributes() {
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

	// First, set the anchor point of the global attributes
	ctx = ctxlog.AnchorGlobalAttrs(ctx)

	// handleRequest uses global attributes to pass back logging attributes to
	// the anchor point
	handleRequest(ctx, "foo")

	// Use slog methods such as InfoContext and the ctxlog handler will automatically
	// attach global attrs from the context to the structured logs.
	slog.InfoContext(ctx, "test")
	// Output:{"level":"INFO","msg":"test","request_id":"foo"}
}

func TestLogger(t *testing.T) {
	tt := []struct {
		fn    func(context.Context, string, ...any)
		level slog.Level
	}{
		{slog.InfoContext, slog.LevelInfo},
		{slog.WarnContext, slog.LevelWarn},
		{slog.ErrorContext, slog.LevelError},
		{slog.DebugContext, slog.LevelDebug},
	}

	for _, tc := range tt {
		t.Run(tc.level.String(), func(t *testing.T) {
			ctx := context.Background()
			ctx = ctxlog.AnchorGlobalAttrs(ctx)
			b := bytes.Buffer{}

			tagHandler := ctxlog.NewHandler(slog.NewJSONHandler(&b, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			}))
			slog.SetDefault(slog.New(tagHandler).With("baz", "mumble"))

			ctx = ctxlog.WithAttrs(ctx, slog.String("hello", "world"))

			// Test global attrs by dropping the returned ctx
			_ = ctxlog.WithGlobalAttrs(ctx, slog.String("global", "attr"))

			tc.fn(ctx, "test", "foo", "bar")

			l := make(map[string]interface{})
			_ = json.Unmarshal(b.Bytes(), &l)
			if len(l) != 7 {
				t.Fatalf("expected 7 keys in log line, but got %d instead", len(l))
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
			if val, ok := l["baz"].(string); !ok || val != "mumble" {
				t.Fatalf(`expect baz tag to be "mumble" but got "%v" instead"`, val)
			}
			if val, ok := l["global"].(string); !ok || val != "attr" {
				t.Fatalf(`expect global tag to be "attr" but got "%v" instead"`, val)
			}
		})
	}
}
