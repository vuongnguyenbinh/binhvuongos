package handler

import (
	"binhvuongos/web/templates/pages"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Companies(c *fiber.Ctx) error {
	return render(c, pages.CompaniesPage())
}

func (h *Handler) CreateCompany(c *fiber.Ctx) error {
	// TODO: Phase 4 — parse form and create company
	return c.Redirect("/companies")
}

func (h *Handler) UpdateCompanyForm(c *fiber.Ctx) error {
	// TODO: Phase 4 — parse form and update company
	return c.Redirect("/companies/" + c.Params("id"))
}
