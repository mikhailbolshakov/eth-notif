package domain

import "context"

type BcTransaction struct {
	Block string // Block number
	From  string // From source address
	To    string // To target address
	Value string // Value of transaction
}

type BlockchainRepository interface {
	// GetCurrentBlockNumber retrieves the latest block number
	GetCurrentBlockNumber(ctx context.Context) (int64, error)
	// GetTransByBlock retrieves transactions by the given block number
	GetTransByBlock(ctx context.Context, blockNum int64) ([]*BcTransaction, error)
}
