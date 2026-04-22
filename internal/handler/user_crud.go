package handler

import (
	"fmt"

	"binhvuongos/internal/db/generated"
	"binhvuongos/internal/middleware"
	"binhvuongos/web/templates/pages"

	"github.com/gofiber/fiber/v2"
)

// EditUserPage renders the user edit form. Only owner/manager with CanManageUser may view.
func (h *Handler) EditUserPage(c *fiber.Ctx) error {
	actor := GetUser(c)
	target, err := h.queries.GetUserByID(c.Context(), middleware.StringToUUID(c.Params("id")))
	if err != nil {
		return c.Status(404).SendString("Không tìm thấy user")
	}
	if !middleware.CanManageUser(actor, target) {
		return c.Status(403).SendString("Không có quyền sửa user này")
	}
	data := pages.UserEditData{
		ID:           middleware.UUIDToString(target.ID),
		Email:        target.Email,
		FullName:     target.FullName,
		Role:         target.Role,
		Phone:        nullStr(target.Phone),
		Status:       target.Status,
		AllowedRoles: middleware.AllowedTargetRoles(actor),
	}
	return render(c, pages.UserEditPage(data))
}

// UpdateUser saves edited user fields with permission + role-whitelist enforcement.
func (h *Handler) UpdateUser(c *fiber.Ctx) error {
	actor := GetUser(c)
	target, err := h.queries.GetUserByID(c.Context(), middleware.StringToUUID(c.Params("id")))
	if err != nil {
		return c.Status(404).SendString("Không tìm thấy user")
	}
	if !middleware.CanManageUser(actor, target) {
		return c.Status(403).SendString("Không có quyền sửa user này")
	}
	newRole := c.FormValue("role")
	if newRole == "" {
		newRole = target.Role
	}
	if !middleware.IsAllowedRole(actor, newRole) {
		return c.Status(403).SendString("Không có quyền gán role này")
	}
	status := c.FormValue("status")
	if status == "" {
		status = target.Status
	}
	_, err = h.queries.UpdateUser(c.Context(), generated.UpdateUserParams{
		ID:       target.ID,
		FullName: c.FormValue("full_name"),
		Role:     newRole,
		Phone:    c.FormValue("phone"),
		Status:   status,
	})
	if err != nil {
		return c.Status(500).SendString("Lỗi cập nhật")
	}
	return c.Redirect("/users")
}

// UserDetail renders the admin-facing user detail page (profile + recent activity).
// Only owner + manager can view; edit button is further gated by CanManageUser.
func (h *Handler) UserDetail(c *fiber.Ctx) error {
	actor := GetUser(c)
	if actor.Role != "owner" && actor.Role != "manager" {
		return c.Status(403).SendString("Forbidden")
	}
	target, err := h.queries.GetUserByID(c.Context(), middleware.StringToUUID(c.Params("id")))
	if err != nil {
		return c.Status(404).SendString("Không tìm thấy user")
	}

	// Recent activity — cap slices to 10 rows even if query returns more.
	tasks, _ := h.queries.ListTasksByAssignee(c.Context(), target.ID)
	if len(tasks) > 10 {
		tasks = tasks[:10]
	}
	logs, _ := h.queries.ListWorkLogsByUser(c.Context(), target.ID, 10, 0)
	assignments, _ := h.queries.ListAssignmentsByUser(c.Context(), target.ID)

	data := pages.UserDetailData{
		ID:        middleware.UUIDToString(target.ID),
		Email:     target.Email,
		FullName:  target.FullName,
		Role:      target.Role,
		Status:    target.Status,
		Phone:     nullStr(target.Phone),
		AvatarURL: nullStr(target.AvatarURL),
		CanEdit:   middleware.CanManageUser(actor, target),
	}
	for _, co := range assignments {
		data.Companies = append(data.Companies, pages.UserDetailCompany{
			Name:      co.CompanyName,
			ShortCode: co.ShortCode,
			Role:      nullStr(co.RoleInCompany),
		})
	}
	for _, t := range tasks {
		data.Tasks = append(data.Tasks, pages.UserDetailTask{
			ID:       middleware.UUIDToString(t.ID),
			Title:    t.Title,
			Status:   t.Status,
			Priority: t.Priority,
			DueDate:  formatDate(t.DueDate),
		})
	}
	for _, l := range logs {
		var qty string
		if l.Quantity.Valid {
			f, _ := l.Quantity.Float64Value()
			qty = fmt.Sprintf("%.0f", f.Float64)
		}
		data.WorkLogs = append(data.WorkLogs, pages.UserDetailWorkLog{
			Date:     formatDate(l.WorkDate),
			Quantity: qty,
			Notes:    nullStr(l.Notes),
		})
	}
	return render(c, pages.UserDetailPage(data))
}

// DeleteUser soft-deletes a user (sets deleted_at). Self-delete is blocked.
func (h *Handler) DeleteUser(c *fiber.Ctx) error {
	actor := GetUser(c)
	target, err := h.queries.GetUserByID(c.Context(), middleware.StringToUUID(c.Params("id")))
	if err != nil {
		return c.Status(404).SendString("Không tìm thấy user")
	}
	if middleware.UUIDToString(actor.ID) == middleware.UUIDToString(target.ID) {
		return c.Status(400).SendString("Không thể tự xoá tài khoản của mình")
	}
	if !middleware.CanManageUser(actor, target) {
		return c.Status(403).SendString("Không có quyền xoá user này")
	}
	if err := h.queries.SoftDeleteUser(c.Context(), target.ID); err != nil {
		return c.Status(500).SendString("Lỗi xoá")
	}
	return c.Redirect("/users")
}
