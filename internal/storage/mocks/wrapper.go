// DONT EDIT: Auto generated

package mock_repository

import (
	"context"
	"homework/internal/model/wrapper"
)

// wrapperStorage ...
type wrapperStorage interface {
	AddWrapper(ctx context.Context, wrapper wrapper.Wrapper, orderId string) error
	Delete(ctx context.Context, orderId string) error
	GetByOrderId(ctx context.Context, orderId string) (wrapper.Wrapper, error)
}
