package eth

import (
	"context"
	"github.com/mikhailbolshakov/eth-notif/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_GetBlockNumber(t *testing.T) {
	num, err := NewEthBcRepository().GetCurrentBlockNumber(context.Background())
	assert.NoError(t, err)
	assert.NotZero(t, num)
}

func Test_GetTransactions(t *testing.T) {
	b, _ := common.HexStrToInt64("0x1105315")
	v, err := NewEthBcRepository().GetTransByBlock(context.Background(), b)
	assert.NoError(t, err)
	assert.NotEmpty(t, v)
}
