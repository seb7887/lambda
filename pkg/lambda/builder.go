package lambda

import "github.com/aws/aws-lambda-go/lambda"

type Builder[I, O any] struct {
	middlewares []MiddlewareFn[I, O]
	fn          HandlerFn[I, O]
}

func NewLambdaBuilder[I, O any](fn HandlerFn[I, O]) *Builder[I, O] {
	return &Builder[I, O]{
		fn: fn,
	}
}

func (b *Builder[I, O]) WithMiddlewares(middlewares ...MiddlewareFn[I, O]) *Builder[I, O] {
	b.middlewares = append(b.middlewares, middlewares...)
	return b
}

func (b *Builder[I, O]) Start() {
	h := NewHandler(b.fn, b.middlewares...)
	lambda.Start(h.EventHandler)
}

func (b *Builder[I, O]) StartBatch() {
	h := NewHandler(b.fn, b.middlewares...)
	h.Batch()
	lambda.Start(h.EventHandler)
}
