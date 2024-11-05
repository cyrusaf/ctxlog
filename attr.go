package ctxlog

import (
	"context"
	"log/slog"
	"sync/atomic"
)

type ctxkey struct{}
type globalctxkey struct{}

func WithAttrs(ctx context.Context, newAttrs ...slog.Attr) context.Context {
	// Get fields if set
	oldAttrs, _ := ctx.Value(ctxkey{}).(map[string]slog.Attr)

	// Make new map so we aren't modifying the previous map resulting in a race condition
	attrs := make(map[string]slog.Attr, len(newAttrs)+len(oldAttrs))
	for _, attr := range oldAttrs {
		attrs[attr.Key] = attr
	}
	for _, attr := range newAttrs {
		attrs[attr.Key] = attr
	}
	ctx = context.WithValue(ctx, ctxkey{}, attrs)
	return ctx
}

func GetAttrs(ctx context.Context) []slog.Attr {
	attrMap, _ := ctx.Value(ctxkey{}).(map[string]slog.Attr)
	attrs := make([]slog.Attr, 0, len(attrMap))
	for _, attr := range attrMap {
		attrs = append(attrs, attr)
	}
	globalAttrMap, ok := ctx.Value(globalctxkey{}).(*atomic.Pointer[map[string]slog.Attr])
	if ok && globalAttrMap != nil {
		for _, attr := range *globalAttrMap.Load() {
			attrs = append(attrs, attr)
		}
	}
	return attrs
}

// Global attributes are used to modify the attributes logged by a parent. One
// example of where this might be useful: Logging middleware that is outside the
// scope of the request handler:
//
// ```
//
//	func LogMiddleware(ctx context.Context, h Handler) {
//	  ctx = ctxlog.InitGlobalAttrs(ctx)
//	  err := h(ctx)
//	  if err != nil {
//	       slog.ErrorContext(ctx, "request failed")
//	  }
//	}
//
//	func myHandler(ctx context.Context) error {
//	  ctx = ctxlog.WithGlobalAttrs(ctx, slog.String("request_id", requestId))
//	  // ...
//	}
//
// ```

// AnchorGlobalAttrs is called to set the "anchor" or root of the global attrs
// attached to the context. `WithGlobalAttrs` will modify the global attrs up to
// this point.
//
// If AnchorGlobalAttrs is called on a context with an existing anchor point, it
// will act as noop and use the original anchor point.
func AnchorGlobalAttrs(ctx context.Context) context.Context {
	_, ok := ctx.Value(globalctxkey{}).(*atomic.Pointer[map[string]slog.Attr])
	if ok {
		return ctx
	}
	ctx, _ = initGlobalAttrs(ctx)
	return ctx
}

// WithGlobalAttrs attaches "global" attributes to the context that can be read
// by a parent function. It is useful when a sub-function parses out data that
// should be logged by the parent later on. For example, logging middleware that
// exists above the scope of the request handler.
//
// If AnchorGlobalAttrs has not been called yet for the given context, the
// returned context will be set as the anchor point.
func WithGlobalAttrs(ctx context.Context, newAttrs ...slog.Attr) context.Context {
	attrsAtomic, ok := ctx.Value(globalctxkey{}).(*atomic.Pointer[map[string]slog.Attr])
	if !ok || attrsAtomic == nil {
		ctx, attrsAtomic = initGlobalAttrs(ctx)
	}

	// FIXME: Should we use a mutex instead of atomics to prevent a gap in read
	// -> write that can lead to missing attrs?
	oldAttrs := attrsAtomic.Load()
	if oldAttrs == nil {
		attrsNew := make(map[string]slog.Attr, len(newAttrs))
		oldAttrs = &attrsNew
	}

	attrs := make(map[string]slog.Attr, len(newAttrs)+len(*oldAttrs))
	for _, attr := range *oldAttrs {
		attrs[attr.Key] = attr
	}
	for _, attr := range newAttrs {
		attrs[attr.Key] = attr
	}
	attrsAtomic.Store(&attrs)

	return ctx
}

func initGlobalAttrs(ctx context.Context) (context.Context, *atomic.Pointer[map[string]slog.Attr]) {
	m := make(map[string]slog.Attr)
	globalAttrs := atomic.Pointer[map[string]slog.Attr]{}
	globalAttrs.Store(&m)
	ctx = context.WithValue(ctx, globalctxkey{}, &globalAttrs)
	return ctx, &globalAttrs
}
