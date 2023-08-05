package eth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/mikhailbolshakov/eth-notif/common"
	"github.com/mikhailbolshakov/eth-notif/domain"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	ethTimeout = time.Second * 10
	ethUrl     = "https://cloudflare-eth.com"
)

type eth struct{}

func NewEthBcRepository() domain.BlockchainRepository {
	return &eth{}
}

func (e *eth) GetCurrentBlockNumber(ctx context.Context) (int64, error) {

	rq := &request{
		JsonRpc: "2.0",
		Method:  "eth_blockNumber",
		Id:      "83",
	}
	rs, err := e.makeRq(ctx, rq)
	if err != nil {
		return 0, e.ethErr(fmt.Errorf("[block-num] invalid response"))
	}

	v := ""
	err = json.Unmarshal(rs["result"], &v)
	if err != nil {
		return 0, e.ethErr(fmt.Errorf("[block-num] invalid response"))
	}

	blockNum, err := common.HexStrToInt64(v)
	if err != nil {
		return 0, e.ethErr(fmt.Errorf("[block-num] invalid format"))
	}

	return blockNum, nil
}

func (e *eth) GetTransByBlock(ctx context.Context, blockNum int64) ([]*domain.BcTransaction, error) {

	rq := &request{
		JsonRpc: "2.0",
		Method:  "eth_getBlockByNumber",
		Id:      "83",
		Params:  []any{common.Int64ToHexStr(blockNum), true},
	}
	rs, err := e.makeRq(ctx, rq)
	if err != nil {
		return nil, e.ethErr(fmt.Errorf("[get-block] invalid response"))
	}

	b := &block{}
	err = json.Unmarshal(rs["result"], &b)
	if err != nil {
		return nil, e.ethErr(fmt.Errorf("[get-block] invalid response"))
	}

	if b == nil || len(b.Transactions) == 0 {
		return nil, nil
	}

	trans := make([]*domain.BcTransaction, 0, len(b.Transactions))
	for _, t := range b.Transactions {
		trans = append(trans, &domain.BcTransaction{
			Block: t.BlockNumber,
			From:  t.From,
			To:    t.To,
			Value: t.Value,
		})
	}

	return trans, nil
}

func (e *eth) makeRq(ctx context.Context, rq *request) (map[string]json.RawMessage, error) {

	// setup timeout
	ctxExec, cancelFn := context.WithTimeout(ctx, ethTimeout)
	defer cancelFn()

	// payload
	var rqReader io.Reader
	if rq != nil {
		bodyB, err := json.Marshal(rq)
		if err != nil {
			return nil, e.ethErr(err)
		}
		rqReader = bytes.NewReader(bodyB)
	}

	// prepare request
	httpRq, err := http.NewRequestWithContext(ctxExec, http.MethodPost, ethUrl, rqReader)
	if err != nil {
		return nil, e.ethErr(err)
	}

	// make request
	resp, err := http.DefaultClient.Do(httpRq)
	if err != nil {
		return nil, e.ethErr(err)
	}

	// check response
	if resp == nil || resp.Body == nil {
		return nil, e.ethErr(fmt.Errorf("empty response"))
	}

	defer func() { _ = resp.Body.Close() }()

	// parse body
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, e.ethErr(err)
	}

	// check response
	var rs map[string]json.RawMessage
	err = json.Unmarshal(data, &rs)
	if err != nil {
		return nil, e.ethErr(err)
	}
	if rs == nil {
		return nil, e.ethErr(fmt.Errorf("empty response object"))
	}

	return rs, nil
}

func (e *eth) ethErr(cause error) error {
	return fmt.Errorf("[eth] %s", cause.Error())
}
