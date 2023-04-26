package ctxlog_test

import (
	"context"
	"testing"

	"github.com/cyrusaf/ctxlog"
	"golang.org/x/exp/slog"
)

func TestTags(t *testing.T) {
	ctx := context.Background()
	attr1 := slog.String("hello", "world")
	attr2 := slog.Int("foo", 5)
	ctx = ctxlog.WithAttrs(ctx, attr1)
	ctx = ctxlog.WithAttrs(ctx, attr2)
	attrs := ctxlog.GetAttrs(ctx)
	if len(attrs) != 2 {
		t.Fatalf("expected 2 attrs, but got %d instead", len(attrs))
	}
	if !attrs[0].Equal(attr1) {
		t.Fatalf(`expected attr to be %+v but got %+v instead`, attr1, attrs[0])
	}
	if !attrs[1].Equal(attr2) {
		t.Fatalf(`expected attr to be %+v but got %+v instead`, attr2, attrs[1])
	}
}
