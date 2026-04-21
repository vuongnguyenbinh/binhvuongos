package handler

import (
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
