package middleware

import (
	"example/hello/internal/util"
	"example/hello/web/templates/layout"
	"example/hello/web/templates/partial"
	"net/http"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
)

func NotFoundMiddleware(c *fiber.Ctx) error {
	return util.Render(c, layout.Base(partial.Error404()), templ.WithStatus(http.StatusNotFound))
}
