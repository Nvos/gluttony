package auth

import (
	"context"
	"fmt"
	"sync"
)

type MemoryStorage[T any] struct {
	mu   sync.RWMutex
	data map[string]T
}

func (s *MemoryStorage[T]) Single(_ context.Context, key string) (T, error) {
	s.mu.RLock()
	defer func() {
		s.mu.RUnlock()
	}()

	t, ok := s.data[key]
	if !ok {
		var t T

		return t, fmt.Errorf("auth: memory storage not found key=%s", key)
	}

	return t, nil
}

func (s *MemoryStorage[T]) Create(_ context.Context, key string, value T) error {
	s.mu.Lock()
	defer func() {
		s.mu.Unlock()
	}()

	s.data[key] = value

	return nil
}

func (s *MemoryStorage[T]) Delete(_ context.Context, key string) error {
	s.mu.Lock()
	defer func() {
		s.mu.Unlock()
	}()

	delete(s.data, key)

	return nil
}

func NewMemoryStorage[T any]() *MemoryStorage[T] {
	return &MemoryStorage[T]{
		mu:   sync.RWMutex{},
		data: map[string]T{},
	}
}
