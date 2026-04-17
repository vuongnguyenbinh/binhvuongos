package handler

import (
	"database/sql"
	"fmt"

	"binhvuongos/internal/db/generated"
	"binhvuongos/internal/middleware"
	"binhvuongos/web/templates/pages"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
)

func (h *Handler) WorkLogs(c *fiber.Ctx) error {
	page, limit, offset := getPage(c)
	items, err := h.queries.ListWorkLogs(c.Context(), limit, offset)
	if err != nil {
		return render(c, pages.WorkLogsListPage(pages.WorkLogsPageData{}))
	}
	total, _ := h.queries.CountWorkLogs(c.Context())
	pendingCount, _ := h.queries.CountPendingWorkLogs(c.Context())
	workTypes, _ := h.queries.ListActiveWorkTypes(c.Context())
	companies, _ := h.queries.ListCompanies(c.Context(), 50, 0)

	data := pages.WorkLogsPageData{
		Items:        toTemplWorkLogs(items, workTypes),
		Total:        total,
		PendingCount: pendingCount,
		Page:         page,
		TotalPages:   totalPages(total),
		Companies:    toTemplCompanies(companies),
		WorkTypes:    toTemplWorkTypes(workTypes),
	}
	return render(c, pages.WorkLogsListPage(data))
}

func (h *Handler) CreateWorkLog(c *fiber.Ctx) error {
	user := GetUser(c)
	workDate := c.FormValue("work_date")
	companyID := c.FormValue("company_id")
	workTypeID := c.FormValue("work_type_id")
	quantity := c.FormValue("quantity")
	notes := c.FormValue("notes")
	sheetURL := c.FormValue("sheet_url")

	if workDate == "" || companyID == "" || workTypeID == "" || quantity == "" {
		return c.Redirect("/work-logs")
	}

	var wd pgtype.Date
	_ = wd.Scan(workDate)
	var qty pgtype.Numeric
	_ = qty.Scan(quantity)

	_, _ = h.queries.CreateWorkLog(c.Context(), generated.CreateWorkLogParams{
		WorkDate:   wd,
		UserID:     user.ID,
		CompanyID:  middleware.StringToUUID(companyID),
		WorkTypeID: middleware.StringToUUID(workTypeID),
		Quantity:   qty,
		Notes:      sql.NullString{String: notes, Valid: notes != ""},
		SheetURL:   sql.NullString{String: sheetURL, Valid: sheetURL != ""},
	})
	return c.Redirect("/work-logs")
}

func (h *Handler) ApproveWorkLogForm(c *fiber.Ctx) error {
	id := c.Params("id")
	user := GetUser(c)
	adminNotes := c.FormValue("admin_notes")
	_, _ = h.queries.ApproveWorkLog(c.Context(), middleware.StringToUUID(id), user.ID,
		sql.NullString{String: adminNotes, Valid: adminNotes != ""})
	// HTMX: return inline status badge
	if c.Get("HX-Request") == "true" {
		return c.SendString(`<span class="pill bg-sage/20 text-forest">ĐÃ DUYỆT</span>`)
	}
	return c.Redirect("/work-logs")
}

func (h *Handler) RejectWorkLogForm(c *fiber.Ctx) error {
	id := c.Params("id")
	user := GetUser(c)
	adminNotes := c.FormValue("admin_notes")
	_, _ = h.queries.RejectWorkLog(c.Context(), middleware.StringToUUID(id), user.ID,
		sql.NullString{String: adminNotes, Valid: adminNotes != ""})
	if c.Get("HX-Request") == "true" {
		return c.SendString(`<span class="pill bg-rust/10 text-rust">TỪ CHỐI</span>`)
	}
	return c.Redirect("/work-logs")
}

func (h *Handler) BatchApproveWorkLogs(c *fiber.Ctx) error {
	user := GetUser(c)
	// Get all submitted work logs and approve them
	submitted, _ := h.queries.ListWorkLogsByStatus(c.Context(), "submitted", 100, 0)
	for _, wl := range submitted {
		_, _ = h.queries.ApproveWorkLog(c.Context(), wl.ID, user.ID, sql.NullString{})
	}
	return c.Redirect("/work-logs")
}

func toTemplWorkLogs(items []generated.WorkLog, workTypes []generated.WorkType) []pages.WorkLogItem {
	// Build work type lookup
	wtMap := make(map[string]generated.WorkType)
	for _, wt := range workTypes {
		wtMap[middleware.UUIDToString(wt.ID)] = wt
	}

	result := make([]pages.WorkLogItem, len(items))
	for i, wl := range items {
		wtID := middleware.UUIDToString(wl.WorkTypeID)
		wt := wtMap[wtID]

		var qtyStr string
		if wl.Quantity.Valid {
			f, _ := wl.Quantity.Float64Value()
			qtyStr = fmt.Sprintf("%.0f", f.Float64)
		}

		result[i] = pages.WorkLogItem{
			ID:           middleware.UUIDToString(wl.ID),
			WorkDate:     formatDate(wl.WorkDate),
			WorkTypeName: wt.Name,
			WorkTypeIcon: nullStr(wt.Icon),
			Unit:         wt.Unit,
			Quantity:     qtyStr,
			Status:       wl.Status,
			Notes:        nullStr(wl.Notes),
		}
	}
	return result
}

func toTemplWorkTypes(wts []generated.WorkType) []pages.WorkTypeItem {
	items := make([]pages.WorkTypeItem, len(wts))
	for i, wt := range wts {
		items[i] = pages.WorkTypeItem{
			ID:   middleware.UUIDToString(wt.ID),
			Name: wt.Name,
			Icon: nullStr(wt.Icon),
			Unit: wt.Unit,
		}
	}
	return items
}
