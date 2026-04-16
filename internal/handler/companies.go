package handler

import (
	"binhvuongos/web/templates/pages"
	"github.com/gofiber/fiber/v2"
)

func Companies(c *fiber.Ctx) error {
	return render(c, pages.CompaniesPage())
}
