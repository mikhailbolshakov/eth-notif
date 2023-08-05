package bootstrap

import (
	"context"
	"github.com/mikhailbolshakov/eth-notif/common"
	"github.com/mikhailbolshakov/eth-notif/domain"
	"github.com/mikhailbolshakov/eth-notif/http"
	"github.com/mikhailbolshakov/eth-notif/repository/eth"
	"github.com/mikhailbolshakov/eth-notif/repository/storage"
	"os"
	"strings"
)

type Service interface {
	// Run all the background processes
	Run(ctx context.Context) error
	// Close all the processes
	Close()
}

type service struct {
	bcScanner domain.BlockChainScanner
	subs      domain.Subscription
	transProc domain.TransactionProcessor
	httpSrv   http.Server
}

func New() Service {
	s := &service{}

	transStorage := storage.NewTransactionMemStorage()
	scanStorage := storage.NewScanMemStorage()
	s.subs = domain.NewSubscription(storage.NewSubscriptionMemStorage())
	s.bcScanner = domain.NewScanner(eth.NewEthBcRepository(), scanStorage)
	s.transProc = domain.NewTransProcessor(s.subs, transStorage)

	s.httpSrv = http.New(os.Getenv("PORT"))
	http.NewRouter(s.httpSrv.Router(), http.NewController(s.bcScanner, s.subs, s.transProc)).SetRoutes()

	ctx := context.Background()
	s.loadDefaultSubscribed(ctx)
	s.setupLastProcessedBlock(ctx, scanStorage)

	return s
}

func (s *service) Run(ctx context.Context) error {

	// start scanner
	transChan, err := s.bcScanner.Run(ctx)
	if err != nil {
		return err
	}

	// start
	err = s.transProc.Run(ctx, transChan)
	if err != nil {
		return err
	}

	// start http server
	s.httpSrv.Listen(ctx)

	return nil
}

func (s *service) Close() {
	s.httpSrv.Stop()
}

// loadDefaultSubscribed allows to set up predefined subscribed address (for test purposes)
func (s *service) loadDefaultSubscribed(ctx context.Context) {
	v := os.Getenv("DEFAULT_SUBSCRIBED")
	if v == "" {
		return
	}
	subscribed := strings.Split(v, ",")
	for _, sub := range subscribed {
		_, _ = s.subs.Subscribe(ctx, sub)
	}
}

// setupLastProcessedBlock allows to start processing from the block in the past (for test purposes)
func (s *service) setupLastProcessedBlock(ctx context.Context, scanStorage domain.ScannerStorage) {
	// start block
	v := os.Getenv("START_BLOCK")
	if v != "" {
		b, _ := common.HexStrToInt64(v)
		_ = scanStorage.SetLastProcessedBlock(ctx, b)
	}
}
