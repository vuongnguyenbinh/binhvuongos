package handler

import (
	"database/sql"

	"binhvuongos/internal/db/generated"
	"binhvuongos/internal/middleware"
	"binhvuongos/web/templates/pages"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Content(c *fiber.Ctx) error {
	page, limit, offset := getPage(c)
	items, err := h.queries.ListContent(c.Context(), limit, offset)
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
		Page:         page,
		TotalPages:   totalPages(total),
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

func (h *Handler) UpdateContentForm(c *fiber.Ctx) error {
	id := c.Params("id")
	title := c.FormValue("title")
	contentType := c.FormValue("content_type")
	status := c.FormValue("status")
	notes := c.FormValue("notes")

	if title == "" {
		return c.Redirect("/content/" + id)
	}

	_, _ = h.queries.UpdateContent(c.Context(), generated.UpdateContentParams{
		ID:          middleware.StringToUUID(id),
		Title:       title,
		ContentType: contentType,
		Status:      status,
		Notes:       sql.NullString{String: notes, Valid: notes != ""},
	})
	return c.Redirect("/content/" + id)
}

func (h *Handler) DeleteContent(c *fiber.Ctx) error {
	id := c.Params("id")
	_ = h.queries.SoftDeleteContent(c.Context(), middleware.StringToUUID(id))
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
