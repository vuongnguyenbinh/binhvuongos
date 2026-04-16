package handler

import (
	"binhvuongos/web/templates/pages"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) WorkLogs(c *fiber.Ctx) error {
	return render(c, pages.WorkLogsPage())
}

func (h *Handler) CreateWorkLog(c *fiber.Ctx) error {
	// TODO: Phase 5 — parse form and create work log
	return c.Redirect("/work-logs")
}

func (h *Handler) ApproveWorkLogForm(c *fiber.Ctx) error {
	// TODO: Phase 5 — approve work log
	return c.Redirect("/work-logs/" + c.Params("id"))
}

func (h *Handler) RejectWorkLogForm(c *fiber.Ctx) error {
	// TODO: Phase 5 — reject work log
	return c.Redirect("/work-logs/" + c.Params("id"))
}
