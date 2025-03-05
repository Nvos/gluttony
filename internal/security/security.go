package security

import (
	"context"
	"errors"
	"net/http"
	"time"
)

const sessionCookieName = "GluttonySession"
const sessionID = "GluttonySession"

var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrSessionNotFound = errors.New("session not found")

type Role int8

const (
	Admin Role = 0
	User  Role = 1
)

type Session struct {
	id string

	UserID   int32
	Username string
	Role     Role
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

func SetSession(ctx context.Context, session Session) context.Context {
	return context.WithValue(ctx, sessionCookieName, session)
}

func GetSession(ctx context.Context) (Session, bool) {
	ctxSession, ok := ctx.Value(sessionID).(Session)
	if !ok {
		return Session{}, false
	}

	return ctxSession, true
}
