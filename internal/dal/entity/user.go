package entity

type User struct {
	ID string `db:"id_user"`

	Secret string `db:"secret,omitempty"`
}
