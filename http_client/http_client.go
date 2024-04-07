package http_client

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

type HTTPClient struct{}

func NewHTTPClient() *HTTPClient {
	return &HTTPClient{}
}

func (c *HTTPClient) IsHttpStatusOK(resp *http.Response) bool {
	return resp.StatusCode == http.StatusOK
}

func (c *HTTPClient) SupportRange(resp *http.Response) bool {
	return resp.Header.Get("Accept-Ranges") == "bytes"
}
func (c *HTTPClient) DownloadPartial(url string, start int, end int) (*http.Response, error) {
	resp, err := c.MakeRequest(url, start, end)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *HTTPClient) MakeRequest(url string, start, end int) (*http.Response, error) {
	if start >= end {
		return nil, errors.New("invalid range")
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return resp, nil
}
