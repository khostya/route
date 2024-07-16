package dto

import (
	"encoding/json"
	"fmt"
	"time"
)

type OnCallMessage struct {
	Args     string
	Method   string
	CalledAt time.Time
}

func (c *OnCallMessage) Marshal() ([]byte, error) {
	return json.Marshal(c)
}

func (c *OnCallMessage) Unmarshal(bytes []byte) error {
	return json.Unmarshal(bytes, c)
}

func (c *OnCallMessage) String() string {
	return fmt.Sprintf("Call(args=%s, method=%s, created_at=%s)", c.Args, c.Method, c.CalledAt)
}
