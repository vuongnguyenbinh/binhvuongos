package handler

import (
	"fmt"

	"binhvuongos/internal/middleware"
	"binhvuongos/web/templates/pages"

	"github.com/gofiber/fiber/v2"
)

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

	var campItems []pages.DashCampaignItem
	for _, camp := range campaigns {
		campItems = append(campItems, pages.DashCampaignItem{
			ID: middleware.UUIDToString(camp.ID), Name: camp.Name,
		})
	}

	data := pages.DashboardPageData{
		UserName:         u.FullName,
		UserRole:         u.Role,
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
	return render(c, pages.DashboardDataPage(data))
}

// userDashboard — personal view for CTV and client_staff
func (h *Handler) userDashboard(c *fiber.Ctx) error {
	u := GetUser(c)

	// My tasks
	myTasks, _ := h.queries.ListTasksByAssignee(c.Context(), u.ID)
	myWorkLogs, _ := h.queries.ListWorkLogsByUser(c.Context(), u.ID, 10, 0)
	workTypes, _ := h.queries.ListActiveWorkTypes(c.Context())

	data := pages.DashboardPageData{
		UserName:   u.FullName,
		UserRole:   u.Role,
		TodayTasks: toTemplTasks(myTasks),
		MyWorkLogs: toTemplWorkLogs(myWorkLogs, workTypes),
	}
	return render(c, pages.DashboardDataPage(data))
}
