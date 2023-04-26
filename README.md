# ctxlog

Wrapper on top of [golang.org/x/exp/slog](https://pkg.go.dev/golang.org/x/exp/slog) that annotates logs with tags set on the context.

## Usage

`ctxlog` can be used as a global logger (recommended) by using `slog.SetDefault()`.

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
