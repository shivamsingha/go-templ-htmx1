package service

import (
	"context"
	"errors"
	"fmt"

	database "example/hello/internal/model"
	"example/hello/internal/util"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func LoginHandler(q *database.Queries) fiber.Handler {
	return func(c *fiber.Ctx) error {
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
	}
}

func SignupHandler(q *database.Queries) fiber.Handler {
	return func(c *fiber.Ctx) error {
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

		validate := validator.New(validator.WithRequiredStructEnabled())
		validate.RegisterValidation("customPassword", util.ValidatePassword)
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
	}
}
