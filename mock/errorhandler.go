package mock

import (
	"github.com/stretchr/testify/mock"

	"github.com/sredni/messenger"
)

type ErrorHandler struct {
	mock.Mock
}

func (eh *ErrorHandler) HandleError(deliveryError messenger.DeliveryError) {
	eh.Called(deliveryError)
}

