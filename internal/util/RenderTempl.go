package util

import (
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

func Render(component templ.Component, options ...func(*templ.ComponentHandler)) fiber.Handler {
	return func(c *fiber.Ctx) error {
		componentHandler := templ.Handler(component)
		for _, o := range options {
			o(componentHandler)
		}
		return adaptor.HTTPHandler(componentHandler)(c)
	}
}
