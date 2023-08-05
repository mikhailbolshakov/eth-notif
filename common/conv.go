package common

import (
	"fmt"
	"strconv"
	"strings"
)

func HexStrToInt64(s string) (int64, error) {
	if s == "" {
		return 0, fmt.Errorf("invalid string")
	}
	return strconv.ParseInt(strings.Replace(s, "0x", "", -1), 16, 64)
}

func Int64ToHexStr(v int64) string {
	return fmt.Sprintf("0x%x", v)
}
