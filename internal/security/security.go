package security

import (
	"errors"
	"net/http"
	"time"
)

var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrSessionNotFound = errors.New("session not found")

type Role int8

const (
	Admin Role = 0
	User  Role = 1
)

type Session struct {
	id string

	UserID   int64
	Username string
}

func (s Session) ToCookie() *http.Cookie {
	return &http.Cookie{
		Name:     sessionCookieName,
		Value:    s.id,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 24),
	}
}

func NewInvalidateCookie() *http.Cookie {
	return &http.Cookie{
		Name:     sessionCookieName,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		Expires:  time.Time{},
	}
}
