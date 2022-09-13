package payloads

import "github.com/ferryvg/hiring-test-go-users-api/internal/domain"

type User struct {
	ID    string   `json:"identity"`
	Roles []string `json:"roles"`
}

func NewUser(user *domain.User) User {
	return User{
		ID:    user.ID,
		Roles: BuildRoles(user.Roles),
	}
}

type UsersList struct {
	Users []User `json:"users"`
}

func BuildRoles(userRoles map[domain.Role]bool) []string {
	roles := make([]string, 0, len(userRoles))

	for role, enabled := range userRoles {
		if enabled {
			roles = append(roles, role.String())
		}
	}

	return roles
}
