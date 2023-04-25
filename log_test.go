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
	b := bytes.Buffer{}
	ctxlog.Logger = slog.New(ctxlog.TagHandler{slog.NewJSONHandler(&b)})

	ctx := context.Background()
	ctx = ctxlog.WithTag(ctx, "hello", "world")
	ctxlog.InfoCtx(ctx, "test", "foo", "bar")

	l := make(map[string]interface{})
	json.Unmarshal(b.Bytes(), &l)
	if len(l) != 5 {
		t.Fatalf("expected 5 keys in log line, but got %d instead", len(l))
	}
	if val, ok := l["level"].(string); !ok || val != "INFO" {
		t.Fatalf(`expect level to be "INFO" but got "%v instead"`, l["level"])
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

}
