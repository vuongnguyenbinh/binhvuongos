package handler

import (
	"binhvuongos/internal/middleware"
	"binhvuongos/web/templates/components"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// allowedConvertTargets gates the ?target= query parameter.
var allowedConvertTargets = map[string]bool{
	"task": true, "content": true, "knowledge": true,
}

// ConvertInbox transactionally creates a target row (task/content/knowledge_items)
// and marks the inbox item done with converted_to_type/id pointing back. Idempotent
// via a status guard — concurrent retries get 409 instead of duplicate inserts.
func (h *Handler) ConvertInbox(c *fiber.Ctx) error {
	actor := GetUser(c)
	target := c.Query("target")
	if !allowedConvertTargets[target] {
		return c.Status(400).SendString("Target không hợp lệ")
	}
	inboxID := middleware.StringToUUID(c.Params("id"))
	if !inboxID.Valid {
		return c.Status(400).SendString("ID không hợp lệ")
	}

	tx, err := h.queries.Pool().BeginTx(c.Context(), pgx.TxOptions{})
	if err != nil {
		return c.Status(500).SendString("Không mở được transaction")
	}
	defer tx.Rollback(c.Context())

	var newID pgtype.UUID
	switch target {
	case "task":
		title := strings.TrimSpace(c.FormValue("title"))
		if title == "" {
			return c.Status(400).SendString("Thiếu title")
		}
		priority := c.FormValue("priority")
		if priority == "" {
			priority = "normal"
		}
		err = tx.QueryRow(c.Context(),
			`INSERT INTO tasks (title, description, company_id, priority, due_date, created_by, status)
			 VALUES ($1, $2, $3, $4, $5, $6, 'todo') RETURNING id`,
			truncate(title, 500),
			nullStringFromForm(c, "description"),
			optionalUUID(c.FormValue("company_id")),
			priority,
			optionalDate(c.FormValue("due_date")),
			actor.ID,
		).Scan(&newID)
	case "content":
		title := strings.TrimSpace(c.FormValue("title"))
		companyID := optionalUUID(c.FormValue("company_id"))
		if title == "" || !companyID.Valid {
			return c.Status(400).SendString("Title và company là bắt buộc với content")
		}
		contentType := c.FormValue("content_type")
		if contentType == "" {
			contentType = "blog"
		}
		err = tx.QueryRow(c.Context(),
			`INSERT INTO content (title, content_type, company_id, author_id, status)
			 VALUES ($1, $2, $3, $4, 'idea') RETURNING id`,
			truncate(title, 500), contentType, companyID, actor.ID,
		).Scan(&newID)
	case "knowledge":
		title := strings.TrimSpace(c.FormValue("title"))
		if title == "" {
			return c.Status(400).SendString("Thiếu title")
		}
		category := strings.TrimSpace(c.FormValue("category"))
		if category == "" {
			category = "note"
		}
		err = tx.QueryRow(c.Context(),
			`INSERT INTO knowledge_items (title, body, category, created_by)
			 VALUES ($1, $2, $3, $4) RETURNING id`,
			truncate(title, 500),
			c.FormValue("body"),
			category,
			actor.ID,
		).Scan(&newID)
	}
	if err != nil {
		return c.Status(500).SendString("Lỗi tạo " + target + ": " + err.Error())
	}

	cmd, err := tx.Exec(c.Context(),
		`UPDATE inbox_items
		 SET status='done', converted_to_type=$2, converted_to_id=$3,
		     processed_at=NOW(), triage_notes=$4
		 WHERE id=$1 AND status != 'done'`,
		inboxID, target, newID, nullStringFromForm(c, "triage_notes"))
	if err != nil {
		return c.Status(500).SendString("Lỗi update inbox")
	}
	if cmd.RowsAffected() == 0 {
		// Already converted by a concurrent request — reject to avoid duplicate target row.
		return c.Status(409).SendString("Inbox item đã được xử lý")
	}

	if err := tx.Commit(c.Context()); err != nil {
		return c.Status(500).SendString("Commit fail")
	}
	return c.Redirect("/inbox")
}

// TriageModalPartial returns the HTMX modal HTML to overlay over the inbox page.
func (h *Handler) TriageModalPartial(c *fiber.Ctx) error {
	target := c.Query("target")
	if !allowedConvertTargets[target] {
		return c.Status(400).SendString("Target không hợp lệ")
	}
	inboxID := c.Params("id")
	item, err := h.queries.GetInboxItemByID(c.Context(), middleware.StringToUUID(inboxID))
	if err != nil {
		return c.Status(404).SendString("Không tìm thấy")
	}
	companies, _ := h.queries.ListCompanies(c.Context(), 200, 0)
	opts := make([]components.CompanyOpt, 0, len(companies))
	for _, co := range companies {
		opts = append(opts, components.CompanyOpt{
			ID:   middleware.UUIDToString(co.ID),
			Name: co.Name,
		})
	}
	prefill := truncate(item.Content, 200)
	switch target {
	case "task":
		return render(c, components.TriageTaskModal(inboxID, prefill, opts))
	case "content":
		return render(c, components.TriageContentModal(inboxID, prefill, opts))
	case "knowledge":
		return render(c, components.TriageKnowledgeModal(inboxID, prefill, item.Content))
	}
	return c.Status(400).SendString("Unreachable")
}

// optionalUUID parses a UUID string but returns invalid pgtype.UUID when blank.
func optionalUUID(s string) pgtype.UUID {
	if strings.TrimSpace(s) == "" {
		return pgtype.UUID{}
	}
	return middleware.StringToUUID(s)
}

// optionalDate returns NULL pgtype.Date for blank input; "2006-01-02" format only.
func optionalDate(s string) pgtype.Date {
	s = strings.TrimSpace(s)
	if s == "" {
		return pgtype.Date{}
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return pgtype.Date{}
	}
	return pgtype.Date{Time: t, Valid: true}
}

// truncate returns at most n bytes of s; safe for ASCII-heavy DB columns with hard cap.
func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}

// nullStringFromForm returns a sql.NullString-like pgtype.Text for empty input.
// Using pgx-friendly *string here would also work; pgtype.Text avoids extra import.
func nullStringFromForm(c *fiber.Ctx, field string) any {
	v := strings.TrimSpace(c.FormValue(field))
	if v == "" {
		return nil
	}
	return v
}
