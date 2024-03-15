module github.com/cyrusaf/ctxlog

go 1.21

retract v1.3.1 // panic if WithAttrs or WithGlobalAttrs called concurrently
