package handler

import (
	"binhvuongos/internal/db/generated"
	"binhvuongos/internal/middleware"
	"binhvuongos/web/templates/pages"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) Users(c *fiber.Ctx) error {
	page, limit, offset := getPage(c)
	users, err := h.queries.ListUsers(c.Context(), limit, offset)
	if err != nil {
		return render(c, pages.UsersListPage(pages.UsersPageData{}))
	}
	total, _ := h.queries.CountUsers(c.Context())

	data := pages.UsersPageData{
		Users:      toTemplUsers(users),
		Total:      total,
		Page:       page,
		TotalPages: totalPages(total),
	}
	return render(c, pages.UsersListPage(data))
}

func (h *Handler) CreateUser(c *fiber.Ctx) error {
	email := c.FormValue("email")
	fullName := c.FormValue("full_name")
	role := c.FormValue("role")
	password := c.FormValue("password")
	phone := c.FormValue("phone")

	if email == "" || fullName == "" || password == "" {
		return c.Redirect("/users")
	}
	if role == "" {
		role = "ctv"
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return c.Redirect("/users")
	}

	_, _ = h.queries.CreateUser(c.Context(), generated.CreateUserParams{
		Email:        email,
		PasswordHash: string(hash),
		FullName:     fullName,
		Role:         role,
		Phone:        phone,
	})
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
