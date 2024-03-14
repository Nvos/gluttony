package auth

import (
	"net/http"
	"time"
)

const cookieName = "session"

type SessionCookieOption func(*http.Cookie)

func WithExpiresAt(expiry time.Time) SessionCookieOption {
	return func(cookie *http.Cookie) {
		cookie.Expires = time.Unix(expiry.Unix()+1, 0)
		cookie.MaxAge = int(time.Until(expiry).Seconds() + 1)
	}
}

func WithExpired() SessionCookieOption {
	return func(cookie *http.Cookie) {
		cookie.Expires = time.Unix(1, 0)
		cookie.MaxAge = -1
	}
}

func NewUnsecureSessionCookie(token string, opts ...SessionCookieOption) *http.Cookie {
	cookie := &http.Cookie{
		Name:     cookieName,
		Value:    token,
		Path:     "/",
		Domain:   "",
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	for i := range opts {
		opts[i](cookie)
	}

	return cookie
}

func NewExpiredSessionCookie(opts ...SessionCookieOption) *http.Cookie {
	cookie := &http.Cookie{
		Name:     cookieName,
		Value:    "",
		SameSite: http.SameSiteLaxMode,
		HttpOnly: false,
		Path:     "/",
		MaxAge:   -1,
	}

	for i := range opts {
		opts[i](cookie)
	}

	return cookie
}
