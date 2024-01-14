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
	if len(middlewares) > 0 {
		b.middlewares = append(b.middlewares, middlewares...)
	}
	return b
}

func (b *Builder[I, O]) Start() {
	lambda.Start(NewHandler(b.fn, b.middlewares...).EventHandler)
}
