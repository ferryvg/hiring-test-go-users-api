package domain

type User struct {
	ID string

	Secret string

	Roles map[Role]bool
}

func NewGuestUser() *User {
	return &User{
		Roles: map[Role]bool{
			GuestRole: true,
			BasicRole: false,
			AdminRole: false,
		},
	}
}
