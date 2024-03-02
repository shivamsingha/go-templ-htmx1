package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"example/hello/database"
	"example/hello/view"
	"example/hello/view/layout"
	"example/hello/view/partial"
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
	app.Static("/public", "./public")

	app.Get("/", func(c *fiber.Ctx) error {
		return Render(c, layout.Base(view.Index()))
	})

	app.Get("/login", func(c *fiber.Ctx) error {
		return Render(c, layout.Base(partial.Login()))
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		type Input struct {
			Email    string `form:"email"`
			Password string `form:"password"`
		}

		inp := new(Input)

		if err := c.BodyParser(inp); err != nil {
			return c.Status(400).SendString(err.Error())
		}

		u, err := q.GetUser(context.Background(), inp.Email)

		if err != nil {
			fmt.Println(inp)
			fmt.Println(err)
			if errors.Is(err, pgx.ErrNoRows) {
				return c.Status(fiber.StatusUnauthorized).SendString("User not found")
			}
			return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
		}

		if inp.Password != u.Password {
			return c.Status(fiber.StatusUnauthorized).SendString("Password wrong")
		}

		return c.Status(fiber.StatusOK).SendString("Logged In")
	})

	app.Get("/signup", func(c *fiber.Ctx) error {
		return Render(c, layout.Base(partial.Signup()))
	})

	app.Post("/signup", func(c *fiber.Ctx) error {
		type Input struct {
			Name            string `form:"name"`
			Email           string `form:"email"`
			Password        string `form:"password"`
			ConfirmPassword string `form:"confirm_password"`
		}

		inp := new(Input)

		if err := c.BodyParser(inp); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		if inp.Password != inp.ConfirmPassword {
			return c.Status(fiber.StatusBadRequest).SendString("Passwords don't match")
		}

		fmt.Println(inp)
		u := database.CreateUserParams{
			Name:     inp.Name,
			Email:    inp.Email,
			Password: inp.Password,
		}
		fmt.Println(u)

		if err := q.CreateUser(context.Background(), u); err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				if pgErr.Code == pgerrcode.UniqueViolation {
					return c.Status(fiber.StatusConflict).SendString("User with this email already exists")
				} else {
					return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
				}
			}
		}

		return c.Status(fiber.StatusCreated).SendString("Created")
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
