package handler

import (
	"binhvuongos/internal/db/generated"
	"binhvuongos/internal/middleware"
	"binhvuongos/web/templates/pages"

	"github.com/gofiber/fiber/v2"
)

// AssignUserToCompany handles POST /companies/:id/assign
func (h *Handler) AssignUserToCompany(c *fiber.Ctx) error {
	companyID := c.Params("id")
	userID := c.FormValue("user_id")
	roleInCompany := c.FormValue("role_in_company")

	if userID == "" {
		return c.Redirect("/companies/" + companyID)
	}

	_, _ = h.queries.CreateAssignment(c.Context(), generated.CreateAssignmentParams{
		UserID:        middleware.StringToUUID(userID),
		CompanyID:     middleware.StringToUUID(companyID),
		RoleInCompany: roleInCompany,
		CanView:       true,
		CanEdit:       true,
		CanApprove:    false,
	})
	return c.Redirect("/companies/" + companyID)
}

// RemoveAssignment handles POST /assignments/:id/delete
func (h *Handler) RemoveAssignment(c *fiber.Ctx) error {
	id := c.Params("id")
	companyID := c.FormValue("company_id")
	_ = h.queries.DeleteAssignment(c.Context(), middleware.StringToUUID(id))
	if companyID != "" {
		return c.Redirect("/companies/" + companyID)
	}
	return c.Redirect("/companies")
}

func toTemplAssignments(assignments []generated.AssignmentWithUser) []pages.AssignmentItem {
	items := make([]pages.AssignmentItem, len(assignments))
	for i, a := range assignments {
		items[i] = pages.AssignmentItem{
			ID:            middleware.UUIDToString(a.ID),
			UserID:        middleware.UUIDToString(a.UserID),
			FullName:      a.FullName,
			Email:         a.Email,
			RoleInCompany: nullStr(a.RoleInCompany),
			UserRole:      a.UserRole,
		}
	}
	return items
}
