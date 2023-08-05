package storage

import (
	"context"
	"github.com/mikhailbolshakov/eth-notif/domain"
	"sync"
)

func NewTransactionMemStorage() domain.TransactionStorage {
	return &transStorage{
		data: make(map[string][]*domain.Transaction),
	}
}

type transStorage struct {
	sync.RWMutex
	data map[string][]*domain.Transaction
}

func (t *transStorage) Save(ctx context.Context, address string, tr *domain.Transaction) error {
	t.Lock()
	defer t.Unlock()
	t.data[address] = append(t.data[address], tr)
	return nil
}

func (t *transStorage) GetByAddress(ctx context.Context, address string) ([]*domain.Transaction, error) {
	t.RLock()
	defer t.RUnlock()
	return t.data[address], nil
}
