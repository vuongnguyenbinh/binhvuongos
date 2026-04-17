package handler

import (
	"binhvuongos/internal/middleware"
	"binhvuongos/web/templates/pages"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Search(c *fiber.Ctx) error {
	q := c.Query("q")
	if q == "" {
		return render(c, pages.SearchPage(pages.SearchPageData{}))
	}

	var results []pages.SearchResult

	// Search tasks by title
	tasks, _ := h.queries.ListTasks(c.Context(), 50, 0)
	for _, t := range tasks {
		if containsInsensitive(t.Title, q) {
			results = append(results, pages.SearchResult{
				Type:  "task",
				ID:    middleware.UUIDToString(t.ID),
				Title: t.Title,
				URL:   "/tasks/" + middleware.UUIDToString(t.ID),
				Extra: t.Status,
			})
		}
	}

	// Search content by title
	content, _ := h.queries.ListContent(c.Context(), 50, 0)
	for _, ct := range content {
		if containsInsensitive(ct.Title, q) {
			results = append(results, pages.SearchResult{
				Type:  "content",
				ID:    middleware.UUIDToString(ct.ID),
				Title: ct.Title,
				URL:   "/content/" + middleware.UUIDToString(ct.ID),
				Extra: ct.Status,
			})
		}
	}

	// Search knowledge (full-text search)
	knowledge, _ := h.queries.SearchKnowledgeItems(c.Context(), q, 20, 0)
	for _, k := range knowledge {
		results = append(results, pages.SearchResult{
			Type:  "knowledge",
			ID:    middleware.UUIDToString(k.ID),
			Title: k.Title,
			URL:   "/knowledge/" + middleware.UUIDToString(k.ID),
			Extra: k.Category,
		})
	}

	// Search bookmarks by title
	bookmarks, _ := h.queries.ListBookmarks(c.Context(), 50, 0)
	for _, b := range bookmarks {
		if containsInsensitive(b.Title, q) || containsInsensitive(b.URL, q) {
			results = append(results, pages.SearchResult{
				Type:  "bookmark",
				ID:    middleware.UUIDToString(b.ID),
				Title: b.Title,
				URL:   "/bookmarks/" + middleware.UUIDToString(b.ID),
				Extra: b.URL,
			})
		}
	}

	data := pages.SearchPageData{
		Query:   q,
		Results: results,
		Total:   len(results),
	}
	return render(c, pages.SearchPage(data))
}

func containsInsensitive(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 &&
		contains(toLower(s), toLower(substr))
}

func toLower(s string) string {
	b := make([]byte, len(s))
	for i := range s {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 32
		}
		b[i] = c
	}
	return string(b)
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || indexStr(s, substr) >= 0)
}

func indexStr(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
