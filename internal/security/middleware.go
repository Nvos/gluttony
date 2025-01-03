package security

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
)

const sessionCookieName = "GluttonySession"
const sessionID = "GluttonySession"

type ReadOnlySessionStore interface {
	Get(ctx context.Context, sessionID string) (Session, error)
}

func AuthenticationMiddleware(
	logger *slog.Logger,
	store ReadOnlySessionStore,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := resolveSession(r, store)
			if errors.Is(err, ErrSessionNotFound) || errors.Is(err, http.ErrNoCookie) {
				next.ServeHTTP(w, r)
				return
			}

			if err != nil {
				logger.Error(
					"Failed to resolve session",
					slog.String("error", err.Error()),
				)

				next.ServeHTTP(w, r)
				return
			}

			if r.URL.Path == "/login" || r.URL.Path == "/" {
				http.Redirect(w, r, "/recipe", http.StatusTemporaryRedirect)
			}

			nextCtx := context.WithValue(r.Context(), sessionID, session)
			next.ServeHTTP(w, r.WithContext(nextCtx))
		})
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
