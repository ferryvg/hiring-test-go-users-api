package entity

import "time"

type JwtAccessToken struct {
	IDUser      string    `db:"id_user"`
	AccessToken string    `db:"access_token"`
	CreatedAt   time.Time `db:"created_at"`
	ExpiredAt   time.Time `db:"expired_at"`
}
