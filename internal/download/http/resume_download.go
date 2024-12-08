package http

import (
	"context"
	"download/internal/utils"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

type ResumeDownload struct {
	fileManager FileManager
	conf        Configurators
}

func NewResumeDownload(
	fileManager FileManager,
	conf Configurators,
) *ResumeDownload {
	return &ResumeDownload{
		fileManager: fileManager,
		conf:        conf,
	}
}

func (d *ResumeDownload) Download() error {
	// Use context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), d.conf.timeout)
	defer cancel()

	// get content length
	contentLength, err := d.fetchContentLength(ctx, d.conf.url)
	if err != nil {
		return fmt.Errorf("failed to get content length: %w", err)
	}

	// init global progress bar
	globalProgressBar := utils.NewProgressBar(int64(contentLength))

	// create temp dir for part files
	if err := d.fileManager.CreatDir(d.conf.fileName); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	defer d.cleanupTempFiles(d.conf.fileName)

	if err := d.resumeDownload(ctx, contentLength, globalProgressBar); err != nil {
		return fmt.Errorf("resume download failed: %w", err)
	}

	// merge part files
	if err := d.fileManager.Merge(d.conf.fileName, d.conf.concurrency); err != nil {
		return fmt.Errorf("merge file failed: %w", err)
	}

	return nil
}

func (d *ResumeDownload) fetchContentLength(ctx context.Context, url string) (int, error) {

	req, err := http.NewRequestWithContext(ctx, "HEAD", url, nil)
	if err != nil {
		return 0, fmt.Errorf("建立 HEAD 請求失敗: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("執行 HEAD 請求失敗: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("HTTP 狀態碼錯誤: %s", resp.Status)
	}

	contentLength := resp.ContentLength
	if contentLength <= 0 {
		return 0, errors.New("無效的檔案長度")
	}

	return int(contentLength), nil
}

func (d *ResumeDownload) resumeDownload(ctx context.Context, contentLength int, progressBar *utils.DefaultProgressBar) error {
	var wg sync.WaitGroup
	wg.Add(d.conf.concurrency)
	partSize := contentLength / d.conf.concurrency

	segmentStart := 0
	for i := 0; i < d.conf.concurrency; i++ {
		segmentEnd := segmentStart + partSize
		if i == d.conf.concurrency-1 {
			segmentEnd = contentLength
		}

		go func(index, start, end int) {
			defer wg.Done()
			if err := d.downloadWorker(ctx, index, start, end, progressBar); err != nil {
				log.Printf("segment %d failed: %v", index, err)
			}
		}(i, segmentStart, segmentEnd)

		segmentStart += partSize + 1
	}

	wg.Wait()
	return nil
}

func (d *ResumeDownload) downloadWorker(ctx context.Context, index, start, end int, progressBar *utils.DefaultProgressBar) error {
	req, err := http.NewRequestWithContext(ctx, "GET", d.conf.url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch segment %d: %w", index, err)
	}
	defer resp.Body.Close()

	partFileName := d.fileManager.GetPartFileName(d.conf.fileName, index)
	partFile, err := d.fileManager.CreateDestFile(partFileName)
	if err != nil {
		return fmt.Errorf("failed to create part file: %w", err)
	}
	defer partFile.Close()

	ctxWriter := &utils.ContextWriter{
		Writer:  io.MultiWriter(partFile, progressBar),
		Context: ctx,
	}

	if _, err = io.Copy(ctxWriter, resp.Body); err != nil {
		return fmt.Errorf("failed to write part %d: %w", index, err)
	}

	return nil
}

func (d *ResumeDownload) cleanupTempFiles(baseName string) {
	tempDir := d.fileManager.GetDirName(baseName)
	if err := os.RemoveAll(tempDir); err != nil {
		log.Printf("failed to cleanup temp files: %v", err)
	}
}
