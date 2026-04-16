package generated

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type AssignmentWithUser struct {
	UserCompanyAssignment
	FullName   string `json:"full_name"`
	Email      string `json:"email"`
	UserRole   string `json:"user_role"`
	UserStatus string `json:"user_status"`
}

func (q *Queries) ListAssignmentsByCompany(ctx context.Context, companyID pgtype.UUID) ([]AssignmentWithUser, error) {
	rows, err := q.pool.Query(ctx,
		`SELECT uca.id, uca.user_id, uca.company_id, uca.role_in_company, uca.can_view, uca.can_edit, uca.can_approve,
		 uca.start_date, uca.end_date, uca.notes, uca.created_at, uca.updated_at,
		 u.full_name, u.email, u.role, u.status
		 FROM user_company_assignments uca JOIN users u ON u.id = uca.user_id
		 WHERE uca.company_id = $1 AND (uca.end_date IS NULL OR uca.end_date > CURRENT_DATE)
		 ORDER BY u.full_name`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []AssignmentWithUser{}
	for rows.Next() {
		var a AssignmentWithUser
		if err := rows.Scan(&a.ID, &a.UserID, &a.CompanyID, &a.RoleInCompany, &a.CanView, &a.CanEdit,
			&a.CanApprove, &a.StartDate, &a.EndDate, &a.Notes, &a.CreatedAt, &a.UpdatedAt,
			&a.FullName, &a.Email, &a.UserRole, &a.UserStatus); err != nil {
			return nil, err
		}
		items = append(items, a)
	}
	return items, rows.Err()
}

type AssignmentWithCompany struct {
	UserCompanyAssignment
	CompanyName string `json:"company_name"`
	ShortCode   string `json:"short_code"`
}

func (q *Queries) ListAssignmentsByUser(ctx context.Context, userID pgtype.UUID) ([]AssignmentWithCompany, error) {
	rows, err := q.pool.Query(ctx,
		`SELECT uca.id, uca.user_id, uca.company_id, uca.role_in_company, uca.can_view, uca.can_edit, uca.can_approve,
		 uca.start_date, uca.end_date, uca.notes, uca.created_at, uca.updated_at,
		 c.name, COALESCE(c.short_code, '')
		 FROM user_company_assignments uca JOIN companies c ON c.id = uca.company_id
		 WHERE uca.user_id = $1 AND (uca.end_date IS NULL OR uca.end_date > CURRENT_DATE)
		 ORDER BY c.name`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []AssignmentWithCompany{}
	for rows.Next() {
		var a AssignmentWithCompany
		if err := rows.Scan(&a.ID, &a.UserID, &a.CompanyID, &a.RoleInCompany, &a.CanView, &a.CanEdit,
			&a.CanApprove, &a.StartDate, &a.EndDate, &a.Notes, &a.CreatedAt, &a.UpdatedAt,
			&a.CompanyName, &a.ShortCode); err != nil {
			return nil, err
		}
		items = append(items, a)
	}
	return items, rows.Err()
}

type CreateAssignmentParams struct {
	UserID        pgtype.UUID `json:"user_id"`
	CompanyID     pgtype.UUID `json:"company_id"`
	RoleInCompany string      `json:"role_in_company"`
	CanView       bool        `json:"can_view"`
	CanEdit       bool        `json:"can_edit"`
	CanApprove    bool        `json:"can_approve"`
}

func (q *Queries) CreateAssignment(ctx context.Context, arg CreateAssignmentParams) (UserCompanyAssignment, error) {
	var a UserCompanyAssignment
	err := q.pool.QueryRow(ctx,
		`INSERT INTO user_company_assignments (user_id, company_id, role_in_company, can_view, can_edit, can_approve)
		 VALUES ($1,$2,$3,$4,$5,$6) RETURNING id, user_id, company_id, role_in_company, can_view, can_edit, can_approve, start_date, end_date, notes, created_at, updated_at`,
		arg.UserID, arg.CompanyID, arg.RoleInCompany, arg.CanView, arg.CanEdit, arg.CanApprove).
		Scan(&a.ID, &a.UserID, &a.CompanyID, &a.RoleInCompany, &a.CanView, &a.CanEdit, &a.CanApprove, &a.StartDate, &a.EndDate, &a.Notes, &a.CreatedAt, &a.UpdatedAt)
	return a, err
}

func (q *Queries) DeleteAssignment(ctx context.Context, id pgtype.UUID) error {
	return q.exec(ctx, "DELETE FROM user_company_assignments WHERE id=$1", id)
}
