package mock

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/sredni/messenger"
)

type Worker struct {
	mock.Mock
}

func (w *Worker) Do(_ context.Context, msg messenger.Message) error {
	args := w.Called(msg)

	return args.Error(0)
}

