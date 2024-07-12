//go:build integration

package kafka

import (
	"homework/internal/dto"
	"time"
)

func NewOnCallMessage() dto.CallMessage {
	return dto.CallMessage{
		CalledAt: time.Now(),
		Args:     []string{"--user=1", "--id=1"},
		Method:   "call",
	}
}
