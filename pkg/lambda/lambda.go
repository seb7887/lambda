package lambda

func Start[I, O any](fn HandlerFn[I, O], middlewares ...MiddlewareFn[I, O]) {
	l := NewLambdaBuilder[I, O](fn).
		WithMiddlewares(LoggerMiddleware[I, O]).
		WithMiddlewares(middlewares...)

	l.Start()
}
