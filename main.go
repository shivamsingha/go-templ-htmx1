package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"

	"github.com/a-h/templ"
	"github.com/go-playground/validator/v10"
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

var validate *validator.Validate

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

		var inp Input

		if err := c.BodyParser(&inp); err != nil {
			return c.Status(400).SendString(err.Error())
		}

		u, err := q.GetUser(context.Background(), inp.Email)

		if err != nil {
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
			Name            string `form:"name" validate:"required,min=3"`
			Email           string `form:"email" validate:"required,email"`
			Password        string `form:"password" validate:"required,min=8,eqfield=ConfirmPassword,customPassword"`
			ConfirmPassword string `form:"confirm_password"`
		}

		var inp Input

		if err := c.BodyParser(&inp); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		validate = validator.New(validator.WithRequiredStructEnabled())
		validate.RegisterValidation("customPassword", validatePassword)
		err := validate.Struct(inp)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
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

func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if matched, _ := regexp.MatchString(`[A-Z]`, password); !matched {
		return false
	}

	if matched, _ := regexp.MatchString(`[a-z]`, password); !matched {
		return false
	}

	if matched, _ := regexp.MatchString(`[0-9]`, password); !matched {
		return false
	}

	if matched, _ := regexp.MatchString(`[!@#$%^&*()_+{}\[\]:;<>,.?~\\/\-=|"']`, password); !matched {
		return false
	}

	return true
}
