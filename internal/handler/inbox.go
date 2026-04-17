package handler

import (
	"database/sql"

	"binhvuongos/internal/db/generated"
	"binhvuongos/internal/middleware"
	"binhvuongos/web/templates/pages"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Inbox(c *fiber.Ctx) error {
	page, limit, offset := getPage(c)
	items, err := h.queries.ListInboxItems(c.Context(), limit, offset)
	if err != nil {
		return render(c, pages.InboxListPage(pages.InboxPageData{}))
	}
	rawCount, _ := h.queries.CountInboxItemsByStatus(c.Context(), "raw")
	total, _ := h.queries.CountInboxItems(c.Context())
	companies, _ := h.queries.ListCompanies(c.Context(), 50, 0)

	data := pages.InboxPageData{
		Items:      toTemplInboxItems(items),
		RawCount:   rawCount,
		Total:      total,
		Companies:  toTemplCompanies(companies),
		Page:       page,
		TotalPages: totalPages(total),
	}
	return render(c, pages.InboxListPage(data))
}

func (h *Handler) CreateInboxItem(c *fiber.Ctx) error {
	user := GetUser(c)
	content := c.FormValue("content")
	url := c.FormValue("url")
	source := c.FormValue("source")
	itemType := c.FormValue("item_type")

	if content == "" {
		return c.Redirect("/inbox")
	}
	if source == "" {
		source = "manual"
	}

	_, _ = h.queries.CreateInboxItem(c.Context(), generated.CreateInboxItemParams{
		Content:     content,
		URL:         sql.NullString{String: url, Valid: url != ""},
		Source:      sql.NullString{String: source, Valid: true},
		ItemType:    sql.NullString{String: itemType, Valid: itemType != ""},
		SubmittedBy: user.ID,
	})
	return c.Redirect("/inbox")
}

func (h *Handler) TriageInbox(c *fiber.Ctx) error {
	id := c.Params("id")
	destination := c.FormValue("destination")
	triageNotes := c.FormValue("triage_notes")

	_, _ = h.queries.TriageInboxItem(c.Context(), generated.TriageInboxItemParams{
		ID:          middleware.StringToUUID(id),
		Destination: sql.NullString{String: destination, Valid: destination != ""},
		TriageNotes: sql.NullString{String: triageNotes, Valid: triageNotes != ""},
	})
	return c.Redirect("/inbox")
}

func toTemplInboxItems(items []generated.InboxItem) []pages.InboxItemData {
	result := make([]pages.InboxItemData, len(items))
	for i, item := range items {
		result[i] = pages.InboxItemData{
			ID:        middleware.UUIDToString(item.ID),
			Content:   item.Content,
			URL:       nullStr(item.URL),
			Source:    nullStr(item.Source),
			ItemType:  nullStr(item.ItemType),
			Status:    item.Status,
			CreatedAt: formatTime(item.CreatedAt),
			TimeAgo:   timeAgo(item.CreatedAt),
		}
	}
	return result
}
