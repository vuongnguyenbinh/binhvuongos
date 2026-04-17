package handler

import (
	"database/sql"
	"strings"

	"binhvuongos/internal/db/generated"
	"binhvuongos/internal/middleware"
	"binhvuongos/web/templates/pages"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Bookmarks(c *fiber.Ctx) error {
	page, limit, offset := getPage(c)
	bookmarks, err := h.queries.ListBookmarks(c.Context(), limit, offset)
	if err != nil {
		return render(c, pages.BookmarksListPage(pages.BookmarksPageData{}))
	}
	total, _ := h.queries.CountBookmarks(c.Context())

	data := pages.BookmarksPageData{
		Bookmarks:  toTemplBookmarks(bookmarks),
		Total:      total,
		Page:       page,
		TotalPages: totalPages(total),
	}
	return render(c, pages.BookmarksListPage(data))
}

func (h *Handler) CreateBookmark(c *fiber.Ctx) error {
	user := GetUser(c)
	title := c.FormValue("title")
	url := c.FormValue("url")
	description := c.FormValue("description")
	tags := c.FormValue("tags")
	notes := c.FormValue("notes")

	if title == "" || url == "" {
		return c.Redirect("/bookmarks")
	}

	var tagSlice []string
	for _, t := range strings.Split(tags, ",") {
		trimmed := strings.TrimSpace(t)
		if trimmed != "" {
			tagSlice = append(tagSlice, trimmed)
		}
	}

	_, _ = h.queries.CreateBookmark(c.Context(), generated.CreateBookmarkParams{
		Title:       title,
		URL:         url,
		Description: sql.NullString{String: description, Valid: description != ""},
		Tags:        tagSlice,
		Notes:       sql.NullString{String: notes, Valid: notes != ""},
		CreatedBy:   user.ID,
	})
	return c.Redirect("/bookmarks")
}

func (h *Handler) UpdateBookmarkForm(c *fiber.Ctx) error {
	id := c.Params("id")
	title := c.FormValue("title")
	url := c.FormValue("url")
	description := c.FormValue("description")
	tags := c.FormValue("tags")
	notes := c.FormValue("notes")

	if title == "" || url == "" {
		return c.Redirect("/bookmarks/" + id)
	}

	var tagSlice []string
	for _, t := range strings.Split(tags, ",") {
		trimmed := strings.TrimSpace(t)
		if trimmed != "" {
			tagSlice = append(tagSlice, trimmed)
		}
	}

	_, _ = h.queries.UpdateBookmark(c.Context(), generated.UpdateBookmarkParams{
		ID:          middleware.StringToUUID(id),
		Title:       title,
		URL:         url,
		Description: sql.NullString{String: description, Valid: description != ""},
		Tags:        tagSlice,
		Notes:       sql.NullString{String: notes, Valid: notes != ""},
	})
	return c.Redirect("/bookmarks/" + id)
}

func (h *Handler) DeleteBookmark(c *fiber.Ctx) error {
	id := c.Params("id")
	_ = h.queries.SoftDeleteBookmark(c.Context(), middleware.StringToUUID(id))
	return c.Redirect("/bookmarks")
}

func toTemplBookmarks(bookmarks []generated.Bookmark) []pages.BookmarkItem {
	items := make([]pages.BookmarkItem, len(bookmarks))
	for i, b := range bookmarks {
		items[i] = pages.BookmarkItem{
			ID:          middleware.UUIDToString(b.ID),
			Title:       b.Title,
			URL:         b.URL,
			Description: nullStr(b.Description),
			Tags:        b.Tags,
			CreatedAt:   formatTime(b.CreatedAt),
		}
	}
	return items
}
