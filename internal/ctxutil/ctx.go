package ctxutil

import (
	"context"
)

type key[T any] struct{}

func Set[T any](ctx context.Context, value T) context.Context {
	return context.WithValue(ctx, key[T]{}, value)
}

func Get[T any](ctx context.Context) (T, bool) {
	v, ok := ctx.Value(key[T]{}).(T)
	return v, ok
}
