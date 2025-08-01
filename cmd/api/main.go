package main

import (
	"fmt"
	"log"
	"my-go-api/internal/config"
	"my-go-api/internal/routes"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	// "github.com/joho/godotenv"
)

func main() {
	// if err := godotenv.Load(); err != nil {
	// 	log.Fatal("❌ Error loading .env file")
	// }

	db := config.InitDB()
	config.InitCloudinary()
	app := fiber.New()

	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = "*"
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins: allowedOrigins,
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders: "Origin, Content-Type, Authorization",
	}))

	routes.SetupRoutes(app, db)

	port := os.Getenv("PORT")
	fmt.Println("✅ Server running on http://localhost:" + port)
	log.Fatal(app.Listen(":" + port))
}
