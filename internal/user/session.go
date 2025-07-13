package user

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"gluttony/internal/config"
	"net/http"
	"sync"
	"time"
)

const SessionCookieName = "GluttonySession"

type Session struct {
	id string

	User User
	Data map[string]any
}

//nolint:ireturn // casts key in session and panics on invalid type
func Get[T any](session Session, key string) (T, bool) {
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
		Name:     SessionCookieName,
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
		Name:     SessionCookieName,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		Expires:  time.Time{},
	}
}

type SessionService struct {
	mu   sync.RWMutex
	data map[string]Session
}

func NewSessionService() *SessionService {
	return &SessionService{
		mu:   sync.RWMutex{},
		data: make(map[string]Session),
	}
}

func (s *SessionService) Create(u User) (Session, error) {
	sessionID, err := generateSessionID()
	if err != nil {
		return Session{}, fmt.Errorf("generate session id: %w", err)
	}

	value := Session{
		id:   sessionID,
		User: u,
		Data: make(map[string]any),
	}

	s.mu.Lock()
	s.data[sessionID] = value
	s.mu.Unlock()

	return value, nil
}

func (s *SessionService) Get(sessionID string) (Session, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, ok := s.data[sessionID]
	if !ok {
		//nolint:exhaustruct // early return when false.
		return Session{}, false
	}

	return value, true
}

func (s *SessionService) Delete(sessionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, sessionID)

	return nil
}

func generateSessionID() (string, error) {
	const sessionIDLength = 32
	b := make([]byte, sessionIDLength)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate session id: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}

func WithContextSession(ctx context.Context, session Session) context.Context {
	return context.WithValue(ctx, sessionKey, session)
}

func GetContextSession(ctx context.Context) (Session, bool) {
	value, ok := ctx.Value(sessionKey).(Session)
	return value, ok
}
