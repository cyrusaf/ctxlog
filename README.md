# ctxlog

[![Go Reference](https://pkg.go.dev/badge/github.com/cyrusaf/ctxlog.svg)](https://pkg.go.dev/github.com/cyrusaf/ctxlog)

Handler for [golang.org/x/exp/slog](https://pkg.go.dev/golang.org/x/exp/slog)
that annotates logs with `slog.Attr` set on the context. Provides methods for
adding `slog.Attr` to the context and a `slog.Handler` for automatically reading
them from the context and adding them to log lines. Useful for adding fields
such as `request_id`, `xray_trace_id`, or `caller` to log lines.

## Usage

`ctxlog` can be used as a global logger (recommended) by using `slog.SetDefault()`.

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

 // Create a tag and json based logger and set it as the default logger
 logger := slog.New(ctxlog.NewHandler(slog.NewJSONHandler(os.Stdout)))
 slog.SetDefault(logger)

 // Can set tags on the context using ctxlog.WithTag(ctx, key, value)
 ctx = ctxlog.WithAttrs(ctx, slog.String("hello", "world"))

 // Can also set tags when logging. Can use slog global methods such as
 // InfoCtx if set as default logger.
 slog.InfoCtx(ctx, "test", slog.Int("foo", 5))
 // Output:{"level":"INFO","msg":"test","foo":5,"hello":"world"}}
```
