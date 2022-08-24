package main

import (
	"math/rand"
	"strconv"
	"time"
)

func newLastAccess() string {
	return strconv.FormatInt(time.Now().UnixNano(), 36) + strconv.FormatUint(rand.Uint64(), 36)
}
