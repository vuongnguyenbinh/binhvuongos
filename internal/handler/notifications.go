package handler

import (
	"fmt"

	"binhvuongos/internal/middleware"
	"binhvuongos/web/templates/pages"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Notifications(c *fiber.Ctx) error {
	user := GetUser(c)
	notifs, _ := h.queries.ListNotifications(c.Context(), user.ID, 30)
	var items []pages.NotifItem
	for _, n := range notifs {
		link := ""
		if n.Link != nil {
			link = *n.Link
		}
		items = append(items, pages.NotifItem{
			ID:    middleware.UUIDToString(n.ID),
			Title: n.Title,
			Link:  link,
			IsRead: n.ReadAt.Valid,
			TimeAgo: timeAgo(n.CreatedAt),
		})
	}
	return render(c, pages.NotificationsPage(items))
}

func (h *Handler) NotificationCount(c *fiber.Ctx) error {
	user := GetUser(c)
	count, _ := h.queries.CountUnreadNotifications(c.Context(), user.ID)
	if count == 0 {
		return c.SendString("")
	}
	return c.SendString(fmt.Sprintf(`<span class="absolute -top-1 -right-1 bg-ember text-white text-[9px] mono rounded-full w-4 h-4 flex items-center justify-center">%d</span>`, count))
}

func (h *Handler) MarkNotificationRead(c *fiber.Ctx) error {
	id := c.Params("id")
	_ = h.queries.MarkNotificationRead(c.Context(), middleware.StringToUUID(id))
	return c.SendString("")
}

func (h *Handler) MarkAllRead(c *fiber.Ctx) error {
	user := GetUser(c)
	_ = h.queries.MarkAllNotificationsRead(c.Context(), user.ID)
	return c.Redirect("/notifications")
}
