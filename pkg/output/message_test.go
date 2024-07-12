package output

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBuildMessageChan_Send(t *testing.T) {
	t.Parallel()

	ch := make(chan int, 10)
	ch <- 1
	close(ch)

	out := BuildMessageChan[int]("1", ch)
	require.Equal(t, Message[int]{"1", 1}, <-out)
}

func TestFilterMessageChan_Send(t *testing.T) {
	t.Parallel()

	m1 := Message[int]{"1", 1}
	m2 := Message[int]{"2", 2}
	ch := make(chan Message[int], 10)

	ch <- m2
	ch <- m1
	close(ch)

	out := FilterMessageChan("1", ch)

	m := <-out
	require.Equal(t, m1, m)
}
