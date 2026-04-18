package handler

import (
	"database/sql"

	"binhvuongos/internal/db/generated"
	"binhvuongos/internal/middleware"
	"binhvuongos/web/templates/pages"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
)

func (h *Handler) Campaigns(c *fiber.Ctx) error {
	filterCompanyID := c.Query("company_id")
	filterType := c.Query("campaign_type")

	items, err := h.queries.ListCampaignsWithCompany(c.Context(),
		middleware.StringToUUID(filterCompanyID), filterType, 50, 0)
	if err != nil {
		return render(c, pages.CampaignsListPage(pages.CampaignsPageData{}))
	}
	total, _ := h.queries.CountCampaigns(c.Context())
	companies, _ := h.queries.ListCompanies(c.Context(), 50, 0)

	var campViews []pages.CampaignItem
	for _, cw := range items {
		campViews = append(campViews, pages.CampaignItem{
			ID:           middleware.UUIDToString(cw.ID),
			Name:         cw.Name,
			Description:  nullStr(cw.Description),
			CampaignType: nullStr(cw.CampaignType),
			Status:       cw.Status,
			StartDate:    formatDate(cw.StartDate),
			EndDate:      formatDate(cw.EndDate),
			CompanyName:  cw.CompanyName,
			CompanyCode:  cw.CompanyCode,
		})
	}

	data := pages.CampaignsPageData{
		Campaigns:       campViews,
		Total:           total,
		Companies:       toTemplCompanies(companies),
		FilterCompanyID: filterCompanyID,
		FilterType:      filterType,
	}
	return render(c, pages.CampaignsListPage(data))
}

func (h *Handler) CreateCampaign(c *fiber.Ctx) error {
	user := GetUser(c)
	name := c.FormValue("name")
	description := c.FormValue("description")
	companyID := c.FormValue("company_id")
	campaignType := c.FormValue("campaign_type")
	startDate := c.FormValue("start_date")
	endDate := c.FormValue("end_date")

	if name == "" || companyID == "" {
		return c.Redirect("/campaigns")
	}

	var sd, ed pgtype.Date
	if startDate != "" {
		_ = sd.Scan(startDate)
	}
	if endDate != "" {
		_ = ed.Scan(endDate)
	}

	_, _ = h.queries.CreateCampaign(c.Context(), generated.CreateCampaignParams{
		Name:         name,
		Description:  sql.NullString{String: description, Valid: description != ""},
		CompanyID:    middleware.StringToUUID(companyID),
		OwnerID:      user.ID,
		CampaignType: sql.NullString{String: campaignType, Valid: campaignType != ""},
		Status:       "planning",
		StartDate:    sd,
		EndDate:      ed,
		TargetJSON:   []byte("{}"),
		CreatedBy:    user.ID,
	})
	return c.Redirect("/campaigns")
}

func (h *Handler) UpdateCampaignForm(c *fiber.Ctx) error {
	id := c.Params("id")
	name := c.FormValue("name")
	description := c.FormValue("description")
	campaignType := c.FormValue("campaign_type")
	status := c.FormValue("status")
	startDate := c.FormValue("start_date")
	endDate := c.FormValue("end_date")

	if name == "" {
		return c.Redirect("/campaigns/" + id)
	}

	var sd, ed pgtype.Date
	if startDate != "" {
		_ = sd.Scan(startDate)
	}
	if endDate != "" {
		_ = ed.Scan(endDate)
	}

	_, _ = h.queries.UpdateCampaign(c.Context(), generated.UpdateCampaignParams{
		ID:           middleware.StringToUUID(id),
		Name:         name,
		Description:  sql.NullString{String: description, Valid: description != ""},
		CampaignType: sql.NullString{String: campaignType, Valid: campaignType != ""},
		Status:       status,
		StartDate:    sd,
		EndDate:      ed,
	})
	return c.Redirect("/campaigns/" + id)
}

func (h *Handler) DeleteCampaign(c *fiber.Ctx) error {
	id := c.Params("id")
	_ = h.queries.SoftDeleteCampaign(c.Context(), middleware.StringToUUID(id))
	return c.Redirect("/campaigns")
}

func toTemplCampaigns(items []generated.Campaign) []pages.CampaignItem {
	result := make([]pages.CampaignItem, len(items))
	for i, c := range items {
		result[i] = pages.CampaignItem{
			ID:           middleware.UUIDToString(c.ID),
			Name:         c.Name,
			Description:  nullStr(c.Description),
			CampaignType: nullStr(c.CampaignType),
			Status:       c.Status,
			StartDate:    formatDate(c.StartDate),
			EndDate:      formatDate(c.EndDate),
		}
	}
	return result
}
