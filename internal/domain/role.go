package domain

type Role int

const (
	// GuestRole is a guest user role with access to non-auth operations.
	GuestRole Role = iota

	// BasicRole is a common user role with access to general operations without extra permissions.
	BasicRole

	// AdminRole is a user role with access to extended operations.
	AdminRole

	UnknownRole = 99998
)

func (r Role) String() string {
	switch r {
	case GuestRole:
		return "guest"

	case BasicRole:
		return "basic"

	case AdminRole:
		return "admin"
	}

	return "unknown"
}

func RoleFromString(role string) Role {
	switch role {
	case "guest":
		return GuestRole

	case "basic":
		return BasicRole

	case "admin":
		return AdminRole
	}

	return UnknownRole
}
