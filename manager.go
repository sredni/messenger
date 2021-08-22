package messenger

import (
	"context"
	"os"
	"os/signal"
	"sync"

	"github.com/pkg/errors"
)

type Manager interface {
	Run(ctx context.Context, errorHandler ErrorHandler)
	Send(msgs []Message) error
	Sync()
}

type Worker interface {
	Do(ctx context.Context, msg Message) error
}

type Message struct {
	Content []byte
}

type ErrorHandler func(deliveryError DeliveryError)

type DeliveryError struct {
	Err error
	Message
}

func (e DeliveryError) Error() string {
	return e.Err.Error()
}

type manager struct {
	worker      Worker
	concurrency int
	exchange    chan Message
	wg          sync.WaitGroup
}

func NewManager(worker Worker, concurrency int) *manager {
	return &manager{
		worker:      worker,
		concurrency: concurrency,
	}
}

func (m *manager) Run(ctx context.Context, errorHandler ErrorHandler) {
	ctx, cancel := context.WithCancel(ctx)

	m.exchange = make(chan Message, m.concurrency)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	limiter := make(chan int, m.concurrency)
	finish := func() {
		<-limiter
		m.wg.Done()
	}

	go func() {
		for {
			select {
			case <-interrupt:
				cancel()
				m.exchange = nil
			case <-ctx.Done():
				signal.Stop(interrupt)
				m.exchange = nil
				return
			case msg := <-m.exchange:
				limiter <- 0

				go func() {
					err := m.worker.Do(ctx, msg)
					if err != nil && errorHandler != nil {
						errorHandler(DeliveryError{
							Message: msg,
							Err:     err,
						})
					}

					finish()
				}()
			}
		}
	}()
}

func (m *manager) Send(msgs []Message) error {
	m.wg.Add(1)
	if m.exchange == nil {
		m.wg.Done()

		return errors.New("exchange doesn't exists, Run manager first")
	}

	go func() {
		for _, msg := range msgs {
			m.wg.Add(1)
			m.exchange <- msg
		}
		m.wg.Done()
	}()

	return nil
}

func (m *manager) Sync() {
	m.wg.Wait()
}
