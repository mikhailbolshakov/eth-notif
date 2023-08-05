package http

import (
	"encoding/json"
	"fmt"
	"github.com/mikhailbolshakov/eth-notif/common"
	"github.com/mikhailbolshakov/eth-notif/domain"
	"net/http"
)

type Controller interface {
	Health(http.ResponseWriter, *http.Request)
	GetCurrentBlock(http.ResponseWriter, *http.Request)
	Subscribe(http.ResponseWriter, *http.Request)
	GetTransactions(http.ResponseWriter, *http.Request)
}

func NewController(bcScanner domain.BlockChainScanner, subs domain.Subscription, trans domain.TransactionProcessor) Controller {
	return &controller{
		bcScanner: bcScanner,
		subs:      subs,
		trans:     trans,
	}
}

type controller struct {
	bcScanner domain.BlockChainScanner
	subs      domain.Subscription
	trans     domain.TransactionProcessor
}

func (c *controller) Health(w http.ResponseWriter, r *http.Request) {
	c.respondJson(w, http.StatusOK, "ok")
}

func (c *controller) GetCurrentBlock(w http.ResponseWriter, r *http.Request) {
	block := c.bcScanner.GetCurrentBlock(r.Context())
	c.respondJson(w, http.StatusOK, common.Int64ToHexStr(block))
}

func (c *controller) Subscribe(w http.ResponseWriter, r *http.Request) {

	address, err := c.queryVar(r, "address", false)
	if err != nil {
		c.respondJson(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := c.subs.Subscribe(r.Context(), address)
	if err != nil {
		c.respondJson(w, http.StatusBadRequest, err.Error())
		return
	}

	if res {
		c.respondJson(w, http.StatusCreated, "created")
	} else {
		c.respondJson(w, http.StatusOK, "already exists")
	}
}

func (c *controller) GetTransactions(w http.ResponseWriter, r *http.Request) {

	address, err := c.queryVar(r, "address", false)
	if err != nil {
		c.respondJson(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := c.trans.GetByAddress(r.Context(), address)
	if err != nil {
		c.respondJson(w, http.StatusBadRequest, err.Error())
		return
	}

	if len(res) == 0 {
		c.respondJson(w, http.StatusNotFound, "not found")
	} else {
		c.respondJson(w, http.StatusOK, res)
	}

}

func (c *controller) respondJson(w http.ResponseWriter, status int, payload any) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(response)
}

func (c *controller) queryVar(r *http.Request, name string, allowEmpty bool) (string, error) {
	val := r.FormValue(name)
	if !allowEmpty && val == "" {
		return "", fmt.Errorf("invalid variable %s", name)
	}
	return val, nil
}
