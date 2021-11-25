package logging

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func ContextParser(c *fiber.Ctx) []zapcore.Field {
	return []zapcore.Field{
		zap.String("method", c.Method()),
		zap.String("localAddr", c.Context().LocalAddr().String()),
		zap.String("remoteAddr", c.Context().RemoteAddr().String()),
		zap.String("uri", c.Request().URI().String()),
	}
}
