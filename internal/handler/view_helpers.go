package handler

import (
	"database/sql"
	"strconv"
	"strings"

	"binhvuongos/internal/db/generated"
	"binhvuongos/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
)

// toCompanyView converts DB Company to view struct
func toCompanyView(c generated.Company) CompanyView {
	return CompanyView{
		ID:        middleware.UUIDToString(c.ID),
		Name:      c.Name,
		ShortCode: nullStr(c.ShortCode),
		Industry:  nullStr(c.Industry),
		MyRole:    c.MyRole,
		Status:    c.Status,
		Health:    nullStr(c.Health),
		Scope:     strings.Join(c.Scope, " · "),
		StartDate: formatDate(c.StartDate),
	}
}

// toCompanyViews converts slice of DB companies to views
func toCompanyViews(companies []generated.Company) []CompanyView {
	views := make([]CompanyView, len(companies))
	for i, c := range companies {
		views[i] = toCompanyView(c)
	}
	return views
}

// toInboxItemView converts DB InboxItem to view struct
func toInboxItemView(item generated.InboxItem) InboxItemView {
	return InboxItemView{
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

// toInboxItemViews converts slice
func toInboxItemViews(items []generated.InboxItem) []InboxItemView {
	views := make([]InboxItemView, len(items))
	for i, item := range items {
		views[i] = toInboxItemView(item)
	}
	return views
}

// toBookmarkView converts DB Bookmark to view struct
func toBookmarkView(b generated.Bookmark) BookmarkView {
	return BookmarkView{
		ID:          middleware.UUIDToString(b.ID),
		Title:       b.Title,
		URL:         b.URL,
		Description: nullStr(b.Description),
		Tags:        b.Tags,
		Notes:       nullStr(b.Notes),
		CreatedAt:   formatTime(b.CreatedAt),
	}
}

// toBookmarkViews converts slice
func toBookmarkViews(bookmarks []generated.Bookmark) []BookmarkView {
	views := make([]BookmarkView, len(bookmarks))
	for i, b := range bookmarks {
		views[i] = toBookmarkView(b)
	}
	return views
}

// toTaskView converts DB Task to view struct
func toTaskView(t generated.Task) TaskView {
	return TaskView{
		ID:          middleware.UUIDToString(t.ID),
		Title:       t.Title,
		Description: nullStr(t.Description),
		Category:    nullStr(t.Category),
		GroupName:   nullStr(t.GroupName),
		Status:      t.Status,
		Priority:    t.Priority,
		DueDate:     formatDate(t.DueDate),
	}
}

// toTaskViews converts slice
func toTaskViews(tasks []generated.Task) []TaskView {
	views := make([]TaskView, len(tasks))
	for i, t := range tasks {
		views[i] = toTaskView(t)
	}
	return views
}

// nullStr extracts string from sql.NullString
func nullStr(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

// formatDate formats pgtype.Date to DD/MM/YYYY
func formatDate(d pgtype.Date) string {
	if !d.Valid {
		return ""
	}
	return d.Time.Format("02/01/2006")
}

const pageSize = 20

// getPage extracts page number from query param, returns (page, limit, offset)
func getPage(c *fiber.Ctx) (int, int32, int32) {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	if page < 1 {
		page = 1
	}
	return page, int32(pageSize), int32((page - 1) * pageSize)
}

// totalPages calculates total pages from total count
func totalPages(total int64) int {
	pages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		pages++
	}
	if pages < 1 {
		pages = 1
	}
	return pages
}
