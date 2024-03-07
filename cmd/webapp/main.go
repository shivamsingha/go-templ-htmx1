package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"

	database "example/hello/internal/database"
	"example/hello/internal/middleware"
	"example/hello/internal/service"
	"example/hello/internal/util"
	view "example/hello/web/templates"
	"example/hello/web/templates/layout"
	"example/hello/web/templates/partial"
)

func init() {
	database.CreatePool()
	database.ConnectDB()
	middleware.CreateSessionStore()
}

func main() {
	app := fiber.New()
	app.Use(helmet.New())
	app.Use(logger.New())
	app.Static("/static", "./web/static")

	app.Get("/", util.Render(layout.Base(view.Index())))

	app.Route("/login", func(router fiber.Router) {
		router.Get("/", util.Render(layout.Base(partial.Login())))
		router.Post("/", service.LoginHandler)
	})

	app.Route("/signup", func(router fiber.Router) {
		router.Get("/", util.Render(layout.Base(partial.Signup())))
		router.Post("/", service.SignupHandler)
	})

	app.Route("/forgot-password", func(router fiber.Router) {
		router.Get("/", util.Render(layout.Base(partial.ForgotPassword())))
		router.Post("/", service.ForgotPassword)
	})

	app.Route("/reset-password", func(router fiber.Router) {
		router.Get("/", util.Render(layout.Base(partial.ResetPassword())))
		// router.Post("/")
	})

	app.Use(middleware.NotFoundMiddleware())

	log.Fatal(
		app.Listen(":3000"))
}
