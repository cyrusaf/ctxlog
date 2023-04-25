# ctxlog

Wrapper on top of [golang.org/x/exp/slog](https://pkg.go.dev/golang.org/x/exp/slog) that annotates logs with tags set on the context.

## Usage

`ctxlog` can be used as a global logger (recommended) or as an initiazed logger that is passed around.

To use as a global logger you can just invoke the following methods:

```golang
package main

import (
    "github.com/cyrusaf/ctxlog"
)

func main() {
    ctx := context.Background()
    // Can set tags on the context using ctxlog.WithTag(ctx, key, value)
    ctx = ctxlog.WithTag(ctx, "hello", "world")
    // Can also set tags when logging. Uses the same interface as slog
    ctxlog.InfoCtx(ctx, "foo", "bar")
    // {"time":"2023-04-25T19:18:13.457009-07:00","level":"INFO","msg":"test","foo":"bar","hello":"world"}
}
```
