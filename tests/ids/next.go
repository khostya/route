package ids

import (
	"strconv"
	"sync/atomic"
)

var nextID atomic.Uint64

func NextID() string {
	return strconv.FormatUint(nextID.Add(1), 10)
}
