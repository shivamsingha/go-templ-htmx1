package service

import (
	"context"
	"errors"
	"log"

	"example/hello/internal/database"
	"example/hello/internal/database/sqlc"
	"example/hello/internal/util"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/matthewhartstonge/argon2"
)

func LoginHandler(c *fiber.Ctx) error {
	type Input struct {
		Email    string `form:"email"`
		Password string `form:"password"`
	}

	var inp Input

	if err := c.BodyParser(&inp); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	q := sqlc.New(database.Conn)
	u, err := q.GetUser(context.Background(), inp.Email)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusUnauthorized).SendString("User not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	passwordMatches, err := argon2.VerifyEncoded([]byte(inp.Password), []byte(u.Password))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	if !passwordMatches {
		return c.Status(fiber.StatusUnauthorized).SendString("Password wrong")
	}

	return c.Status(fiber.StatusOK).SendString("Logged In")
}

func SignupHandler(c *fiber.Ctx) error {
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

	argon := argon2.RecommendedDefaults()

	encoded, err := argon.HashEncoded([]byte(inp.Password))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	u := sqlc.CreateUserParams{
		Name:     inp.Name,
		Email:    inp.Email,
		Password: string(encoded),
	}

	q := sqlc.New(database.Conn)
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

type ForgotPasswordTask struct {
	Email string
}

func (bt *ForgotPasswordTask) Process() error {
	q := sqlc.New(database.Conn)
	u, err := q.GetUser(context.Background(), bt.Email)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("Error while generating reset link, " + bt.Email + " doesn't exist. ")
		}
		return err
	}

	count, err := q.CountResetTokensByUser(context.Background(), u.ID)
	if err != nil {
		return err
	}

	if count >= 10 {
		return errors.New("Error while generating reset link, " + bt.Email + " reset limit exceeded.")
	}

	token, err := util.GenerateResetToken()
	if err != nil {
		return err
	}

	ResetTokenRow := sqlc.CreateResetTokenParams{
		UserID: u.ID,
		Token:  token,
	}

	if err := q.CreateResetToken(context.Background(), ResetTokenRow); err != nil {
		return err
	}

	return nil
}

func ForgotPassword(c *fiber.Ctx) error {
	type Input struct {
		Email string `form:"email"`
	}

	var inp Input

	if err := c.BodyParser(&inp); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	task := &ForgotPasswordTask{Email: inp.Email}

	go func() {
		if err := task.Process(); err != nil {
			log.Printf("Error creating reset password link: %v", err)
		}
	}()

	return c.Status(fiber.StatusOK).SendString("Password Reset Link sent to email")

}
