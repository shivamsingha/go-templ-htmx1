package middleware

import (
	"example/hello/internal/database"
	"time"

	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/postgres/v2"
)

var Store *session.Store

func CreateSessionStore() {
	storage := postgres.New(postgres.Config{
		DB:         database.Pool,
		Table:      "fiber_storage",
		Reset:      false,
		GCInterval: 10 * time.Second,
	})

	store := session.New(session.Config{
		Storage: storage,
	})

	Store = store
}
