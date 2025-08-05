package routes

import (
	"github.com/aukawut/BackendCompoundHt/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupPartRoutes(app *fiber.App) {
	part := app.Group("/api/part")
	part.Get("/compound", handlers.GetPartCodeByCompoundTags)
}
