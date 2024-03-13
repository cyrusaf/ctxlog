# ctxlog

[![Go Reference](https://pkg.go.dev/badge/github.com/cyrusaf/ctxlog.svg)](https://pkg.go.dev/github.com/cyrusaf/ctxlog)

Handler for [log/slog](https://pkg.go.dev/log/slog)
that annotates logs with `slog.Attr` set on the context. Provides methods for
adding `slog.Attr` to the context and a `slog.Handler` for automatically reading
them from the context and adding them to log lines. Useful for adding fields
such as `request_id`, `xray_trace_id`, or `caller` to log lines.

## Usage

Use `ctxlog.WithAttrs(ctx, attrs...)` to add `slog.Attr` to the context. Use
`ctxlog.NewHandler(baseHandler)` to create a new `slog.Handler` that reads attrs
from the context and adds them to log lines automatically. 

```golang
package main

import (
    "context"
    "log/slog"
    "github.com/cyrusaf/ctxlog"
)

func main() {
 ctx := context.Background()

 // Create ctxlog and json logger and set it as the default logger
 logger := slog.New(ctxlog.NewHandler(slog.NewJSONHandler(os.Stdout)))
 slog.SetDefault(logger)

 // Can set attrs on the context using ctxlog.WithAttrs(ctx, ...slog.Attr)
 ctx = ctxlog.WithAttrs(ctx, slog.String("hello", "world"))

 // Use slog methods such as InfoContext and the ctxlog handler will automatically
 // attach attrs from the context to the structured logs.
 slog.InfoContext(ctx, "test")
 // Output:{"level":"INFO","msg":"test","hello":"world"}
}
```

## Global Attributes

Sometimes you want a sub-function to be able to "pass back" attributes to the
parent. An example of this is when using middleware to log errors. In the
situation below, the handler `h` may handle parsing the `req` and pull out
certain fields that should be logged if there is an error. Unfortunately,
because `WithAttrs(...)` only attaches the attrs to a child context, the parent
function will not have access to these attrs.

```golang
func logMiddleware(ctx context.Context, h Handler, req []byte) {
    err := h(ctx, req)
    if err != nil {
        slog.ErrorContext(ctx, "request error")
    }
}
```

`ctxlog` provides a way to pass back these attributes for logging through
__global attributes__.

```golang
func logMiddleware(ctx context.Context, h Handler, req []byte) {
    // First, set the anchor point/root of global attrs. This allows us to
    // scope global attrs to each request.
    ctx = ctxlog.AnchorGlobalAttrs(ctx) 
    err := h(ctx, req)
    if err != nil {
        slog.ErrorContext(ctx, "request error")
    }
}

func myHandler(ctx context.Context, req []byte) error {
    parsedReq := parseRequest(req)

    // Use ctxlog.WithGlobalAttrs to "pass back" attrs to the anchor point
    ctx = ctxlog.WithGlobalAttrs(ctx, slog.String("request_id", parsedReq.request_id))
    // ...
    return nil
}
```
