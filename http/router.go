package http

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Router interface {
	SetRoutes()
}

func NewRouter(r *mux.Router, controller Controller) Router {
	return &router{
		routes:     r,
		controller: controller,
	}
}

type router struct {
	routes     *mux.Router
	controller Controller
}

func (r *router) SetRoutes() {
	r.routes.HandleFunc(`/health`, r.controller.Health).Methods(http.MethodGet)
	r.routes.HandleFunc(`/subscriptions`, r.controller.Subscribe).Methods(http.MethodPost)
	r.routes.HandleFunc(`/transactions`, r.controller.GetTransactions).Methods(http.MethodGet)
	r.routes.HandleFunc(`/blocks/current`, r.controller.GetCurrentBlock).Methods(http.MethodGet)
}
