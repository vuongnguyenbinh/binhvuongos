package handler

import (
	"database/sql"
	"strings"

	"binhvuongos/internal/db/generated"
	"binhvuongos/internal/middleware"
	"binhvuongos/web/templates/pages"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
)

func (h *Handler) Knowledge(c *fiber.Ctx) error {
	q := c.Query("q")
	category := c.Query("category")
	page, limit, offset := getPage(c)

	var items []generated.KnowledgeItem
	var err error

	if q != "" {
		items, err = h.queries.SearchKnowledgeItems(c.Context(), q, limit, offset)
	} else if category != "" {
		items, err = h.queries.ListKnowledgeItemsByCategory(c.Context(), category, limit, offset)
	} else {
		items, err = h.queries.ListKnowledgeItems(c.Context(), limit, offset)
	}
	if err != nil {
		return render(c, pages.KnowledgeListPage(pages.KnowledgePageData{}))
	}

	total, _ := h.queries.CountKnowledgeItems(c.Context())

	data := pages.KnowledgePageData{
		Items:      toTemplKnowledge(items),
		Total:      total,
		Query:      q,
		Category:   category,
		Page:       page,
		TotalPages: totalPages(total),
	}
	return render(c, pages.KnowledgeListPage(data))
}

func (h *Handler) CreateKnowledgeItem(c *fiber.Ctx) error {
	user := GetUser(c)
	title := c.FormValue("title")
	description := c.FormValue("description")
	category := c.FormValue("category")
	topics := c.FormValue("topics")
	scope := c.FormValue("scope")
	sourceURL := c.FormValue("source_url")

	if title == "" || category == "" {
		return c.Redirect("/knowledge")
	}
	if scope == "" {
		scope = "shared"
	}

	var topicSlice []string
	for _, t := range strings.Split(topics, ",") {
		trimmed := strings.TrimSpace(t)
		if trimmed != "" {
			topicSlice = append(topicSlice, trimmed)
		}
	}

	_, _ = h.queries.CreateKnowledgeItem(c.Context(), generated.CreateKnowledgeItemParams{
		Title:       title,
		Description: sql.NullString{String: description, Valid: description != ""},
		Category:    category,
		Topics:      topicSlice,
		Scope:       scope,
		SourceURL:   sql.NullString{String: sourceURL, Valid: sourceURL != ""},
		CreatedBy:   user.ID,
		QualityRating: pgtype.Int4{},
	})
	return c.Redirect("/knowledge")
}

func (h *Handler) UpdateKnowledgeForm(c *fiber.Ctx) error {
	id := c.Params("id")
	title := c.FormValue("title")
	description := c.FormValue("description")
	category := c.FormValue("category")
	topics := c.FormValue("topics")
	scope := c.FormValue("scope")
	sourceURL := c.FormValue("source_url")

	if title == "" || category == "" {
		return c.Redirect("/knowledge/" + id)
	}

	var topicSlice []string
	for _, t := range strings.Split(topics, ",") {
		trimmed := strings.TrimSpace(t)
		if trimmed != "" {
			topicSlice = append(topicSlice, trimmed)
		}
	}

	_, _ = h.queries.UpdateKnowledgeItem(c.Context(), generated.UpdateKnowledgeItemParams{
		ID:          middleware.StringToUUID(id),
		Title:       title,
		Description: sql.NullString{String: description, Valid: description != ""},
		Category:    category,
		Topics:      topicSlice,
		Scope:       scope,
		SourceURL:   sql.NullString{String: sourceURL, Valid: sourceURL != ""},
	})
	return c.Redirect("/knowledge/" + id)
}

func (h *Handler) DeleteKnowledge(c *fiber.Ctx) error {
	id := c.Params("id")
	_ = h.queries.SoftDeleteKnowledgeItem(c.Context(), middleware.StringToUUID(id))
	return c.Redirect("/knowledge")
}

func toTemplKnowledge(items []generated.KnowledgeItem) []pages.KnowledgeItemData {
	result := make([]pages.KnowledgeItemData, len(items))
	for i, k := range items {
		result[i] = pages.KnowledgeItemData{
			ID:          middleware.UUIDToString(k.ID),
			Title:       k.Title,
			Description: nullStr(k.Description),
			Category:    k.Category,
			Topics:      k.Topics,
			Scope:       k.Scope,
			SourceURL:   nullStr(k.SourceURL),
			CreatedAt:   formatTime(k.CreatedAt),
		}
	}
	return result
}
