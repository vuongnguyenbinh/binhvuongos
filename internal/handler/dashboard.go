package handler

import (
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
	}
	return render(c, pages.DashboardDataPage(data))
}
