package main

import (
	"context"
	"log"

	"binhvuongos/internal/config"
	"binhvuongos/internal/db"
	"binhvuongos/internal/db/generated"
	"binhvuongos/internal/handler"
	"binhvuongos/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	cfg := config.Load()

	app := fiber.New(fiber.Config{
		AppName: "Bình Vương OS v2.0",
	})

	app.Use(logger.New())

	// Static files
	app.Static("/static", "./web/static")

	// Connect to database
	ctx := context.Background()
	pool, err := db.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()
	log.Println("Connected to database")

	queries := generated.New(pool)
	h := handler.NewHandler(queries, cfg)

	// Public routes
	app.Get("/login", h.LoginPage)
	app.Post("/auth/login", h.Login)
	app.Post("/auth/logout", h.Logout)

	// Protected routes (require authentication)
	protected := app.Group("", middleware.AuthRequired(queries, cfg))

	// Auth info
	protected.Get("/auth/me", h.AuthMe)

	// Pages
	protected.Get("/", h.Dashboard)
	protected.Get("/inbox", h.Inbox)
	protected.Get("/work-logs", h.WorkLogs)
	protected.Get("/tasks", h.Tasks)
	protected.Get("/content", h.Content)
	protected.Get("/companies", h.Companies)
	protected.Get("/campaigns", h.Campaigns)
	protected.Get("/knowledge", h.Knowledge)
	protected.Get("/bookmarks", h.Bookmarks)
	protected.Get("/prompts", h.Prompts)

	// Create pages
	protected.Get("/inbox/new", h.InboxCreate)

	// Detail pages
	protected.Get("/content/:id", h.ContentDetail)
	protected.Get("/companies/:id", h.CompanyDetail)
	protected.Get("/tasks/:id", h.TaskDetail)
	protected.Get("/work-logs/:id", h.WorkLogDetail)
	protected.Get("/campaigns/:id", h.CampaignDetail)
	protected.Get("/knowledge/:id", h.KnowledgeDetail)
	protected.Get("/inbox/:id", h.InboxDetail)
	protected.Get("/bookmarks/:id", h.BookmarkDetail)

	// CRUD POST/PUT/DELETE routes
	protected.Post("/companies", h.CreateCompany)
	protected.Post("/companies/:id", h.UpdateCompanyForm)
	protected.Post("/inbox", h.CreateInboxItem)
	protected.Post("/inbox/:id/triage", h.TriageInbox)
	protected.Post("/tasks", h.CreateTask)
	protected.Post("/tasks/:id/status", h.UpdateTaskStatusForm)
	protected.Post("/content", h.CreateContent)
	protected.Post("/content/:id/review", h.ReviewContentForm)
	protected.Post("/work-logs", h.CreateWorkLog)
	protected.Post("/work-logs/:id/approve", h.ApproveWorkLogForm)
	protected.Post("/work-logs/:id/reject", h.RejectWorkLogForm)
	protected.Post("/campaigns", h.CreateCampaign)
	protected.Post("/knowledge", h.CreateKnowledgeItem)
	protected.Post("/bookmarks", h.CreateBookmark)
	protected.Post("/bookmarks/:id/delete", h.DeleteBookmark)

	// JSON API routes (API key auth)
	api := app.Group("/api/v1", middleware.APIKeyAuth(cfg.APIKey))
	api.Get("/dashboard", h.APIDashboard)
	api.Get("/companies", h.APIListCompanies)
	api.Get("/tasks", h.APIListTasks)
	api.Post("/inbox", h.APICreateInbox)
	api.Post("/bookmarks", h.APICreateBookmark)
	api.Post("/work-logs", h.APICreateWorkLog)
	api.Post("/knowledge", h.APICreateKnowledge)

	log.Printf("Bình Vương OS starting on :%s", cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
}
