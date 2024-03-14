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

	attrs := ctxlog.GetAttrs(ctx)
	if len(attrs) != 4 {
		t.Fatalf("expected 4 attrs, but got %d instead", len(attrs))
	}
	if !attrs[0].Equal(attr1) {
		t.Fatalf(`expected attr to be %+v but got %+v instead`, attr1, attrs[0])
	}
	if !attrs[1].Equal(attr2) {
		t.Fatalf(`expected attr to be %+v but got %+v instead`, attr2, attrs[1])
	}
	if !attrs[2].Equal(attr3) {
		t.Fatalf(`expected attr to be %+v but got %+v instead`, attr3, attrs[2])
	}
	if !attrs[3].Equal(attr4) {
		t.Fatalf(`expected attr to be %+v but got %+v instead`, attr4, attrs[3])
	}
}

func TestWithGlobalAttrsWithoutAnchor(t *testing.T) {
	ctx := context.Background()
	attr1 := slog.String("foo", "bar")
	ctx2 := ctxlog.WithGlobalAttrs(ctx, attr1)

	attrs := ctxlog.GetAttrs(ctx2)
	if len(attrs) != 1 {
		t.Fatalf("expected 1 attr, but got %d instead", len(attrs))
	}
	if !attrs[0].Equal(attr1) {
		t.Fatalf(`expected attr to be %+v but got %+v instead`, attr1, attrs[0])
	}

	attrs = ctxlog.GetAttrs(ctx)
	if len(attrs) != 0 {
		t.Fatalf("expected 0 attrs, but got %d instead", len(attrs))
	}
}
