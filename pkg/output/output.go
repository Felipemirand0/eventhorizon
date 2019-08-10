package output

import (
	"context"
)

type Output interface {
	Send(context.Context, interface{}) error
	Close() error
}
