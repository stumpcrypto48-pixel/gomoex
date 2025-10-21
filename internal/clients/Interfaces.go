package client

import (
	"context"
)

type GetInterface interface {
	GetData(context.Context, string) ([]byte, error)
}
