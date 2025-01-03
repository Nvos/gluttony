package security

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"sync"
)

type SessionStore struct {
	mu   sync.RWMutex
	data map[string]Session
}

func NewSessionStore() *SessionStore {
	return &SessionStore{
		data: map[string]Session{},
	}
}

func (s *SessionStore) Save(_ context.Context, session Session) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[session.id] = session

	return nil
}

func (s *SessionStore) New(_ context.Context) (Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	sessionID, err := generateRandomBytes(32)
	if err != nil {
		return Session{}, fmt.Errorf("generate session id: %w", err)
	}

	session := Session{
		id: sessionID,
	}
	s.data[sessionID] = session

	return session, nil
}

func (s *SessionStore) Get(_ context.Context, sessionID string) (Session, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, ok := s.data[sessionID]
	if !ok {
		return Session{}, ErrSessionNotFound
	}

	return data, nil
}

func (s *SessionStore) Delete(_ context.Context, session Session) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.data[session.id]
	if !ok {
		return ErrSessionNotFound
	}

	delete(s.data, session.id)

	return nil
}

func generateRandomBytes(n uint32) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(b), nil
}
