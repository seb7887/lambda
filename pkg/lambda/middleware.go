package lambda

import (
	"context"
	"go.uber.org/zap"
	"seb7887/lambda/pkg/logger"
)

type MiddlewareFn[I, O any] func(next HandlerFn[I, O]) HandlerFn[I, O]

func LoggerMiddleware[I, O any](next HandlerFn[I, O]) HandlerFn[I, O] {
	return func(ctx context.Context, in *I) (*O, error) {
		ctx, log := logger.NewContextWithLogger(ctx)

		log.With(zap.Any("event", in)).Info("starting execution...")

		out, err := next(ctx, in)
		if err != nil {
			log.With(zap.Any("error", err)).Error("execution has an error")
		} else {
			log.With(zap.Any("response", out)).Info("success")
		}

		return out, err
	}
}
