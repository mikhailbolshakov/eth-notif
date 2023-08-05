package domain

import (
	"context"
	"github.com/mikhailbolshakov/eth-notif/common"
)

const (
	TransTypeIn           = "in"
	TransTypeOut          = "out"
	cfgTransProcessorsNum = 4
)

type Transaction struct {
	Address     string
	BlockNumber string
	Type        string
	Value       string
}

type TransactionProcessor interface {
	// Run executes transaction processing in background
	Run(ctx context.Context, transChan TransChannel) error
	// GetByAddress retrieves transactions by address
	GetByAddress(ctx context.Context, address string) ([]*Transaction, error)
}

type TransactionStorage interface {
	// Save saves transactions in storage
	Save(ctx context.Context, address string, tr *Transaction) error
	// GetByAddress retrieves transactions by address
	GetByAddress(ctx context.Context, address string) ([]*Transaction, error)
}

func NewTransProcessor(subscription Subscription, storage TransactionStorage) TransactionProcessor {
	return &transProcessor{
		subscription: subscription,
		storage:      storage,
	}
}

type transProcessor struct {
	subscription Subscription
	storage      TransactionStorage
}

func (t *transProcessor) GetByAddress(ctx context.Context, address string) ([]*Transaction, error) {
	return t.storage.GetByAddress(ctx, address)
}

func (t *transProcessor) Run(ctx context.Context, transChan TransChannel) error {
	// run processing workers
	for i := 0; i < cfgTransProcessorsNum; i++ {
		t.worker(ctx, transChan)
	}
	common.LogDbg("[trans]: processing started with %d workers", cfgTransProcessorsNum)
	return nil
}

func (t *transProcessor) worker(ctx context.Context, transChan TransChannel) {
	go func() {
		for {
			select {
			case trans := <-transChan:
				var err error
				// check subscription and save to storage
				if t.subscription.IsSubscribed(ctx, trans.To) {
					err = t.storage.Save(ctx, trans.To, t.convert(trans.To, trans.Block, TransTypeIn, trans.Value))
				} else if t.subscription.IsSubscribed(ctx, trans.From) {
					err = t.storage.Save(ctx, trans.From, t.convert(trans.From, trans.Block, TransTypeOut, trans.Value))
				}
				if err != nil {
					common.LogErr(err)
					continue
				}
				//common.LogDbg("[trans]: block: %s val: %s", trans.Block, trans.Value)
			case <-ctx.Done():
				common.LogDbg("[trans]: stopped")
				return
			}
		}
	}()
}

func (t *transProcessor) convert(addr, block, tp, value string) *Transaction {
	return &Transaction{
		Address:     addr,
		BlockNumber: block,
		Type:        tp,
		Value:       value,
	}
}
