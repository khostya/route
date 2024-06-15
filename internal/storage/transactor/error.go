package transactor

import (
	"fmt"
)

type TransactionError struct {
	Inner, Rollback error
}

func (t TransactionError) Error() string {
	return fmt.Sprintf("inner=%s rollbacl=%s", t.Inner, t.Rollback)
}
