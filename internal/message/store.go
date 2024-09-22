package message

import (
	"context"
	"sync"
)

type Store interface {
	GetMessages(ctx context.Context, roomName string) ([]*Message, error)
	SaveMessage(ctx context.Context, roomName string, msg *Message) error
}

type InMemoryStore struct {
	history map[string][]*Message
	mu      sync.RWMutex
}

func NewInMemoryStore() Store {
	return &InMemoryStore{
		history: make(map[string][]*Message),
	}
}

func (s *InMemoryStore) GetMessages(ctx context.Context, roomName string) ([]*Message, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.history[roomName], nil
}

func (s *InMemoryStore) SaveMessage(ctx context.Context, roomName string, msg *Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.history[roomName]; !ok {
		s.history[roomName] = []*Message{msg}
	} else {
		s.history[roomName] = append(s.history[roomName], msg)
	}
	return nil
}
