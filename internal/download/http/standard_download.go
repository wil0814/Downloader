package http

import (
	"context"
	"download/internal/utils"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type StandardDownload struct {
	fileManager FileManager
	url         string
	fileName    string
	timeout     time.Duration
}

func NewStandardDownload(
	fileManager FileManager,
	url string,
	fileName string,
	timeout time.Duration,
) *StandardDownload {
	return &StandardDownload{
		fileManager: fileManager,
		url:         url,
		fileName:    fileName,
		timeout:     timeout,
	}
}

func (d *StandardDownload) Download() error {
	// Use context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", d.url, nil)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform HTTP request: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			err = fmt.Errorf("failed to close response body: %w", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected HTTP status: %s", resp.Status)
	}

	return d.writeResponseToFile(ctx, resp.Body, d.fileName, resp.ContentLength)
}

func (d *StandardDownload) writeResponseToFile(ctx context.Context, respBody io.ReadCloser, filename string, contentLength int64) error {
	localFile, err := d.fileManager.CreateDestFile(filename)
	if err != nil {
		return err
	}
	defer func() {
		localFile.Close()
	}()

	// 初始化進度條
	progressBar := utils.NewProgressBar(contentLength)
	ctxWriter := &utils.ContextWriter{
		Writer:  io.MultiWriter(localFile, progressBar),
		Context: ctx,
	}

	if _, err := io.Copy(ctxWriter, respBody); err != nil {
		_ = os.Remove(filename)
		return fmt.Errorf("failed to write response to file: %w", err)
	}

	return nil
}
