package main

import (
	"log"
	"os"

	"binhvuongos/internal/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New(fiber.Config{
		AppName: "Bình Vương OS v1.0",
	})

	app.Use(logger.New())

	// Static files
	app.Static("/static", "./web/static")

	// Pages
	app.Get("/", handler.Dashboard)
	app.Get("/work-logs", handler.WorkLogs)
	app.Get("/inbox", handler.Inbox)
	app.Get("/tasks", handler.Tasks)
	app.Get("/content", handler.Content)
	app.Get("/companies", handler.Companies)
	app.Get("/campaigns", handler.Campaigns)
	app.Get("/knowledge", handler.Knowledge)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Bình Vương OS starting on :%s", port)
	log.Fatal(app.Listen(":" + port))
}
