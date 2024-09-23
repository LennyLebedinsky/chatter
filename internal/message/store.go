package message

import (
	"context"
	"sync"
)

// Store is intended to provide message retention.
// Message history is stored in chronological order by rooms.
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

	// TODO: implement eviction to control memory usage.
	// Simplest implementation would require using linked list of messages instead of array
	// to maintain constant maximum size of history for each room and control it in O(1).
	if _, ok := s.history[roomName]; !ok {
		s.history[roomName] = []*Message{msg}
	} else {
		s.history[roomName] = append(s.history[roomName], msg)
	}
	return nil
}
