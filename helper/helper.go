package helper

import (
	"strconv"
)

const MaxMsgChanLen = 1024

func ParseUint64OrPanic(item string) uint64 {
	value, err := strconv.ParseUint(item, 10, 64)
	if err != nil {
		panic(err)
	}

	return value
}
