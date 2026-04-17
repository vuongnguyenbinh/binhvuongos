package handler

import (
	"fmt"

	"binhvuongos/internal/middleware"
	"binhvuongos/web/templates/pages"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Dashboard(c *fiber.Ctx) error {
	user := GetUser(c)
	counts, err := h.queries.GetDashboardCounts(c.Context())
	if err != nil {
		return render(c, pages.DashboardPage())
	}

	companies, _ := h.queries.ListCompanies(c.Context(), 10, 0)
	todayTasks, _ := h.queries.ListTasksDueToday(c.Context())
	monthOutput, _ := h.queries.GetDashboardOutputThisMonth(c.Context())
	campaigns, _ := h.queries.ListCampaignsByStatus(c.Context(), "running")

	// Convert output to view
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
			Name:  o.Name,
			Icon:  o.Icon,
			Unit:  o.Unit,
			Total: total,
		})
	}

	// Convert campaigns
	var campItems []pages.DashCampaignItem
	for _, camp := range campaigns {
		campItems = append(campItems, pages.DashCampaignItem{
			ID:   middleware.UUIDToString(camp.ID),
			Name: camp.Name,
		})
	}

	data := pages.DashboardPageData{
		UserName:         user.FullName,
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
