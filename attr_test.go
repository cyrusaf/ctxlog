package ctxlog_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/cyrusaf/ctxlog"
)

func TestTags(t *testing.T) {
	ctx := context.Background()
	attr1 := slog.String("hello", "world")
	attr2 := slog.Int("foo", 5)
	ctx = ctxlog.WithAttrs(ctx, attr1)
	ctx = ctxlog.WithAttrs(ctx, attr2)

	// Test global attrs. Use _ = to drop returned context so the global flow
	// is properly tested.
	ctx = ctxlog.AnchorGlobalAttrs(ctx)
	attr3 := slog.Int("bar", 1)
	attr4 := slog.Int("baz", 2)
	_ = ctxlog.WithGlobalAttrs(ctx, attr3)
	_ = ctxlog.WithGlobalAttrs(ctx, attr4)

	expectedAttrs := []slog.Attr{attr1, attr2, attr3, attr4}

	attrs := ctxlog.GetAttrs(ctx)
	assertAttrs(t, expectedAttrs, attrs)
}

func TestWithGlobalAttrsWithoutAnchor(t *testing.T) {
	ctx := context.Background()
	attr1 := slog.String("foo", "bar")
	ctx2 := ctxlog.WithGlobalAttrs(ctx, attr1)

	attrs := ctxlog.GetAttrs(ctx2)
	assertAttrs(t, []slog.Attr{attr1}, attrs)

	attrs = ctxlog.GetAttrs(ctx)
	if len(attrs) != 0 {
		t.Fatalf("expected 0 attrs, but got %d instead", len(attrs))
	}
}

func TestOverwriteAttr(t *testing.T) {
	ctx := context.Background()
	attr1 := slog.String("foo", "bar")
	ctx = ctxlog.WithAttrs(ctx, attr1)

	attr2 := slog.String("foo", "baz")
	ctx = ctxlog.WithAttrs(ctx, attr2)

	attrs := ctxlog.GetAttrs(ctx)
	assertAttrs(t, []slog.Attr{attr2}, attrs)
}

func assertAttrs(t *testing.T, expected, actual []slog.Attr) {
	t.Helper()

	if len(expected) != len(actual) {
		t.Fatalf("expected %d attrs, but got %d instead", len(expected), len(actual))
	}

	expectedMap := toMap(expected)
	actualMap := toMap(actual)
	for key, expected := range expectedMap {
		actual, ok := actualMap[key]
		if !ok {
			t.Fatalf("missing attr %+v", expected)
		}
		if !expected.Equal(actual) {
			t.Fatalf("attrs not equal - expected: %+v, actual: %+v", expected, actual)
		}
	}

}

func toMap(attrs []slog.Attr) map[string]slog.Attr {
	m := make(map[string]slog.Attr, len(attrs))
	for _, attr := range attrs {
		m[attr.Key] = attr
	}
	return m
}
