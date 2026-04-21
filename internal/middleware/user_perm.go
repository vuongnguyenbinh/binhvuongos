package middleware

import "binhvuongos/internal/db/generated"

// CanManageUser returns true if the actor is allowed to view/edit/delete the target user.
// Owner manages everyone; manager only manages staff; staff has no admin power.
func CanManageUser(actor, target generated.User) bool {
	switch actor.Role {
	case "owner":
		return true
	case "manager":
		return target.Role == "staff"
	default:
		return false
	}
}

// AllowedTargetRoles returns the role values an actor is allowed to assign when
// creating or updating users. Used to enforce server-side whitelist on form input.
func AllowedTargetRoles(actor generated.User) []string {
	switch actor.Role {
	case "owner":
		return []string{"owner", "manager", "staff"}
	case "manager":
		return []string{"staff"}
	default:
		return nil
	}
}

// IsAllowedRole checks membership in AllowedTargetRoles without importing slices.
func IsAllowedRole(actor generated.User, role string) bool {
	for _, r := range AllowedTargetRoles(actor) {
		if r == role {
			return true
		}
	}
	return false
}
