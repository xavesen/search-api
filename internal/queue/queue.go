package queue

import "context"

type Queue interface {
	WriteMessage(ctx context.Context, message []byte) error
}