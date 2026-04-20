package handler

import (
	"database/sql"
	"encoding/json"
	"strings"

	"binhvuongos/internal/db/generated"
	"binhvuongos/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
)

// APIDashboard returns dashboard stats as JSON
func (h *Handler) APIDashboard(c *fiber.Ctx) error {
	counts, err := h.queries.GetDashboardCounts(c.Context())
	if err != nil {
		return c.Status(500).JSON(APIError("DB_ERROR", err.Error()))
	}
	return c.JSON(APISuccess(counts))
}

// APIListCompanies returns companies list as JSON
func (h *Handler) APIListCompanies(c *fiber.Ctx) error {
	companies, err := h.queries.ListCompanies(c.Context(), 100, 0)
	if err != nil {
		return c.Status(500).JSON(APIError("DB_ERROR", err.Error()))
	}
	return c.JSON(APISuccess(companies))
}

// APIListTasks returns tasks list as JSON
func (h *Handler) APIListTasks(c *fiber.Ctx) error {
	tasks, err := h.queries.ListTasks(c.Context(), 100, 0)
	if err != nil {
		return c.Status(500).JSON(APIError("DB_ERROR", err.Error()))
	}
	return c.JSON(APISuccess(tasks))
}

// APICreateInbox is the unified inbox webhook endpoint.
// Accepts application/json (text + pre-uploaded attachment URLs) or multipart/form-data (file upload).
// Idempotency: when external_ref is provided, retries with the same (source, external_ref) return the
// existing row instead of inserting a duplicate.
// Attribution: submitted_by is always set to the configured owner user (single-tenant by design).
func (h *Handler) APICreateInbox(c *fiber.Ctx) error {
	var (
		content, url, source, itemType, externalRef string
		attachments                                  []map[string]any
	)

	contentType := c.Get("Content-Type")
	if strings.HasPrefix(contentType, "multipart/form-data") {
		content = c.FormValue("content")
		url = c.FormValue("url")
		source = c.FormValue("source")
		itemType = c.FormValue("item_type")
		externalRef = c.FormValue("external_ref")
		if source == "" {
			source = "api"
		}

		// Optional file → Google Drive
		if file, err := c.FormFile("file"); err == nil && file != nil {
			att, uerr := h.uploadAttachmentToDrive(c, file)
			if uerr != nil {
				return c.Status(400).JSON(APIError("UPLOAD_ERROR", uerr.Error()))
			}
			attachments = append(attachments, att)
		}

		// Pre-uploaded attachment URLs via repeatable form field
		if mf, err := c.MultipartForm(); err == nil && mf != nil {
			for _, u := range mf.Value["attachment_urls"] {
				attachments = append(attachments, map[string]any{"url": u})
			}
		}
	} else {
		var input struct {
			Content        string   `json:"content"`
			URL            string   `json:"url"`
			Source         string   `json:"source"`
			ItemType       string   `json:"item_type"`
			AttachmentURLs []string `json:"attachment_urls"`
			ExternalRef    string   `json:"external_ref"`
		}
		if err := c.BodyParser(&input); err != nil {
			return c.Status(400).JSON(APIError("VALIDATION_ERROR", err.Error()))
		}
		content, url, source, itemType, externalRef = input.Content, input.URL, input.Source, input.ItemType, input.ExternalRef
		if source == "" {
			source = "api"
		}
		for _, u := range input.AttachmentURLs {
			attachments = append(attachments, map[string]any{"url": u})
		}
	}

	// Collect attachment URLs for validation
	attachURLs := make([]string, 0, len(attachments))
	for _, a := range attachments {
		if u, ok := a["url"].(string); ok {
			attachURLs = append(attachURLs, u)
		}
	}
	if err := validateInboxInput(content, source, itemType, externalRef, attachURLs); err != nil {
		return c.Status(400).JSON(APIError("VALIDATION_ERROR", err.Error()))
	}

	if itemType == "" {
		itemType = detectItemType(content, url, attachments)
	}

	attachJSON, _ := json.Marshal(attachments)
	if attachments == nil {
		attachJSON = []byte("[]")
	}

	params := generated.CreateInboxItemParams{
		Content:     content,
		URL:         sql.NullString{String: url, Valid: url != ""},
		Source:      sql.NullString{String: source, Valid: true},
		ItemType:    sql.NullString{String: itemType, Valid: true},
		SubmittedBy: h.ownerUserID,
		Attachments: attachJSON,
		ExternalRef: sql.NullString{String: externalRef, Valid: externalRef != ""},
	}

	// Atomic idempotent insert when external_ref supplied; regular insert otherwise.
	if externalRef != "" {
		item, inserted, err := h.queries.UpsertInboxItemByExternalRef(c.Context(), params)
		if err != nil {
			return c.Status(500).JSON(APIError("DB_ERROR", "failed to persist inbox item"))
		}
		if !inserted {
			return c.Status(200).JSON(fiber.Map{
				"success":   true,
				"duplicate": true,
				"data":      item,
			})
		}
		return c.Status(201).JSON(APISuccess(item))
	}

	item, err := h.queries.CreateInboxItem(c.Context(), params)
	if err != nil {
		return c.Status(500).JSON(APIError("DB_ERROR", "failed to persist inbox item"))
	}
	return c.Status(201).JSON(APISuccess(item))
}

