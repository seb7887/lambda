package lambda

func Start[I, O any](fn HandlerFn[I, O], middlewares ...MiddlewareFn[I, O]) {
	buildLambda[I, O](fn, middlewares...).Start()
}

func StartBatch[I, O any](fn HandlerFn[I, O], middlewares ...MiddlewareFn[I, O]) {
	buildLambda[I, O](fn, middlewares...).StartBatch()
}

func buildLambda[I, O any](fn HandlerFn[I, O], middlewares ...MiddlewareFn[I, O]) *Builder[I, O] {
	mw := append([]MiddlewareFn[I, O]{LoggerMiddleware[I, O]}, middlewares...)
	return NewLambdaBuilder[I, O](fn).WithMiddlewares(mw...)
}
