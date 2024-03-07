package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/postgres/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"example/hello/internal/middleware"
	database "example/hello/internal/model"
	"example/hello/internal/service"
	"example/hello/internal/util"
	view "example/hello/web/templates"
	"example/hello/web/templates/layout"
	"example/hello/web/templates/partial"
)

func main() {
	dbConfig, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))

	pool, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	conn, err := pool.Acquire(context.Background())
	if err != nil {
		log.Fatal("Error while acquiring connection from the database pool!!")
	}

	defer conn.Release()

	q := database.New(conn)

	storage := postgres.New(postgres.Config{
		DB:         pool,
		Table:      "fiber_storage",
		Reset:      false,
		GCInterval: 10 * time.Second,
	})

	store := session.New(session.Config{
		Storage: storage,
	})

	app := fiber.New()
	app.Use(helmet.New())
	app.Use(logger.New())
	app.Static("/static", "./web/static")

	app.Get("/", util.Render(layout.Base(view.Index())))

	app.Route("/login", func(router fiber.Router) {
		router.Get("/", util.Render(layout.Base(partial.Login())))
		router.Post("/", service.LoginHandler(q))
	})

	app.Route("/signup", func(router fiber.Router) {
		router.Get("/", util.Render(layout.Base(partial.Signup())))
		router.Post("/", service.SignupHandler(q))
	})

	app.Route("/forgot-password", func(router fiber.Router) {
		router.Get("/", util.Render(layout.Base(partial.ForgotPassword())))
		router.Post("/", service.ForgotPassword(q))
	})

	app.Route("/reset-password", func(router fiber.Router) {
		router.Get("/", util.Render(layout.Base(partial.ResetPassword())))
		// router.Post("/")
	})

	app.Use(middleware.NotFoundMiddleware())

	log.Fatal(
		app.Listen(":3000"))
}
