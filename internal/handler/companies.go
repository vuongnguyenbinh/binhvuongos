package handler

import (
	"database/sql"

	"binhvuongos/internal/db/generated"
	"binhvuongos/internal/drive"
	"binhvuongos/internal/middleware"
	"binhvuongos/web/templates/pages"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Companies(c *fiber.Ctx) error {
	// Filter: show=active (default) | all | archived
	show := c.Query("show", "active")
	var companies []generated.Company
	var err error
	switch show {
	case "all":
		companies, err = h.queries.ListCompanies(c.Context(), 200, 0)
	case "archived":
		companies, err = h.queries.ListCompaniesByStatus(c.Context(), "archived")
	default:
		show = "active"
		companies, err = h.queries.ListCompaniesByStatus(c.Context(), "active")
	}
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
		Show:      show,
	}
	return render(c, pages.CompaniesListPage(data))
}

// ArchiveCompany soft-archives a company by flipping status; data liên kết giữ nguyên.
// Owner + manager only.
func (h *Handler) ArchiveCompany(c *fiber.Ctx) error {
	actor := GetUser(c)
	if actor.Role != "owner" && actor.Role != "manager" {
		return c.Status(403).SendString("Forbidden")
	}
	id := middleware.StringToUUID(c.Params("id"))
	if err := h.queries.UpdateCompanyStatus(c.Context(), id, "archived"); err != nil {
		return c.Status(500).SendString("Lỗi lưu trữ")
	}
	return c.Redirect("/companies")
}

// UnarchiveCompany restores status='active'.
func (h *Handler) UnarchiveCompany(c *fiber.Ctx) error {
	actor := GetUser(c)
	if actor.Role != "owner" && actor.Role != "manager" {
		return c.Status(403).SendString("Forbidden")
	}
	id := middleware.StringToUUID(c.Params("id"))
	if err := h.queries.UpdateCompanyStatus(c.Context(), id, "active"); err != nil {
		return c.Status(500).SendString("Lỗi khôi phục")
	}
	return c.Redirect("/companies?show=archived")
}

// allowedLogoMimes is the closed whitelist for company logo uploads.
// SVG renders via <img> (not <object>/<iframe>), so XSS risk is negligible.
var allowedLogoMimes = map[string]bool{
	"image/png":     true,
	"image/jpeg":    true,
	"image/svg+xml": true,
}

// UploadCompanyLogo uploads an image to Drive and stores the URL in companies.logo_url.
// Owner + manager only. PNG/JPEG/SVG accepted, max 5MB.
func (h *Handler) UploadCompanyLogo(c *fiber.Ctx) error {
	actor := GetUser(c)
	if actor.Role != "owner" && actor.Role != "manager" {
		return c.Status(403).SendString("Forbidden")
	}
	idStr := c.Params("id")
	id := middleware.StringToUUID(idStr)
	file, err := c.FormFile("logo")
	if err != nil {
		return c.Status(400).SendString("Thiếu file logo")
	}
	if file.Size > 5*1024*1024 {
		return c.Status(400).SendString("Logo quá lớn (tối đa 5MB)")
	}
	mime := file.Header.Get("Content-Type")
	if !allowedLogoMimes[mime] {
		return c.Status(400).SendString("Định dạng không hỗ trợ. Chỉ nhận PNG, JPEG, SVG.")
	}
	if h.config.GoogleRefreshToken == "" {
		return c.Status(503).SendString("Drive chưa được cấu hình")
	}
	src, err := file.Open()
	if err != nil {
		return c.Status(500).SendString("Không mở được file")
	}
	defer src.Close()

	cfg := &drive.Config{
		ClientID:     h.config.GoogleClientID,
		ClientSecret: h.config.GoogleClientSecret,
		RefreshToken: h.config.GoogleRefreshToken,
		FolderID:     h.config.GoogleDriveFolderID,
	}
	result, err := drive.UploadFile(c.Context(), cfg, file.Filename, mime, src)
	if err != nil {
		return c.Status(500).SendString("Upload fail: " + err.Error())
	}
	logoURL := result.WebViewLink
	if logoURL == "" {
		logoURL = "https://drive.google.com/file/d/" + result.FileID + "/view"
	}
	if err := h.queries.UpdateCompanyLogo(c.Context(), id, logoURL); err != nil {
		return c.Status(500).SendString("Lỗi lưu URL")
	}
	return c.Redirect("/companies/" + idStr)
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
		label, class := deadlineBadge(c.EndDate)
		items[i] = pages.CompanyItem{
			ID:            middleware.UUIDToString(c.ID),
			Name:          c.Name,
			ShortCode:     nullStr(c.ShortCode),
			Industry:      nullStr(c.Industry),
			MyRole:        c.MyRole,
			Status:        c.Status,
			Health:        nullStr(c.Health),
			LogoURL:       nullStr(c.LogoURL),
			DeadlineLabel: label,
			DeadlineClass: class,
			EndDate:       formatDate(c.EndDate),
		}
	}
	return items
}
