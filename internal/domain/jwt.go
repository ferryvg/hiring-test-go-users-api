package domain

import "time"

type JwtAccessToken struct {
	IDUser      string
	AccessToken string
	CreatedAt   time.Time
	ExpiredAt   time.Time
}
