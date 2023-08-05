package domain

import (
	"context"
	"github.com/mikhailbolshakov/eth-notif/common"
)

type Subscription interface {
	// Subscribe address
	Subscribe(ctx context.Context, address string) (bool, error)
	// IsSubscribed returns true if address is subscribed
	IsSubscribed(ctx context.Context, address string) bool
}

type SubscriptionStorage interface {
	// CreateOrUpdate subscription
	CreateOrUpdate(ctx context.Context, address string) error
	// Exists checks if address subscription exists
	Exists(ctx context.Context, address string) bool
}

type subscription struct {
	storage SubscriptionStorage
}

func NewSubscription(storage SubscriptionStorage) Subscription {
	return &subscription{
		storage: storage,
	}
}

func (s *subscription) Subscribe(ctx context.Context, address string) (bool, error) {
	// check if already exists
	if s.storage.Exists(ctx, address) {
		return false, nil
	}

	// subscribe
	err := s.storage.CreateOrUpdate(ctx, address)
	if err != nil {
		return false, err
	}

	common.LogDbg("[subs] %s subscribed", address)

	return true, nil
}

func (s *subscription) IsSubscribed(ctx context.Context, address string) bool {
	return s.storage.Exists(ctx, address)
}
