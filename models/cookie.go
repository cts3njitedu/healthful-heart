package models

import (
	"time"
)


type JwtCookie struct {
	Name string
	Value string
	Expires time.Time
	HttpOnly bool
}