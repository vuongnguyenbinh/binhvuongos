package handler

import (
	"fmt"

	"binhvuongos/internal/db/generated"
	"binhvuongos/internal/middleware"
	"binhvuongos/web/templates/pages"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
)

func (h *Handler) CompanyDetail(c *fiber.Ctx) error {
	id := middleware.StringToUUID(c.Params("id"))
	company, err := h.queries.GetCompanyByID(c.Context(), id)
	if err != nil {
		return c.Redirect("/companies")
	}
	data := pages.CompanyDetailData{
		Company: companyToDetail(company),
	}
	// Load tasks for this company
	tasks, _ := h.queries.ListTasksByCompany(c.Context(), id, 10, 0)
	data.RecentTasks = toTemplTasks(tasks)
	return render(c, pages.CompanyDetailDataPage(data))
}

func (h *Handler) TaskDetail(c *fiber.Ctx) error {
	id := middleware.StringToUUID(c.Params("id"))
	task, err := h.queries.GetTaskByID(c.Context(), id)
	if err != nil {
		return c.Redirect("/tasks")
	}
	item := pages.TaskItem{
		ID:       middleware.UUIDToString(task.ID),
		Title:    task.Title,
		Status:   task.Status,
		Priority: task.Priority,
		Category: nullStr(task.Category),
		DueDate:  formatDate(task.DueDate),
	}
	desc := nullStr(task.Description)
	return render(c, pages.TaskDetailDataPage(item, desc))
}

func (h *Handler) ContentDetail(c *fiber.Ctx) error {
	id := middleware.StringToUUID(c.Params("id"))
	content, err := h.queries.GetContentByID(c.Context(), id)
	if err != nil {
		return c.Redirect("/content")
	}
	item := pages.ContentItem{
		ID:          middleware.UUIDToString(content.ID),
		Title:       content.Title,
		ContentType: content.ContentType,
		Status:      content.Status,
		PublishDate: formatDate(content.PublishDate),
	}
	return render(c, pages.ContentDetailDataPage(item))
}

func (h *Handler) InboxDetail(c *fiber.Ctx) error {
	id := middleware.StringToUUID(c.Params("id"))
	item, err := h.queries.GetInboxItemByID(c.Context(), id)
	if err != nil {
		return c.Redirect("/inbox")
	}
	data := pages.InboxItemData{
		ID:        middleware.UUIDToString(item.ID),
		Content:   item.Content,
		URL:       nullStr(item.URL),
		Source:    nullStr(item.Source),
		ItemType:  nullStr(item.ItemType),
		Status:    item.Status,
		CreatedAt: formatTime(item.CreatedAt),
		TimeAgo:   timeAgo(item.CreatedAt),
	}
	return render(c, pages.InboxDetailDataPage(data))
}

func (h *Handler) InboxCreate(c *fiber.Ctx) error {
	return render(c, pages.InboxCreatePage())
}

func (h *Handler) WorkLogDetail(c *fiber.Ctx) error {
	id := middleware.StringToUUID(c.Params("id"))
	wl, err := h.queries.GetWorkLogByID(c.Context(), id)
	if err != nil {
		return c.Redirect("/work-logs")
	}
	wt, _ := h.queries.GetWorkTypeByID(c.Context(), wl.WorkTypeID)
	item := pages.WorkLogItem{
		ID:           middleware.UUIDToString(wl.ID),
		WorkDate:     formatDate(wl.WorkDate),
		WorkTypeName: wt.Name,
		WorkTypeIcon: nullStr(wt.Icon),
		Unit:         wt.Unit,
		Quantity:     numericToStr(wl.Quantity),
		Status:       wl.Status,
		Notes:        nullStr(wl.Notes),
	}
	return render(c, pages.WorkLogDetailDataPage(item))
}

func (h *Handler) CampaignDetail(c *fiber.Ctx) error {
	id := middleware.StringToUUID(c.Params("id"))
	camp, err := h.queries.GetCampaignByID(c.Context(), id)
	if err != nil {
		return c.Redirect("/campaigns")
	}
	item := pages.CampaignItem{
		ID:           middleware.UUIDToString(camp.ID),
		Name:         camp.Name,
		Description:  nullStr(camp.Description),
		CampaignType: nullStr(camp.CampaignType),
		Status:       camp.Status,
		StartDate:    formatDate(camp.StartDate),
		EndDate:      formatDate(camp.EndDate),
	}
	return render(c, pages.CampaignDetailDataPage(item))
}

func (h *Handler) KnowledgeDetail(c *fiber.Ctx) error {
	id := middleware.StringToUUID(c.Params("id"))
	ki, err := h.queries.GetKnowledgeItemByID(c.Context(), id)
	if err != nil {
		return c.Redirect("/knowledge")
	}
	item := pages.KnowledgeItemData{
		ID:          middleware.UUIDToString(ki.ID),
		Title:       ki.Title,
		Description: nullStr(ki.Description),
		Category:    ki.Category,
		Topics:      ki.Topics,
		Scope:       ki.Scope,
		SourceURL:   nullStr(ki.SourceURL),
		CreatedAt:   formatTime(ki.CreatedAt),
	}
	return render(c, pages.KnowledgeDetailDataPage(item))
}

func (h *Handler) BookmarkDetail(c *fiber.Ctx) error {
	id := middleware.StringToUUID(c.Params("id"))
	bm, err := h.queries.GetBookmarkByID(c.Context(), id)
	if err != nil {
		return c.Redirect("/bookmarks")
	}
	item := pages.BookmarkItem{
		ID:          middleware.UUIDToString(bm.ID),
		Title:       bm.Title,
		URL:         bm.URL,
		Description: nullStr(bm.Description),
		Tags:        bm.Tags,
		CreatedAt:   formatTime(bm.CreatedAt),
	}
	return render(c, pages.BookmarkDetailDataPage(item))
}

func companyToDetail(c generated.Company) pages.CompanyItem {
	return pages.CompanyItem{
		ID:        middleware.UUIDToString(c.ID),
		Name:      c.Name,
		ShortCode: nullStr(c.ShortCode),
		Industry:  nullStr(c.Industry),
		MyRole:    c.MyRole,
		Status:    c.Status,
		Health:    nullStr(c.Health),
	}
}

func numericToStr(n pgtype.Numeric) string {
	if !n.Valid {
		return "0"
	}
	f, err := n.Float64Value()
	if err != nil {
		return "0"
	}
	return fmt.Sprintf("%.0f", f.Float64)
}
