package storage

import (
	"context"
	"github.com/mikhailbolshakov/eth-notif/domain"
	"go.uber.org/atomic"
)

func NewScanMemStorage() domain.ScannerStorage {
	return &scanStorage{
		processedBlock: atomic.NewInt64(0),
	}
}

type scanStorage struct {
	processedBlock *atomic.Int64
}

func (s *scanStorage) GetLastProcessedBlock(ctx context.Context) int64 {
	return s.processedBlock.Load()
}

func (s *scanStorage) SetLastProcessedBlock(ctx context.Context, blockNum int64) error {
	s.processedBlock.Store(blockNum)
	return nil
}
