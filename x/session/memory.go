package session

import (
	"context"
	"sync"
)

type StoreMemory struct {
	mu   sync.RWMutex
	data map[string]Session
}

func NewStoreMemory() *StoreMemory {
	return &StoreMemory{
		mu:   sync.RWMutex{},
		data: make(map[string]Session),
	}
}

func (s *StoreMemory) Delete(_ context.Context, sessionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, sessionID)

	return nil
}

func (s *StoreMemory) Create(_ context.Context, value Session) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[value.id] = value

	return nil
}

func (s *StoreMemory) Get(_ context.Context, sessionID string) (Session, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, ok := s.data[sessionID]
	if !ok {
		return Session{}, ErrNotFound
	}

	return data, nil
}
