package ids

import (
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

func TestNextNumber(t *testing.T) {
	for i := 1; i < 10; i++ {
		require.Equal(t, strconv.FormatUint(uint64(i), 10), NextID())
	}
}
