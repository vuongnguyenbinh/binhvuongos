package handler

import (
	"database/sql"

	"binhvuongos/internal/db/generated"
	"binhvuongos/internal/middleware"
	"binhvuongos/web/templates/pages"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
)

// AdminSettings shows settings page
func (h *Handler) AdminSettings(c *fiber.Ctx) error {
	settings, _ := h.queries.ListSettings(c.Context())
	settingsMap := make(map[string]string)
	for _, s := range settings {
		settingsMap[s.Key] = s.Value
	}
	return render(c, pages.AdminSettingsPage(settingsMap, ""))
}

// SaveSettings handles POST /admin/settings
func (h *Handler) SaveSettings(c *fiber.Ctx) error {
	keys := []string{
		"smtp_host", "smtp_port", "smtp_user", "smtp_pass", "smtp_from",
		"notion_api_key", "notion_database_ids",
		"n8n_webhook_url",
		"google_oauth_client_id", "google_oauth_client_secret",
		"unsplash_keywords",
	}
	for _, key := range keys {
		val := c.FormValue(key)
		if val != "" || c.FormValue(key+"_clear") == "1" {
			_ = h.queries.SetSetting(c.Context(), key, val)
		}
	}

	settings, _ := h.queries.ListSettings(c.Context())
	settingsMap := make(map[string]string)
	for _, s := range settings {
		settingsMap[s.Key] = s.Value
	}
	return render(c, pages.AdminSettingsPage(settingsMap, "Đã lưu cài đặt!"))
}

// AdminWorkTypes shows work types management
func (h *Handler) AdminWorkTypes(c *fiber.Ctx) error {
	types, _ := h.queries.ListActiveWorkTypes(c.Context())
	return render(c, pages.AdminWorkTypesPage(toAdminWorkTypes(types)))
}

// CreateWorkType handles POST /admin/work-types
func (h *Handler) CreateWorkType(c *fiber.Ctx) error {
	name := c.FormValue("name")
	slug := c.FormValue("slug")
	unit := c.FormValue("unit")
	icon := c.FormValue("icon")
	color := c.FormValue("color")

	if name == "" || slug == "" || unit == "" {
		return c.Redirect("/admin/work-types")
	}

	_ = h.queries.CreateWorkType(c.Context(), name, slug, unit, icon, color)
	return c.Redirect("/admin/work-types")
}

// UpdateWorkType handles POST /admin/work-types/:id
func (h *Handler) UpdateWorkType(c *fiber.Ctx) error {
	id := c.Params("id")
	name := c.FormValue("name")
	unit := c.FormValue("unit")
	icon := c.FormValue("icon")
	color := c.FormValue("color")

	if name == "" {
		return c.Redirect("/admin/work-types")
	}

	_ = h.queries.UpdateWorkType(c.Context(), middleware.StringToUUID(id), name, unit, icon, color)
	return c.Redirect("/admin/work-types")
}

// DeleteWorkType handles POST /admin/work-types/:id/delete
func (h *Handler) DeleteWorkType(c *fiber.Ctx) error {
	id := c.Params("id")
	_ = h.queries.DeactivateWorkType(c.Context(), middleware.StringToUUID(id))
	return c.Redirect("/admin/work-types")
}

func toAdminWorkTypes(types []generated.WorkType) []pages.AdminWorkTypeItem {
	items := make([]pages.AdminWorkTypeItem, len(types))
	for i, wt := range types {
		items[i] = pages.AdminWorkTypeItem{
			ID:    middleware.UUIDToString(wt.ID),
			Name:  wt.Name,
			Slug:  wt.Slug,
			Unit:  wt.Unit,
			Icon:  nullStr(wt.Icon),
			Color: nullStr(wt.Color),
		}
	}
	return items
}

// Suppress unused import
var _ = sql.NullString{}
var _ = pgtype.UUID{}
