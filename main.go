package main

import (
	"fmt"
	"os"

	"github.com/aukawut/BackendCompoundHt/routes"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {

	// Load environments
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	// API Config
	app := fiber.New(fiber.Config{
		BodyLimit: 100 * 1024 * 1024, //100 MB
	})

	// Allow Cors
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	// Setup API Routes
	routes.SetupPartRoutes(app)
	routes.SetupTagsRoutes(app)

	// Load environment API Port
	PORT := os.Getenv("PORT")

	// Run Services
	app.Listen(":" + PORT)
}
