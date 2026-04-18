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

	companies, _ := h.queries.ListCompanies(c.Context(), 50, 0)
	users, _ := h.queries.ListUsers(c.Context(), 50, 0)
	userNames := h.getUserNameMap(c)

	data := pages.TasksPageData{
		Todo:         toTemplTasksWithNames(todo, userNames),
		InProgress:   toTemplTasksWithNames(inProgress, userNames),
		Waiting:      toTemplTasksWithNames(waiting, userNames),
		Review:       toTemplTasksWithNames(review, userNames),
		Done:         toTemplTasksWithNames(done, userNames),
		StatusCounts: counts,
		Companies:    toTemplCompanies(companies),
		Users:        toTemplUsers(users),
		ViewMode:     c.Query("view", "kanban"),
	}
	// For table view, also provide flat list of all tasks
	if data.ViewMode == "table" {
		allTasks, _ := h.queries.ListTasks(c.Context(), 100, 0)
		data.AllTasks = toTemplTasksWithNames(allTasks, userNames)
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
	assigneeID := c.FormValue("assignee_id")
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
		AssigneeID:  middleware.StringToUUID(assigneeID),
		DueDate:     dd,
		CreatedBy:   user.ID,
	})
	return c.Redirect("/tasks")
}

func (h *Handler) UpdateTaskForm(c *fiber.Ctx) error {
	id := c.Params("id")
	title := c.FormValue("title")
	description := c.FormValue("description")
	category := c.FormValue("category")
	priority := c.FormValue("priority")
	dueDate := c.FormValue("due_date")

	if title == "" {
		return c.Redirect("/tasks/" + id)
	}

	var dd pgtype.Date
	if dueDate != "" {
		_ = dd.Scan(dueDate)
	}

	_, _ = h.queries.UpdateTask(c.Context(), generated.UpdateTaskParams{
		ID:          middleware.StringToUUID(id),
		Title:       title,
		Description: sql.NullString{String: description, Valid: description != ""},
		Category:    sql.NullString{String: category, Valid: category != ""},
		Priority:    priority,
		DueDate:     dd,
	})
	return c.Redirect("/tasks/" + id)
}

func (h *Handler) DeleteTask(c *fiber.Ctx) error {
	id := c.Params("id")
	_ = h.queries.SoftDeleteTask(c.Context(), middleware.StringToUUID(id))
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

// userNameMap builds a lookup from UUID to full_name
var cachedUserNames map[string]string

func (h *Handler) getUserNameMap(c *fiber.Ctx) map[string]string {
	users, _ := h.queries.ListUsers(c.Context(), 100, 0)
	m := make(map[string]string)
	for _, u := range users {
		m[middleware.UUIDToString(u.ID)] = u.FullName
	}
	return m
}

func toTemplTasks(tasks []generated.Task) []pages.TaskItem {
	return toTemplTasksWithNames(tasks, nil)
}

func toTemplTasksWithNames(tasks []generated.Task, userNames map[string]string) []pages.TaskItem {
	items := make([]pages.TaskItem, len(tasks))
	for i, t := range tasks {
		assignee := ""
		if userNames != nil && t.AssigneeID.Valid {
			assignee = userNames[middleware.UUIDToString(t.AssigneeID)]
		}
		items[i] = pages.TaskItem{
			ID:           middleware.UUIDToString(t.ID),
			Title:        t.Title,
			Status:       t.Status,
			StatusVi:     LabelVi("task_status", t.Status),
			Priority:     t.Priority,
			PriorityVi:   LabelVi("priority", t.Priority),
			Category:     nullStr(t.Category),
			DueDate:      formatDate(t.DueDate),
			AssigneeName: assignee,
		}
	}
	return items
}
