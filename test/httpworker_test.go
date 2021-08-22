package test

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	testingmock "github.com/stretchr/testify/mock"

	"github.com/sredni/messenger"
	"github.com/sredni/messenger/mock"
)

const testUrl = "https://nothing.com/test"

func TestHttpPostWorker(t *testing.T) {
	type testCase struct {
		timeout           time.Duration
		shouldReturnError bool
		setup             func() (*mock.HttpClient, messenger.Message)
	}

	testCases := map[string]testCase{
		"sent": {
			timeout:           1 * time.Second,
			shouldReturnError: false,
			setup: func() (*mock.HttpClient, messenger.Message) {
				msg := messenger.Message{
					Content: []byte("test"),
				}
				hc := &mock.HttpClient{}
				hc.On("Do", testingmock.Anything).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader("")),
				})

				return hc, msg
			},
		},
		"timeout error": {
			timeout:           1 * time.Second,
			shouldReturnError: true,
			setup: func() (*mock.HttpClient, messenger.Message) {
				msg := messenger.Message{
					Content: []byte("test"),
				}
				hc := &mock.HttpClient{}
				hc.
					On("Do", testingmock.Anything).
					Return(nil, errors.New("timeout"))

				return hc, msg
			},
		},
		"wrong response status": {
			timeout:           1 * time.Second,
			shouldReturnError: true,
			setup: func() (*mock.HttpClient, messenger.Message) {
				msg := messenger.Message{
					Content: []byte("test"),
				}
				hc := &mock.HttpClient{}
				hc.On("Do", testingmock.Anything).Return(&http.Response{
					StatusCode: http.StatusNotFound,
					Body:       io.NopCloser(strings.NewReader("")),
				})

				return hc, msg
			},
		},
	}

	ctx := context.Background()

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc := tc
			t.Parallel()

			client, msg := tc.setup()

			worker := messenger.NewHttpPostWorker(client, testUrl, tc.timeout)
			err := worker.Do(ctx, msg)
			if tc.shouldReturnError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			client.AssertExpectations(t)
		})
	}
}
