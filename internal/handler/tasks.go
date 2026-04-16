package handler

import (
	"binhvuongos/web/templates/pages"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Tasks(c *fiber.Ctx) error {
	return render(c, pages.TasksPage())
}

func (h *Handler) CreateTask(c *fiber.Ctx) error {
	// TODO: Phase 5 — parse form and create task
	return c.Redirect("/tasks")
}

func (h *Handler) UpdateTaskStatusForm(c *fiber.Ctx) error {
	// TODO: Phase 5 — update task status (HTMX)
	return c.Redirect("/tasks")
}
