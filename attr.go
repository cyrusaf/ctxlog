package ctxlog

import (
	"context"
	"log/slog"
)

type ctxkey struct{}
type globalctxkey struct{}

func WithAttrs(ctx context.Context, newAttrs ...slog.Attr) context.Context {
	// Get fields if set
	attrs, _ := ctx.Value(ctxkey{}).([]slog.Attr)
	attrs = append(attrs, newAttrs...)
	ctx = context.WithValue(ctx, ctxkey{}, attrs)
	return ctx
}

func GetAttrs(ctx context.Context) []slog.Attr {
	attrs, _ := ctx.Value(ctxkey{}).([]slog.Attr)
	globalAttrs, ok := ctx.Value(globalctxkey{}).(*[]slog.Attr)
	if ok && globalAttrs != nil {
		attrs = append(attrs, *globalAttrs...)
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
// AnchorGlobalAttrs can be used to anchor global attrs to a given request,
// rather than the context created at the beginning of the server.
func AnchorGlobalAttrs(ctx context.Context) context.Context {
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
	attrs, ok := ctx.Value(globalctxkey{}).(*[]slog.Attr)
	if !ok {
		ctx, attrs = initGlobalAttrs(ctx)
	}
	*attrs = append(*attrs, newAttrs...)
	return ctx
}

func initGlobalAttrs(ctx context.Context) (context.Context, *[]slog.Attr) {
	globalAttrs := &[]slog.Attr{}
	ctx = context.WithValue(ctx, globalctxkey{}, globalAttrs)
	return ctx, globalAttrs
}
