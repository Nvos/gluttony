package session

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const cookieName = "GluttonySession"

type Key string

var ErrNotFound = errors.New("session not found")

type Store interface {
	Delete(ctx context.Context, sessionID string) error
	Create(ctx context.Context, session Session) error
	Get(ctx context.Context, sessionID string) (Session, error)
}

type Session struct {
	id   string
	Data map[Key]any
}

//nolint:ireturn // casts key in session and panics on invalid type
func Get[T any](session Session, key Key) (T, bool) {
	var value T
	got, ok := session.Data[key]
	if !ok {
		return value, false
	}

	value, ok = got.(T)
	if !ok {
		panic(fmt.Sprintf("cast %+v to invalid type", got))
	}

	return value, true
}

func (s Session) ToCookie(expiresAt time.Time) *http.Cookie {
	return &http.Cookie{
		Name:     cookieName,
		Value:    s.id,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		Expires:  expiresAt,
		// TODO: should be secure in prod
		// TODO: concrete domain in prod
	}
}

func NewInvalidateCookie() *http.Cookie {
	return &http.Cookie{
		Name:     cookieName,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		Expires:  time.Time{},
	}
}
