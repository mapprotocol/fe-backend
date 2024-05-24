package http

import (
	"fmt"
	"github.com/mapprotocol/ceffu-fe-backend/utils/reqerror"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

const defaultHTTPTimeout = 20 * time.Second

func Request(url, method string, headers http.Header, body io.Reader) ([]byte, error) {
	client := http.Client{
		Timeout: defaultHTTPTimeout,
	}

	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, errors.Wrap(err, "request creation failed")
	}

	if headers != nil {
		request.Header = headers
	}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, fmt.Errorf("response is nil")
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, reqerror.NewExternalRequestError(
			url,
			reqerror.WithCode(strconv.Itoa(resp.StatusCode)),
			reqerror.WithMessage(string(data)),
			reqerror.WithError(errors.New("response status code is not 200")),
		)
	}

	return data, nil
}

func Get(url string, headers http.Header, body io.Reader) ([]byte, error) {
	return Request(url, http.MethodGet, headers, body)
}

func Post(url string, headers http.Header, body io.Reader) ([]byte, error) {
	return Request(url, http.MethodPost, headers, body)
}
