package mock

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

type HttpClient struct {
	mock.Mock
}

func (hc *HttpClient) Do(req *http.Request) (*http.Response, error) {
	args := hc.Called(req)

	if args.Get(0) != nil {
		return args.Get(0).(*http.Response), nil
	}

	return nil, args.Error(1)
}
