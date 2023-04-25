package ctxlog_test

import (
	"context"
	"testing"

	"github.com/cyrusaf/ctxlog"
)

func TestTags(t *testing.T) {
	ctx := context.Background()
	ctx = ctxlog.WithTag(ctx, "hello", "world")
	ctx = ctxlog.WithTag(ctx, "foo", "bar")
	tags := ctxlog.GetTags(ctx)
	if len(tags) != 2 {
		t.Fatalf("expected 2 tags, but got %d instead", len(tags))
	}
	if tag, ok := tags["hello"].(string); !ok || tag != "world" {
		t.Fatalf(`expected value to be "hello" but got "%v" instead`, tags["hello"])
	}
	if tag, ok := tags["foo"].(string); !ok || tag != "bar" {
		t.Fatalf(`expected value to be "hello" but got "%v" instead`, tags["foo"])
	}
}
