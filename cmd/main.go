package main

import (
	"context"
	"github.com/mikhailbolshakov/eth-notif/bootstrap"
	"github.com/mikhailbolshakov/eth-notif/common"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	// init context
	ctx, cancelFn := context.WithCancel(context.Background())

	// create a new service
	s := bootstrap.New()

	// start listening
	if err := s.Run(ctx); err != nil {
		common.LogErr(err)
		os.Exit(1)
	}

	common.LogDbg("listening")

	// handle app close
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	cancelFn()
	s.Close()
	common.LogDbg("closed")
	os.Exit(0)
}
