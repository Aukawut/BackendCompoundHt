package routes

import (
	"github.com/aukawut/BackendCompoundHt/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupTagsRoutes(app *fiber.App) {
	part := app.Group("/api/tags")

	part.Get("/check/duplicate", handlers.CheckDuplicatedTag)
	part.Get("/list", handlers.GetTagsByDate)
	part.Get("/compound", handlers.GetCompoundTagDetail)
	part.Get("/generate/qrcode", handlers.GenerateQRCode)
	part.Put("/cancel", handlers.CancelTags)
	part.Post("/save", handlers.SaveTags)
}
