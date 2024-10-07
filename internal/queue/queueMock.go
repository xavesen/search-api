package queue

import "context"

type QueueMock struct {
	Error		error
}

func (qm *QueueMock) WriteMessage(ctx context.Context, message []byte) error {
	return qm.Error
}