// APICreateBookmark creates a bookmark via JSON API
func (h *Handler) APICreateBookmark(c *fiber.Ctx) error {
	var input struct {
		Title       string   `json:"title"`
		URL         string   `json:"url"`
		Description string   `json:"description"`
		Tags        []string `json:"tags"`
		Notes       string   `json:"notes"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(APIError("VALIDATION_ERROR", err.Error()))
	}
	if input.Title == "" || input.URL == "" {
		return c.Status(400).JSON(APIError("VALIDATION_ERROR", "title and url are required"))
	}

	bm, err := h.queries.CreateBookmark(c.Context(), generated.CreateBookmarkParams{
		Title:       input.Title,
		URL:         input.URL,
		Description: sql.NullString{String: input.Description, Valid: input.Description != ""},
		Tags:        input.Tags,
		Notes:       sql.NullString{String: input.Notes, Valid: input.Notes != ""},
	})
	if err != nil {
		return c.Status(500).JSON(APIError("DB_ERROR", err.Error()))
	}
	return c.Status(201).JSON(APISuccess(bm))
}

// APICreateWorkLog creates a work log via JSON API
func (h *Handler) APICreateWorkLog(c *fiber.Ctx) error {
	var input struct {
		WorkDate   string  `json:"work_date"`
		UserID     string  `json:"user_id"`
		CompanyID  string  `json:"company_id"`
		WorkTypeID string  `json:"work_type_id"`
		CampaignID string  `json:"campaign_id"`
		Quantity   float64 `json:"quantity"`
		Notes      string  `json:"notes"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(APIError("VALIDATION_ERROR", err.Error()))
	}

	var workDate pgtype.Date
	if err := workDate.Scan(input.WorkDate); err != nil {
		return c.Status(400).JSON(APIError("VALIDATION_ERROR", "invalid work_date format"))
	}

	var qty pgtype.Numeric
	_ = qty.Scan(input.Quantity)

	wl, err := h.queries.CreateWorkLog(c.Context(), generated.CreateWorkLogParams{
		WorkDate:   workDate,
		UserID:     middleware.StringToUUID(input.UserID),
		CompanyID:  middleware.StringToUUID(input.CompanyID),
		WorkTypeID: middleware.StringToUUID(input.WorkTypeID),
		CampaignID: middleware.StringToUUID(input.CampaignID),
		Quantity:   qty,
		Notes:      sql.NullString{String: input.Notes, Valid: input.Notes != ""},
	})
	if err != nil {
		return c.Status(500).JSON(APIError("DB_ERROR", err.Error()))
	}
	return c.Status(201).JSON(APISuccess(wl))
}

// APICreateKnowledge creates a knowledge item via JSON API
func (h *Handler) APICreateKnowledge(c *fiber.Ctx) error {
	var input struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Body        string   `json:"body"`
		Category    string   `json:"category"`
		Topics      []string `json:"topics"`
		Scope       string   `json:"scope"`
		SourceURL   string   `json:"source_url"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(APIError("VALIDATION_ERROR", err.Error()))
	}
	if input.Title == "" || input.Category == "" {
		return c.Status(400).JSON(APIError("VALIDATION_ERROR", "title and category are required"))
	}

	scope := "shared"
	if input.Scope != "" {
		scope = input.Scope
	}

	ki, err := h.queries.CreateKnowledgeItem(c.Context(), generated.CreateKnowledgeItemParams{
		Title:       input.Title,
		Description: sql.NullString{String: input.Description, Valid: input.Description != ""},
		Body:        sql.NullString{String: input.Body, Valid: input.Body != ""},
		Category:    input.Category,
		Topics:      input.Topics,
		Scope:       scope,
		SourceURL:   sql.NullString{String: input.SourceURL, Valid: input.SourceURL != ""},
	})
	if err != nil {
		return c.Status(500).JSON(APIError("DB_ERROR", err.Error()))
	}
	return c.Status(201).JSON(APISuccess(ki))
}
