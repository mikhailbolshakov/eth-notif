package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mikhailbolshakov/eth-notif/common"
	"net/http"
	"time"
)

type Server interface {
	// Listen start listening
	Listen(ctx context.Context)
	// Stop stops listening
	Stop()
	// Router returns router
	Router() *mux.Router
}

func New(port string) Server {
	r := mux.NewRouter()
	s := &server{
		srv: &http.Server{
			Addr:    fmt.Sprintf(":%s", port),
			Handler: r,
		},
		router: r,
	}
	return s
}

type server struct {
	srv    *http.Server
	router *mux.Router
}

func (s *server) Router() *mux.Router {
	return s.router
}

func (s *server) Listen(ctx context.Context) {
	go func() {
	start:
		if err := s.srv.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				common.LogErr(err)
				time.Sleep(time.Second * 5)
				goto start
			} else {
				common.LogDbg("server closed")
			}
			return
		}
	}()
}

func (s *server) Stop() {
	_ = s.srv.Close()
}
