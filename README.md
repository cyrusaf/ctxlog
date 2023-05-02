# ctxlog

[![Go Reference](https://pkg.go.dev/badge/github.com/cyrusaf/ctxlog.svg)](https://pkg.go.dev/github.com/cyrusaf/ctxlog)

Handler for [golang.org/x/exp/slog](https://pkg.go.dev/golang.org/x/exp/slog)
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
    "golang.org/x/exp/slog"
    "github.com/cyrusaf/ctxlog"
)

func main() {
 ctx := context.Background()

 // Create ctxlog and json logger and set it as the default logger
 logger := slog.New(ctxlog.NewHandler(slog.NewJSONHandler(os.Stdout)))
 slog.SetDefault(logger)

 // Can set attrs on the context using ctxlog.WithAttrs(ctx, ...slog.Attr)
 ctx = ctxlog.WithAttrs(ctx, slog.String("hello", "world"))

 // Use slog methods such as InfoCtx and the ctxlog handler will automatically
 // attach attrs from the context to the structured logs.
 slog.InfoCtx(ctx, "test")
 // Output:{"level":"INFO","msg":"test","hello":"world"}
}
```
