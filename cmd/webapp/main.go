package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jackc/pgx/v5"

	"example/hello/internal/middleware"
	database "example/hello/internal/model"
	"example/hello/internal/service"
	"example/hello/internal/util"
	view "example/hello/web/templates"
	"example/hello/web/templates/layout"
	"example/hello/web/templates/partial"
)

func main() {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())
	q := database.New(conn)

	app := fiber.New()
	app.Use(helmet.New())
	app.Use(logger.New())
	app.Static("/static", "./web/static")

	app.Get("/", util.Render(layout.Base(view.Index())))

	app.Route("/login", func(login fiber.Router) {
		login.Get("/", util.Render(layout.Base(partial.Login())))
		login.Post("/", service.LoginHandler(q))
	})

	app.Route("/signup", func(signup fiber.Router) {
		signup.Get("/", util.Render(layout.Base(partial.Signup())))
		signup.Post("/", service.SignupHandler(q))
	})

	app.Route("/forgot-password", func(forgot fiber.Router) {
		forgot.Get("/", util.Render(layout.Base(partial.ForgotPassword())))
		forgot.Post("/", service.ForgotPassword(q))
	})

	app.Use(middleware.NotFoundMiddleware())

	log.Fatal(
		app.Listen(":3000"))
}
