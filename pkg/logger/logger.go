package logger

import (
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

type key int

const (
	_key key = iota
)

func NewContextWithLogger(ctx context.Context) (context.Context, *zap.Logger) {
	encoderCfg := zapcore.EncoderConfig{
		MessageKey:  "msg",
		LevelKey:    "level",
		TimeKey:     "ts",
		EncodeLevel: zapcore.LowercaseLevelEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format(time.RFC3339))
		},
		EncodeDuration: zapcore.StringDurationEncoder,
	}
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.AddSync(os.Stdout),
		zapcore.DebugLevel,
	)

	log := zap.New(core)
	ctx = context.WithValue(ctx, _key, log)

	return ctx, log
}

func New(ctx context.Context) *zap.Logger {
	log, _ := ctx.Value(_key).(*zap.Logger)
	return log
}
