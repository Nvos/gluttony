package security

import (
	"context"
	"errors"
	"fmt"
	"gluttony/x/httpx"
	"net/http"
)

const sessionCookieName = "GluttonySession"
const sessionID = "GluttonySession"

type ReadOnlySessionStore interface {
	Get(ctx context.Context, sessionID string) (Session, error)
}

func NewAuthenticationMiddleware(store ReadOnlySessionStore) httpx.MiddlewareFunc {
	return func(next httpx.HandlerFunc) httpx.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) error {
			session, err := resolveSession(r, store)
			if err != nil {
				return next(w, r)
			}

			if r.URL.Path == "/login" || r.URL.Path == "/" {
				http.Redirect(w, r, "/recipe", http.StatusTemporaryRedirect)
			}

			nextCtx := context.WithValue(r.Context(), sessionID, session)
			return next(w, r.WithContext(nextCtx))
		}
	}
}

func resolveSession(
	r *http.Request,
	store ReadOnlySessionStore,
) (Session, error) {
	c, err := r.Cookie(sessionCookieName)
	if errors.Is(err, http.ErrNoCookie) {
		return Session{}, err
	}

	if err != nil {
		panic(fmt.Sprintf("unexpected cookie err: %v", err))
	}

	session, err := store.Get(r.Context(), c.Value)
	if err != nil {
		return Session{}, err
	}

	return session, err
}

func GetSession(ctx context.Context) (Session, bool) {
	ctxSession, ok := ctx.Value(sessionID).(Session)
	if !ok {
		return Session{}, false
	}

	return ctxSession, true
}
