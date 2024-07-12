package dto

import (
	"encoding/json"
	"fmt"
	"time"
)

type CallMessage struct {
	Args     []string
	Method   string
	CalledAt time.Time
}

func (c *CallMessage) Marshal() ([]byte, error) {
	return json.Marshal(c)
}

func (c *CallMessage) Unmarshal(bytes []byte) error {
	return json.Unmarshal(bytes, c)
}

func (c *CallMessage) String() string {
	return fmt.Sprintf("Call(args=%s, method=%s, created_at=%s)", c.Args, c.Method, c.CalledAt)
}
