package http_handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

type HttpHandlerInterface interface {
	IsHttpStatusOK(resp *http.Response) bool
	SupportRange(resp *http.Response) bool
	DownloadPartial(url string, start int, end int) (*http.Response, error)
	CreateRangeRequest(url string, start, end int) (*http.Response, error)
}
type HTTPHandler struct{}

func NewHTTPClient() *HTTPHandler {
	return &HTTPHandler{}
}

func (c *HTTPHandler) IsHttpStatusOK(resp *http.Response) bool {
	return resp.StatusCode == http.StatusOK
}

func (c *HTTPHandler) SupportRange(resp *http.Response) bool {
	return resp.Header.Get("Accept-Ranges") == "bytes"
}
func (c *HTTPHandler) DownloadPartial(url string, start int, end int) (*http.Response, error) {
	resp, err := c.CreateRangeRequest(url, start, end)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *HTTPHandler) CreateRangeRequest(url string, start, end int) (*http.Response, error) {
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
