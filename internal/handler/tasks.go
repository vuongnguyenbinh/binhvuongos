package handler

import (
	"database/sql"

	"binhvuongos/internal/db/generated"
	"binhvuongos/internal/middleware"
	"binhvuongos/web/templates/pages"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
)

func (h *Handler) Tasks(c *fiber.Ctx) error {
	// Get tasks grouped by status for kanban
	todo, _ := h.queries.ListTasksByStatus(c.Context(), "todo")
	inProgress, _ := h.queries.ListTasksByStatus(c.Context(), "in_progress")
	waiting, _ := h.queries.ListTasksByStatus(c.Context(), "waiting")
	review, _ := h.queries.ListTasksByStatus(c.Context(), "review")
	done, _ := h.queries.ListTasksByStatus(c.Context(), "done")
	statusCounts, _ := h.queries.CountTasksByStatus(c.Context())

	counts := make(map[string]int64)
	for _, sc := range statusCounts {
		counts[sc.Status] = sc.Count
	}

	data := pages.TasksPageData{
		Todo:         toTemplTasks(todo),
		InProgress:   toTemplTasks(inProgress),
		Waiting:      toTemplTasks(waiting),
		Review:       toTemplTasks(review),
		Done:         toTemplTasks(done),
		StatusCounts: counts,
	}
	return render(c, pages.TasksListPage(data))
}

func (h *Handler) CreateTask(c *fiber.Ctx) error {
	user := GetUser(c)
	title := c.FormValue("title")
	description := c.FormValue("description")
	category := c.FormValue("category")
	priority := c.FormValue("priority")
	companyID := c.FormValue("company_id")
	dueDate := c.FormValue("due_date")

	if title == "" {
		return c.Redirect("/tasks")
	}
	if priority == "" {
		priority = "normal"
	}

	var dd pgtype.Date
	if dueDate != "" {
		_ = dd.Scan(dueDate)
	}

	_, _ = h.queries.CreateTask(c.Context(), generated.CreateTaskParams{
		Title:       title,
		Description: sql.NullString{String: description, Valid: description != ""},
		Category:    sql.NullString{String: category, Valid: category != ""},
		Status:      "todo",
		Priority:    priority,
		CompanyID:   middleware.StringToUUID(companyID),
		DueDate:     dd,
		CreatedBy:   user.ID,
	})
	return c.Redirect("/tasks")
}

func (h *Handler) UpdateTaskStatusForm(c *fiber.Ctx) error {
	id := c.Params("id")
	status := c.FormValue("status")
	if status == "" {
		return c.Redirect("/tasks")
	}
	_, _ = h.queries.UpdateTaskStatus(c.Context(), middleware.StringToUUID(id), status)
	return c.Redirect("/tasks")
}

func toTemplTasks(tasks []generated.Task) []pages.TaskItem {
	items := make([]pages.TaskItem, len(tasks))
	for i, t := range tasks {
		items[i] = pages.TaskItem{
			ID:       middleware.UUIDToString(t.ID),
			Title:    t.Title,
			Status:   t.Status,
			Priority: t.Priority,
			Category: nullStr(t.Category),
			DueDate:  formatDate(t.DueDate),
		}
	}
	return items
}
