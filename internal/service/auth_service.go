package service

import (
	"context"
	"errors"
	"log"
	"time"

	"example/hello/internal/database"
	"example/hello/internal/database/sqlc"
	"example/hello/internal/middleware"
	"example/hello/internal/util"
	"example/hello/web/templates/components"
	"example/hello/web/templates/layout"
	"example/hello/web/templates/partial"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/matthewhartstonge/argon2"
)

func LoginHandler(c *fiber.Ctx) error {
	type Input struct {
		Email    string `form:"email"`
		Password string `form:"password"`
	}

	var inp Input

	if err := c.BodyParser(&inp); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
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
	time.Sleep(1 * time.Second)
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
		c.Status(fiber.StatusBadRequest)
		if c.Get("HX-Request") == "true" {
			c.Set("HX-Select", "#SignupError")
			return util.Render(c, components.SignupError(err.Error()))
		}
		return util.Render(c, layout.Base(partial.Signup(err.Error())))
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
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	task := &ForgotPasswordTask{Email: inp.Email}

	go func() {
		if err := task.Process(); err != nil {
			log.Printf("Error creating reset password link: %v", err)
		}
	}()

	return c.Status(fiber.StatusOK).SendString("Password Reset Link sent to email")

}

func ResetPasswordPage(c *fiber.Ctx) error {
	token := c.Query("token")

	if token == "" {
		return c.Status(400).SendString("Token required")
	}

	q := sqlc.New(database.Conn)
	reset, err := q.PopResetToken(context.Background(), token)

	if err != nil {
		log.Printf("Error popping reset token at ResetPasswordPage: %v", err)
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusUnauthorized).SendString("Invalid Token")
		}
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	if time.Now().After(reset.ExpiresAt.Time) {
		return c.Status(fiber.StatusForbidden).SendString("Token expired")
	}

	sess, err := middleware.Store.Get(c)
	if err != nil {
		log.Printf("Error getting session at ResetPasswordPage: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	sess.Set("resetuser", reset.UserID)
	sess.SetExpiry(10 * time.Minute)
	if err := sess.Save(); err != nil {
		log.Printf("Error saving session at ResetPasswordPage: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	return util.Render(c, layout.Base(partial.ResetPassword()))
}

func ResetPassword(c *fiber.Ctx) error {
	sess, err := middleware.Store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	resetuser, ok := sess.Get("resetuser").(pgtype.UUID)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	if err := sess.Destroy(); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	type Input struct {
		Password        string `form:"password" validate:"required,min=8,eqfield=ConfirmPassword,customPassword"`
		ConfirmPassword string `form:"confirm_password"`
	}

	var inp Input

	if err := c.BodyParser(&inp); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterValidation("customPassword", util.ValidatePassword)
	if err := validate.Struct(inp); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	argon := argon2.RecommendedDefaults()

	encoded, err := argon.HashEncoded([]byte(inp.Password))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	u := sqlc.UpdatePasswordUserParams{
		ID:       resetuser,
		Password: string(encoded),
	}

	q := sqlc.New(database.Conn)
	if err := q.UpdatePasswordUser(context.Background(), u); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	return c.Status(fiber.StatusCreated).SendString("Created")
}
