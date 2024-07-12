package cli

import (
	"context"
	"testing"
)

func TestCli_RunExit(t *testing.T) {
	t.Parallel()

	mocks := newMocks(t)
	cli := NewCLI(Deps{Service: mocks.mockOrderService})
	ctx := context.Background()

	cli.Run(ctx, []string{exit})

	<-cli.Exit()
}
