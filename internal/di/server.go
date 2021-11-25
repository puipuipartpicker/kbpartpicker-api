package di

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/puipuipartpicker/kbpartpicker/api/pkg/di"
	appErr "github.com/puipuipartpicker/kbpartpicker/api/pkg/error"
	"go.uber.org/zap"
)

// Server object
type Server struct {
	server *fiber.App
}

func newService() *Server {
	return &Server{
		server: fiber.New(),
	}
}

func (s *Server) start() {
	s.setupRoutes()
}

func (s *Server) healthCheck(ctx *fiber.Ctx) error {
	return ctx.SendStatus(http.StatusOK)
}

func (s *Server) handleError(c *fiber.Ctx) (err error) {
	err = c.Next()
	logger := di.GetContextLogger()

	// managed error
	var managed *appErr.Error
	if errors.As(err, &managed) {
		c.Status(http.StatusBadRequest)
		return c.JSON(managed)
	} else if err != nil {
		logger.Error(c, "Received unmanaged error", zap.Error(err))
	}

	return err
}

// func (s *server) installBot(r fiber.Router, sc model.Bot) {
// 	r.Get("/calendars/"+sc.Type()+"/:option", func(ctx *fiber.Ctx) error {
// 		q := string(ctx.Request().URI().QueryString())

// 		c, err := sc.GetOptions(ctx.UserContext(), q, sc.GetRouteParams(ctx))
// 		if err != nil {
// 			return err
// 		}

// 		return ctx.JSON(c)
// 	})

// 	r.Get("/calendars/"+sc.Type(), func(ctx *fiber.Ctx) error {
// 		q := string(ctx.Request().URI().QueryString())

// 		c, err := sc.GetCalendar(ctx.UserContext(), q, sc.GetRouteParams(ctx))
// 		if err != nil {
// 			return err
// 		}

// 		return ctx.JSON(c)
// 	})
// }
