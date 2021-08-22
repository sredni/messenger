package test

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/sredni/messenger"
	"github.com/sredni/messenger/mock"
)

func TestManager(t *testing.T) {
	type testCase struct {
		concurrency int
		setup       func() (*mock.Worker, []messenger.Message, *mock.ErrorHandler)
	}

	testCases := map[string]testCase{
		"one message success": {
			concurrency: 1,
			setup: func() (*mock.Worker, []messenger.Message, *mock.ErrorHandler) {
				msg1 := messenger.Message{
					Content: []byte("test1"),
				}
				eh := &mock.ErrorHandler{}
				w := &mock.Worker{}
				w.On("Do", msg1).Return(nil).Once()

				return w, []messenger.Message{msg1}, eh
			},
		},
		"one message error": {
			concurrency: 1,
			setup: func() (*mock.Worker, []messenger.Message, *mock.ErrorHandler) {
				msg1 := messenger.Message{
					Content: []byte("test1"),
				}
				de1 := messenger.DeliveryError{
					Message: msg1,
					Err:     errors.New("Unable to deliver message1"),
				}
				eh := &mock.ErrorHandler{}
				w := &mock.Worker{}
				w.On("Do", msg1).Return(de1.Err).Once()
				eh.On("HandleError", de1).Once()

				return w, []messenger.Message{msg1}, eh
			},
		},
		"multiple messages no error": {
			concurrency: 3,
			setup: func() (*mock.Worker, []messenger.Message, *mock.ErrorHandler) {
				msg1 := messenger.Message{
					Content: []byte("test1"),
				}
				msg2 := messenger.Message{
					Content: []byte("test2"),
				}
				msg3 := messenger.Message{
					Content: []byte("test3"),
				}
				eh := &mock.ErrorHandler{}
				w := &mock.Worker{}
				w.On("Do", msg1).Return(nil).Once()
				w.On("Do", msg2).Return(nil).Once()
				w.On("Do", msg3).Return(nil).Once()

				return w, []messenger.Message{msg1, msg2, msg3}, eh
			},
		},
		"multiple messages some errors": {
			concurrency: 3,
			setup: func() (*mock.Worker, []messenger.Message, *mock.ErrorHandler) {
				msg1 := messenger.Message{
					Content: []byte("test1"),
				}
				msg2 := messenger.Message{
					Content: []byte("test2"),
				}
				de2 := messenger.DeliveryError{
					Message: msg2,
					Err:     errors.New("Unable to deliver message2"),
				}
				msg3 := messenger.Message{
					Content: []byte("test3"),
				}
				eh := &mock.ErrorHandler{}
				w := &mock.Worker{}
				w.On("Do", msg1).Return(nil).Once()
				w.On("Do", msg2).Return(de2.Err).Once()
				eh.On("HandleError", de2).Once()
				w.On("Do", msg3).Return(nil).Once()

				return w, []messenger.Message{msg1, msg2, msg3}, eh
			},
		},
		"multiple messages all errors": {
			concurrency: 3,
			setup: func() (*mock.Worker, []messenger.Message, *mock.ErrorHandler) {
				msg1 := messenger.Message{
					Content: []byte("test1"),
				}
				de1 := messenger.DeliveryError{
					Message: msg1,
					Err:     errors.New("Unable to deliver message1"),
				}
				msg2 := messenger.Message{
					Content: []byte("test2"),
				}
				de2 := messenger.DeliveryError{
					Message: msg2,
					Err:     errors.New("Unable to deliver message2"),
				}
				msg3 := messenger.Message{
					Content: []byte("test3"),
				}
				de3 := messenger.DeliveryError{
					Message: msg3,
					Err:     errors.New("Unable to deliver message3"),
				}
				eh := &mock.ErrorHandler{}
				w := &mock.Worker{}
				w.On("Do", msg1).Return(de1.Err).Once()
				eh.On("HandleError", de1).Once()
				w.On("Do", msg2).Return(de2.Err).Once()
				eh.On("HandleError", de2).Once()
				w.On("Do", msg3).Return(de3.Err).Once()
				eh.On("HandleError", de3).Once()

				return w, []messenger.Message{msg1, msg2, msg3}, eh
			},
		},
	}
	ctx := context.Background()

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc := tc
			t.Parallel()

			worker, messages, errorHandler := tc.setup()
			mngr := messenger.NewManager(worker, tc.concurrency)

			mngr.Run(ctx, errorHandler.HandleError)
			err := mngr.Send(messages)
			assert.NoError(t, err)
			mngr.Sync()

			worker.AssertExpectations(t)
			errorHandler.AssertExpectations(t)
		})
	}
}

func TestManager_Run(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	worker := &mock.Worker{}
	msg := messenger.Message{
		Content: []byte("test1"),
	}
	mngr := messenger.NewManager(worker, 1)

	err := mngr.Send([]messenger.Message{msg})
	assert.Error(t, err)
	mngr.Run(ctx, nil)

	worker.On("Do", msg).Return(nil).Once()

	err = mngr.Send([]messenger.Message{
		{
			Content: []byte("test1"),
		},
	})
	assert.NoError(t, err)
	mngr.Sync()
	worker.AssertExpectations(t)
}
