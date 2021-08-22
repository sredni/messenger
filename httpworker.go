package messenger

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type httpPostWorker struct {
	client  HttpClient
	url     string
	timeout time.Duration
}

func NewHttpPostWorker(client HttpClient, url string, timeout time.Duration) *httpPostWorker {
	return &httpPostWorker{
		client:  client,
		url:     url,
		timeout: timeout,
	}
}

func (hpw *httpPostWorker) Do(ctx context.Context, msg Message) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, hpw.timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(
		timeoutCtx,
		http.MethodPost,
		hpw.url,
		bytes.NewBuffer(msg.Content),
	)
	if err != nil {
		return errors.Wrap(err, "unable to create request")
	}

	resp, err := hpw.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "unable to do request")
	}
	defer resp.Body.Close()
	if resp.StatusCode >= http.StatusMultipleChoices {
		return errors.New(fmt.Sprintf("Invalid response status code: %d", resp.StatusCode))
	}

	return nil
}
