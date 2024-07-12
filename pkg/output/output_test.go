package output

import (
	"github.com/stretchr/testify/require"
	"sort"
	"testing"
)

func TestController_Subscribe_Push(t *testing.T) {
	t.Parallel()

	controller := NewController[int]()

	res := controller.Subscribe()

	controller.push(1)
	controller.push(2)

	expected := []int{<-res, <-res}
	sort.Ints(expected)
	require.Equal(t, []int{1, 2}, expected)

	controller.Close()
	_, ok := <-res
	require.False(t, ok)
}

func TestController_Subscribe_Add(t *testing.T) {
	t.Parallel()

	controller := NewController[int]()

	res := controller.Subscribe()

	chan1 := make(chan int, 1)
	chan1 <- 1
	close(chan1)

	chan2 := make(chan int, 1)
	chan2 <- 2
	close(chan2)

	controller.add(chan1)
	controller.add(chan2)

	expected := []int{<-res, <-res}
	sort.Ints(expected)

	require.Equal(t, []int{1, 2}, expected)

	controller.Close()
	_, ok := <-res
	require.False(t, ok)
}
