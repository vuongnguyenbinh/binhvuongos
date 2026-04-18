package handler

import (
	"database/sql"

	"binhvuongos/internal/db/generated"
	"binhvuongos/internal/middleware"
	"binhvuongos/web/templates/pages"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Companies(c *fiber.Ctx) error {
	companies, err := h.queries.ListCompanies(c.Context(), 50, 0)
	if err != nil {
		return render(c, pages.CompaniesListPage(pages.CompaniesPageData{}))
	}
	total, _ := h.queries.CountCompanies(c.Context())

	// Get task stats per company
	taskStats, _ := h.queries.GetCompanyTaskStats(c.Context())
	statsMap := make(map[string]pages.CompanyStats)
	for _, s := range taskStats {
		cid := middleware.UUIDToString(s.CompanyID)
		pct := 0
		if s.TotalTasks > 0 {
			pct = int(s.DoneTasks * 100 / s.TotalTasks)
		}
		statsMap[cid] = pages.CompanyStats{
			OpenTasks:      s.OpenTasks,
			CompletionPct:  pct,
		}
	}

	compViews := toTemplCompanies(companies)
	for i := range compViews {
		if stats, ok := statsMap[compViews[i].ID]; ok {
			compViews[i].OpenTasks = stats.OpenTasks
			compViews[i].CompletionPct = stats.CompletionPct
		}
	}

	data := pages.CompaniesPageData{
		Companies: compViews,
		Total:     total,
	}
	return render(c, pages.CompaniesListPage(data))
}

func (h *Handler) CreateCompany(c *fiber.Ctx) error {
	user := GetUser(c)
	name := c.FormValue("name")
	shortCode := c.FormValue("short_code")
	industry := c.FormValue("industry")
	myRole := c.FormValue("my_role")
	description := c.FormValue("description")

	if name == "" || myRole == "" {
		return c.Redirect("/companies")
	}

	_, _ = h.queries.CreateCompany(c.Context(), generated.CreateCompanyParams{
		Name:        name,
		ShortCode:   sql.NullString{String: shortCode, Valid: shortCode != ""},
		Slug:        sql.NullString{String: "", Valid: false},
		Industry:    sql.NullString{String: industry, Valid: industry != ""},
		MyRole:      myRole,
		Status:      "active",
		Health:      sql.NullString{String: "ok", Valid: true},
		Description: sql.NullString{String: description, Valid: description != ""},
		CreatedBy:   user.ID,
	})
	return c.Redirect("/companies")
}

func (h *Handler) UpdateCompanyForm(c *fiber.Ctx) error {
	id := c.Params("id")
	name := c.FormValue("name")
	shortCode := c.FormValue("short_code")
	industry := c.FormValue("industry")
	health := c.FormValue("health")
	status := c.FormValue("status")

	if name != "" {
		_, _ = h.queries.UpdateCompany(c.Context(), generated.UpdateCompanyParams{
			ID:        middleware.StringToUUID(id),
			Name:      name,
			ShortCode: sql.NullString{String: shortCode, Valid: shortCode != ""},
			Industry:  sql.NullString{String: industry, Valid: industry != ""},
			MyRole:    "owner", // keep existing
			Status:    status,
			Health:    sql.NullString{String: health, Valid: health != ""},
		})
	}
	return c.Redirect("/companies/" + id)
}

func toTemplCompanies(companies []generated.Company) []pages.CompanyItem {
	items := make([]pages.CompanyItem, len(companies))
	for i, c := range companies {
		items[i] = pages.CompanyItem{
			ID:        middleware.UUIDToString(c.ID),
			Name:      c.Name,
			ShortCode: nullStr(c.ShortCode),
			Industry:  nullStr(c.Industry),
			MyRole:    c.MyRole,
			Status:    c.Status,
			Health:    nullStr(c.Health),
		}
	}
	return items
}
