package ctxlog

import (
	"context"
)

type ctxkey string

const (
	tagsKey ctxkey = "ctxlogtags"
)

func WithTag(ctx context.Context, key string, value any) context.Context {
	// Get fields if set
	tags, ok := ctx.Value(tagsKey).(map[string]any)
	if tags == nil || !ok {
		tags = make(map[string]interface{})
	}

	tags[key] = value
	ctx = context.WithValue(ctx, tagsKey, tags)
	return ctx
}

func GetTags(ctx context.Context) map[string]interface{} {
	tags, ok := ctx.Value(tagsKey).(map[string]any)
	if tags == nil || !ok {
		return make(map[string]interface{})
	}
	return tags
}
