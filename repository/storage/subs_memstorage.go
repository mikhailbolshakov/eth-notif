package storage

import (
	"context"
	"github.com/mikhailbolshakov/eth-notif/domain"
	"sync"
)

func NewSubscriptionMemStorage() domain.SubscriptionStorage {
	return &subsStorage{
		data: make(map[string]struct{}),
	}
}

type subsStorage struct {
	sync.RWMutex
	data map[string]struct{}
}

func (s *subsStorage) CreateOrUpdate(ctx context.Context, address string) error {
	s.Lock()
	defer s.Unlock()
	s.data[address] = struct{}{}
	return nil
}

func (s *subsStorage) Exists(ctx context.Context, address string) bool {
	s.RLock()
	defer s.RUnlock()
	_, ok := s.data[address]
	return ok
}
