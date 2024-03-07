package middleware

import (
	"example/hello/internal/util"
	"example/hello/web/templates/layout"
	"example/hello/web/templates/partial"
	"net/http"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
)

func NotFoundMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Assuming layout.Base(partial.Error404()) returns a templ.Component
		// and templ.WithStatus(http.StatusNotFound) is a valid option function
		handler := util.Render(layout.Base(partial.Error404()), templ.WithStatus(http.StatusNotFound))
		return handler(c)
	}
}
