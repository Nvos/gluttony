package session

import (
	"context"
	"errors"
	"fmt"
	"gluttony/internal/config"
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

func (s Session) ToCookie(cfg *config.Config) *http.Cookie {
	const monthDuration = 30 * 24 * time.Hour // month

	cookie := &http.Cookie{
		Name:     cookieName,
		Value:    s.id,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		Expires:  time.Now().UTC().Add(monthDuration),
	}

	if cfg != nil && cfg.Mode == config.ModeProd {
		cookie.Secure = true
		cookie.SameSite = http.SameSiteStrictMode
		cookie.Domain = cfg.Domain
	}

	return cookie
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
