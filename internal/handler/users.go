package handler

import (
	"binhvuongos/internal/db/generated"
	"binhvuongos/internal/middleware"
	"binhvuongos/web/templates/pages"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) Users(c *fiber.Ctx) error {
	actor := GetUser(c)
	page, limit, offset := getPage(c)
	users, err := h.queries.ListUsers(c.Context(), limit, offset)
	if err != nil {
		return render(c, pages.UsersListPage(pages.UsersPageData{}))
	}
	total, _ := h.queries.CountUsers(c.Context())

	items := toTemplUsers(users)
	for i, u := range users {
		items[i].CanEdit = middleware.CanManageUser(actor, u)
	}

	data := pages.UsersPageData{
		Users:        items,
		Total:        total,
		Page:         page,
		TotalPages:   totalPages(total),
		ActorRole:    actor.Role,
		AllowedRoles: middleware.AllowedTargetRoles(actor),
	}
	return render(c, pages.UsersListPage(data))
}

func (h *Handler) CreateUser(c *fiber.Ctx) error {
	actor := GetUser(c)
	email := c.FormValue("email")
	fullName := c.FormValue("full_name")
	role := c.FormValue("role")
	password := c.FormValue("password")
	phone := c.FormValue("phone")

	if email == "" || fullName == "" || password == "" {
		return c.Status(400).SendString("Thiếu thông tin bắt buộc")
	}
	if role == "" {
		role = "staff"
	}
	// Server-side whitelist — never trust the form's role value.
	if !middleware.IsAllowedRole(actor, role) {
		return c.Status(403).SendString("Không có quyền tạo user với role này")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return c.Status(500).SendString("Lỗi mã hoá mật khẩu")
	}

	_, err = h.queries.CreateUser(c.Context(), generated.CreateUserParams{
		Email:        email,
		PasswordHash: string(hash),
		FullName:     fullName,
		Role:         role,
		Phone:        phone,
	})
	if err != nil {
		return c.Status(500).SendString("Lỗi tạo user (email trùng?)")
	}
	return c.Redirect("/users")
}

func toTemplUsers(users []generated.User) []pages.UserItem {
	items := make([]pages.UserItem, len(users))
	for i, u := range users {
		items[i] = pages.UserItem{
			ID:       middleware.UUIDToString(u.ID),
			Email:    u.Email,
			FullName: u.FullName,
			Role:     u.Role,
			Status:   u.Status,
			Phone:    nullStr(u.Phone),
		}
	}
	return items
}
