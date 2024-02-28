package main

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"example/hello/view"
	"example/hello/view/layout"
	"example/hello/view/partial"
)

func main() {
	app := fiber.New()
	app.Use(helmet.New())
	app.Use(logger.New())
	app.Static("/public", "./public")

	app.Get("/", func(c *fiber.Ctx) error {
		return Render(c, layout.Base(view.Index()))
	})

	app.Get("/login", func(c *fiber.Ctx) error {
		return Render(c, layout.Base(partial.Login()))
	})


	app.Get("/signup", func(c *fiber.Ctx) error {
		return Render(c, layout.Base(partial.Signup()))
	})

	app.Use(NotFoundMiddleware)

	log.Fatal(
		app.Listen(":3000"))
}

func NotFoundMiddleware(c *fiber.Ctx) error {
	return Render(c, layout.Base(partial.Error404()), templ.WithStatus(http.StatusNotFound))
}

func Render(c *fiber.Ctx, component templ.Component, options ...func(*templ.ComponentHandler)) error {
	componentHandler := templ.Handler(component)
	for _, o := range options {
		o(componentHandler)
	}
	return adaptor.HTTPHandler(componentHandler)(c)
}
