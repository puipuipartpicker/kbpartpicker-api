package logging

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ContextLogger interface {
	Info(ctx *fiber.Ctx, msg string, fields ...zapcore.Field)
	Error(ctx *fiber.Ctx, msg string, fields ...zapcore.Field)
	Fatal(ctx *fiber.Ctx, msg string, fields ...zapcore.Field)
	Warn(ctx *fiber.Ctx, msg string, fields ...zapcore.Field)
	Debug(ctx *fiber.Ctx, msg string, fields ...zapcore.Field)
	Access(ctx *fiber.Ctx) error
	Named(name string) ContextLogger
	With(fields ...zapcore.Field) ContextLogger
	Unwrap() *zap.Logger
}

type contextParser func(ctx *fiber.Ctx) []zap.Field

type contextLogger struct {
	logger        *zap.Logger
	contextParser contextParser
}

func NewContextLogger(logLevel string, contextParser contextParser, opts ...zap.Option) (ContextLogger, error) {
	opts = append(opts, zap.AddCallerSkip(1))

	l, err := newZapLogger(logLevel, opts...)
	if err != nil {
		return nil, err
	}

	return &contextLogger{logger: l, contextParser: contextParser}, nil
}

func NewDevelopmentContextLogger(contextParser contextParser, opts ...zap.Option) (ContextLogger, error) {
	opts = append(opts, zap.AddCallerSkip(1))

	l, err := zap.NewDevelopment(opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to init logger for development: %w", err)
	}

	return &contextLogger{logger: l, contextParser: contextParser}, nil
}

func (l *contextLogger) Info(ctx *fiber.Ctx, msg string, fields ...zapcore.Field) {
	f := append(l.parseContext(ctx), fields...)
	l.logger.Info(msg, f...)
}

func (l *contextLogger) Error(ctx *fiber.Ctx, msg string, fields ...zapcore.Field) {
	f := append(l.parseContext(ctx), fields...)
	l.logger.Error(msg, f...)
}

func (l *contextLogger) Fatal(ctx *fiber.Ctx, msg string, fields ...zapcore.Field) {
	f := append(l.parseContext(ctx), fields...)
	l.logger.Fatal(msg, f...)
}

func (l *contextLogger) Warn(ctx *fiber.Ctx, msg string, fields ...zapcore.Field) {
	f := append(l.parseContext(ctx), fields...)
	l.logger.Warn(msg, f...)
}

func (l *contextLogger) Debug(ctx *fiber.Ctx, msg string, fields ...zapcore.Field) {
	f := append(l.parseContext(ctx), fields...)
	l.logger.Debug(msg, f...)
}

var traceIDCtxKey = struct{}{}

func (l *contextLogger) Access(ctx *fiber.Ctx) error {
	ctx.SetUserContext(context.WithValue(ctx.Context(), traceIDCtxKey, ctx.Get("X-Trace-ID")))
	l.logger.Info("access", l.parseContext(ctx)...)
	return ctx.Next()
}

func (l *contextLogger) Named(name string) ContextLogger {
	return &contextLogger{logger: l.logger.Named(name), contextParser: l.contextParser}
}

func (l *contextLogger) Unwrap() *zap.Logger {
	return l.logger
}

func (l *contextLogger) With(fields ...zapcore.Field) ContextLogger {
	return &contextLogger{logger: l.logger.With(fields...), contextParser: l.contextParser}
}

func (l *contextLogger) parseContext(ctx *fiber.Ctx) []zapcore.Field {
	return l.contextParser(ctx)
}
