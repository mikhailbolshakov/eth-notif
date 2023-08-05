package domain

import (
	"context"
	"github.com/mikhailbolshakov/eth-notif/common"
	"time"
)

const (
	cfgBcScanPeriod                 = time.Second * 3
	cfgTransactionChannelBufferSize = 4
	cfgBcReadersNum                 = 4
)

type TransChannel chan *BcTransaction

type BlockChainScanner interface {
	// Run executes blockchain scanner in background and push obtained transactions to the output channel
	Run(ctx context.Context) (TransChannel, error)
	// GetCurrentBlock retrieves the latest scanned block
	GetCurrentBlock(ctx context.Context) int64
}

type ScannerStorage interface {
	// GetLastProcessedBlock retrieves the latest processed block
	GetLastProcessedBlock(ctx context.Context) int64
	// SetLastProcessedBlock sets last processed block
	SetLastProcessedBlock(ctx context.Context, blockNum int64) error
}

type scanner struct {
	bcRep   BlockchainRepository
	storage ScannerStorage
}

func NewScanner(bcRep BlockchainRepository, storage ScannerStorage) BlockChainScanner {
	return &scanner{
		bcRep:   bcRep,
		storage: storage,
	}
}

func (s *scanner) GetCurrentBlock(ctx context.Context) int64 {
	return s.storage.GetLastProcessedBlock(ctx)
}

func (s *scanner) Run(ctx context.Context) (TransChannel, error) {

	// first running, start from the current
	if s.storage.GetLastProcessedBlock(ctx) <= 0 {

		// get current from blockchain
		startWithBlock, err := s.bcRep.GetCurrentBlockNumber(ctx)
		if err != nil {
			return nil, err
		}

		// save the latest
		err = s.storage.SetLastProcessedBlock(ctx, startWithBlock)
		if err != nil {
			return nil, err
		}
	}

	// create output channel
	transChan := make(TransChannel, cfgTransactionChannelBufferSize)
	blockNumChan := make(chan int64, cfgBcReadersNum*2)

	// block obtain workers
	for i := 0; i < cfgBcReadersNum; i++ {
		s.transactionReadWorker(ctx, blockNumChan, transChan)
	}

	// run ticker
	go func() {
		ticker := time.NewTicker(cfgBcScanPeriod)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				s.scan(ctx, blockNumChan)
			case <-ctx.Done():
				common.LogDbg("[scanner] ticker stopped")
				return
			}
		}
	}()

	return transChan, nil
}

func (s *scanner) scan(ctx context.Context, blockNumChan chan<- int64) {

	// get current block
	curBlock, err := s.bcRep.GetCurrentBlockNumber(ctx)
	if err != nil {
		common.LogErr(err)
		return
	}

	// get last processed
	lastBlock := s.storage.GetLastProcessedBlock(ctx)

	// no block changed
	if lastBlock >= curBlock {
		return
	}

	// go through range of blocks and send for processing
	for i := lastBlock + 1; i <= curBlock; i++ {
		blockNumChan <- i
	}

	// save the latest
	err = s.storage.SetLastProcessedBlock(ctx, curBlock)
	if err != nil {
		common.LogErr(err)
	}

	common.LogDbg("[scanner] scanned blocks: %s - %s", common.Int64ToHexStr(lastBlock), common.Int64ToHexStr(curBlock))

}

func (s *scanner) transactionReadWorker(ctx context.Context, blockNumChan <-chan int64, transChan TransChannel) {
	go func() {
		for {
			select {
			case num := <-blockNumChan:
				trans, err := s.bcRep.GetTransByBlock(ctx, num)
				if err != nil {
					common.LogErr(err)
					continue
				}
				for _, tr := range trans {
					transChan <- tr
				}
				common.LogDbg("[scanner]: block scanned, found %d trans: %d", num, len(trans))
			case <-ctx.Done():
				common.LogDbg("[scanner] worker stopped")
				return
			}
		}
	}()
}
