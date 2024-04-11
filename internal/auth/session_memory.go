package auth

import (
	"context"
	"fmt"
	"sync"
)

type SessionMemory struct {
	mu   sync.RWMutex
	data map[string]Session
}

func (s *SessionMemory) Single(_ context.Context, key string) (Session, error) {
	s.mu.RLock()
	defer func() {
		s.mu.RUnlock()
	}()

	t, ok := s.data[key]
	if !ok {
		return Session{}, fmt.Errorf("auth: memory storage not found key=%s", key)
	}

	return t, nil
}

func (s *SessionMemory) Create(_ context.Context, key string, value Session) error {
	s.mu.Lock()
	defer func() {
		s.mu.Unlock()
	}()

	s.data[key] = value

	return nil
}

func (s *SessionMemory) Delete(_ context.Context, key string) error {
	s.mu.Lock()
	defer func() {
		s.mu.Unlock()
	}()

	delete(s.data, key)

	return nil
}

func NewMemoryStorage() *SessionMemory {
	return &SessionMemory{
		mu:   sync.RWMutex{},
		data: map[string]Session{},
	}
}
