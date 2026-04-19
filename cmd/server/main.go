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
	app.Post("/auth/login", middleware.LoginRateLimit(), h.Login)
	app.Post("/auth/logout", h.Logout)

	// JSON API routes (API key auth — must be before protected group)
	api := app.Group("/api/v1", middleware.APIKeyAuth(cfg.APIKey))
	api.Get("/dashboard", h.APIDashboard)
	api.Get("/companies", h.APIListCompanies)
	api.Get("/tasks", h.APIListTasks)
	api.Post("/inbox", h.APICreateInbox)
	api.Post("/bookmarks", h.APICreateBookmark)
	api.Post("/work-logs", h.APICreateWorkLog)
	api.Post("/knowledge", h.APICreateKnowledge)
	// Integration endpoints
	api.Get("/notion/status", h.NotionSyncStatus)
	api.Post("/notion/sync", h.NotionSyncTrigger)
	api.Post("/telegram/webhook", h.TelegramWebhook)

	// Protected routes (require authentication)
	app.Use(middleware.AuthRequired(queries, cfg))

	// Auth info + Profile
	app.Get("/auth/me", h.AuthMe)
	app.Get("/profile", h.ProfilePageHandler)
	app.Post("/profile/password", h.ChangePassword)

	// Dashboard notes
	app.Post("/dashboard/notes", h.SaveDashboardNotes)

	// Comments
	app.Get("/comments", h.LoadComments)
	app.Post("/comments", h.CreateComment)
	app.Post("/comments/:id/delete", h.DeleteComment)

	// Notifications
	app.Get("/notifications", h.Notifications)
	app.Get("/notifications/count", h.NotificationCount)
	app.Post("/notifications/:id/read", h.MarkNotificationRead)
	app.Post("/notifications/read-all", h.MarkAllRead)

	// Search
	app.Get("/search", h.Search)

	// Admin-only routes (owner + core_staff)
	admin := app.Group("", middleware.RequireRole("owner", "core_staff"))
	admin.Get("/users", h.Users)
	admin.Get("/admin/settings", h.AdminSettings)
	admin.Post("/admin/settings", h.SaveSettings)
	admin.Get("/admin/work-types", h.AdminWorkTypes)
	admin.Post("/admin/work-types", h.CreateWorkType)
	admin.Post("/admin/work-types/:id", h.UpdateWorkType)
	admin.Post("/admin/work-types/:id/delete", h.DeleteWorkType)
	admin.Post("/users", h.CreateUser)

	// Pages — Inbox is home
	app.Get("/", func(c *fiber.Ctx) error { return c.Redirect("/inbox") })
	app.Get("/dashboard", h.Dashboard)
	app.Get("/inbox", h.Inbox)
	app.Get("/work-logs", h.WorkLogs)
	app.Get("/tasks", h.Tasks)
	app.Get("/content", h.Content)
	app.Get("/companies", h.Companies)
	app.Get("/campaigns", h.Campaigns)
	app.Get("/knowledge", h.Knowledge)
	app.Get("/bookmarks", h.Bookmarks)
	app.Get("/prompts", h.Prompts)

	// Create pages
	app.Get("/inbox/new", h.InboxCreate)

	// Detail pages
	app.Get("/content/:id", h.ContentDetail)
	app.Get("/companies/:id", h.CompanyDetail)
	app.Get("/tasks/:id", h.TaskDetail)
	app.Get("/work-logs/:id", h.WorkLogDetail)
	app.Get("/campaigns/:id", h.CampaignDetail)
	app.Get("/knowledge/:id", h.KnowledgeDetail)
	app.Get("/inbox/:id", h.InboxDetail)
	app.Get("/bookmarks/:id", h.BookmarkDetail)

	// CRUD POST/PUT/DELETE routes
	app.Post("/companies", h.CreateCompany)
	app.Post("/companies/:id", h.UpdateCompanyForm)
	app.Post("/companies/:id/assign", h.AssignUserToCompany)
	app.Post("/assignments/:id/delete", h.RemoveAssignment)
	app.Post("/inbox", h.CreateInboxItem)
	app.Post("/inbox/batch-triage", h.BatchTriageInbox)
	app.Post("/inbox/:id/triage", h.TriageInbox)
	app.Post("/inbox/:id/archive", h.ArchiveInbox)
	app.Post("/tasks", h.CreateTask)
	app.Post("/tasks/:id", h.UpdateTaskForm)
	app.Post("/tasks/:id/status", h.UpdateTaskStatusForm)
	app.Post("/tasks/:id/delete", h.DeleteTask)
	app.Post("/content", h.CreateContent)
	app.Post("/content/:id", h.UpdateContentForm)
	app.Post("/content/:id/review", h.ReviewContentForm)
	app.Post("/content/:id/publish", h.PublishContent)
	app.Post("/content/:id/delete", h.DeleteContent)
	app.Get("/work-logs/chart", h.WorkLogChart)
	app.Post("/work-logs", h.CreateWorkLog)
	app.Post("/work-logs/batch-approve", h.BatchApproveWorkLogs)
	app.Post("/work-logs/:id/approve", h.ApproveWorkLogForm)
	app.Post("/work-logs/:id/reject", h.RejectWorkLogForm)
	app.Post("/campaigns", h.CreateCampaign)
	app.Post("/campaigns/:id", h.UpdateCampaignForm)
	app.Post("/campaigns/:id/delete", h.DeleteCampaign)
	app.Post("/knowledge", h.CreateKnowledgeItem)
	app.Post("/knowledge/:id", h.UpdateKnowledgeForm)
	app.Post("/knowledge/:id/delete", h.DeleteKnowledge)
	// File upload
	app.Post("/upload", h.Upload)

	app.Post("/bookmarks", h.CreateBookmark)
	app.Post("/bookmarks/:id", h.UpdateBookmarkForm)
	app.Post("/bookmarks/:id/delete", h.DeleteBookmark)

	log.Printf("Bình Vương OS starting on :%s", cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
}
