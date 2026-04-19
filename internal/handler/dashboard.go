package handler

import (
	"fmt"
	"time"

	"binhvuongos/internal/middleware"
	"binhvuongos/web/templates/pages"

	"github.com/gofiber/fiber/v2"
)

// getGreeting returns Vietnamese greeting + background image path based on Hanoi time
func getGreeting() (string, string) {
	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		loc = time.FixedZone("ICT", 7*3600)
	}
	hour := time.Now().In(loc).Hour()

	switch {
	case hour >= 5 && hour < 11:
		return "Chào buổi sáng", "/static/img/morning.jpg"
	case hour >= 11 && hour < 13:
		return "Chào buổi trưa", "/static/img/noon.jpg"
	case hour >= 13 && hour < 18:
		return "Chào buổi chiều", "/static/img/afternoon.jpg"
	case hour >= 18 && hour < 22:
		return "Chào buổi tối", "/static/img/evening.jpg"
	default:
		return "Chào đêm muộn", "/static/img/night.jpg"
	}
}

func (h *Handler) Dashboard(c *fiber.Ctx) error {
	user := GetUser(c)

	// Role-based dashboard
	switch user.Role {
	case "owner", "core_staff":
		return h.ownerDashboard(c)
	default:
		return h.userDashboard(c)
	}
}

func (h *Handler) ownerDashboard(c *fiber.Ctx) error {
	u := GetUser(c)
	counts, err := h.queries.GetDashboardCounts(c.Context())
	if err != nil {
		return render(c, pages.DashboardPage())
	}

	companies, _ := h.queries.ListCompanies(c.Context(), 10, 0)
	todayTasks, _ := h.queries.ListTasksDueToday(c.Context())
	monthOutput, _ := h.queries.GetDashboardOutputThisMonth(c.Context())
	campaigns, _ := h.queries.ListCampaignsByStatus(c.Context(), "running")

	var outputItems []pages.DashOutputItem
	for _, o := range monthOutput {
		var total string
		if o.Total.Valid {
			f, _ := o.Total.Float64Value()
			total = fmt.Sprintf("%.0f", f.Float64)
		} else {
			total = "0"
		}
		outputItems = append(outputItems, pages.DashOutputItem{
			Name: o.Name, Icon: o.Icon, Unit: o.Unit, Total: total,
		})
	}

	// Build company name lookup
	companyNames := make(map[string]string)
	for _, c := range companies {
		companyNames[middleware.UUIDToString(c.ID)] = c.Name
	}

	var campItems []pages.DashCampaignItem
	for _, camp := range campaigns {
		cid := middleware.UUIDToString(camp.CompanyID)
		campItems = append(campItems, pages.DashCampaignItem{
			ID:          middleware.UUIDToString(camp.ID),
			Name:        camp.Name,
			CompanyName: companyNames[cid],
			ProgressPct: 0, // Will be enhanced when v_campaign_progress is queried
		})
	}

	greeting, bgImage := getGreeting()

	data := pages.DashboardPageData{
		UserName:         u.FullName,
		UserRole:         u.Role,
		Greeting:         greeting,
		BgImage:          bgImage,
		PendingReviews:   counts.PendingReviews,
		ContentReview:    counts.ContentReview,
		OverdueTasks:     counts.OverdueTasks,
		RawInbox:         counts.RawInbox,
		OpenTasks:        counts.OpenTasks,
		DoneTasks:        counts.DoneTasks,
		RunningCampaigns: counts.RunningCampaigns,
		Companies:        toTemplCompanies(companies),
		TodayTasks:       toTemplTasks(todayTasks),
		MonthOutput:      outputItems,
		RunCampaigns:     campItems,
	}
	noteContent, _ := h.queries.GetUserNote(c.Context(), u.ID)
	data.NoteContent = noteContent
	return render(c, pages.DashboardDataPage(data))
}

// SaveDashboardNotes handles POST /dashboard/notes (HTMX auto-save)
func (h *Handler) SaveDashboardNotes(c *fiber.Ctx) error {
	user := GetUser(c)
	content := c.FormValue("content")
	_ = h.queries.UpsertUserNote(c.Context(), user.ID, content)
	return c.SendStatus(204)
}

// userDashboard — personal view for CTV and client_staff
func (h *Handler) userDashboard(c *fiber.Ctx) error {
	u := GetUser(c)

	// My tasks
	myTasks, _ := h.queries.ListTasksByAssignee(c.Context(), u.ID)
	myWorkLogs, _ := h.queries.ListWorkLogsByUser(c.Context(), u.ID, 10, 0)
	workTypes, _ := h.queries.ListActiveWorkTypes(c.Context())

	greeting, bgImage := getGreeting()

	data := pages.DashboardPageData{
		UserName:   u.FullName,
		UserRole:   u.Role,
		Greeting:   greeting,
		BgImage:    bgImage,
		TodayTasks: toTemplTasks(myTasks),
		MyWorkLogs: toTemplWorkLogs(myWorkLogs, workTypes),
	}
	return render(c, pages.DashboardDataPage(data))
}
