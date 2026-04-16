package handler

import (
	"database/sql"

	"binhvuongos/internal/db/generated"
	"binhvuongos/internal/middleware"
	"binhvuongos/web/templates/pages"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Content(c *fiber.Ctx) error {
	items, err := h.queries.ListContent(c.Context(), 50, 0)
	if err != nil {
		return render(c, pages.ContentListPage(pages.ContentPageData{}))
	}
	total, _ := h.queries.CountContent(c.Context())
	statusCounts, _ := h.queries.CountContentByStatus(c.Context())

	counts := make(map[string]int64)
	for _, sc := range statusCounts {
		counts[sc.Status] = sc.Count
	}

	companies, _ := h.queries.ListCompanies(c.Context(), 50, 0)

	data := pages.ContentPageData{
		Items:        toTemplContents(items),
		Total:        total,
		StatusCounts: counts,
		Companies:    toTemplCompanies(companies),
	}
	return render(c, pages.ContentListPage(data))
}

func (h *Handler) CreateContent(c *fiber.Ctx) error {
	user := GetUser(c)
	title := c.FormValue("title")
	contentType := c.FormValue("content_type")
	companyID := c.FormValue("company_id")

	if title == "" || contentType == "" || companyID == "" {
		return c.Redirect("/content")
	}

	_, _ = h.queries.CreateContent(c.Context(), generated.CreateContentParams{
		Title:       title,
		ContentType: contentType,
		CompanyID:   middleware.StringToUUID(companyID),
		AuthorID:    user.ID,
		Status:      "idea",
		CreatedBy:   user.ID,
	})
	return c.Redirect("/content")
}

func (h *Handler) ReviewContentForm(c *fiber.Ctx) error {
	id := c.Params("id")
	status := c.FormValue("status")
	reviewNotes := c.FormValue("review_notes")
	if status == "" {
		return c.Redirect("/content/" + id)
	}
	_, _ = h.queries.ReviewContent(c.Context(), middleware.StringToUUID(id), status,
		sql.NullString{String: reviewNotes, Valid: reviewNotes != ""})
	return c.Redirect("/content/" + id)
}

func toTemplContents(items []generated.Content) []pages.ContentItem {
	result := make([]pages.ContentItem, len(items))
	for i, c := range items {
		result[i] = pages.ContentItem{
			ID:          middleware.UUIDToString(c.ID),
			Title:       c.Title,
			ContentType: c.ContentType,
			Status:      c.Status,
			PublishDate: formatDate(c.PublishDate),
		}
	}
	return result
}
