package http

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

type configurators struct {
	url         string
	fileName    string
	concurrency int
	timeout     time.Duration
}

type Configurator func(d *configurators) error

type Downloader struct {
	Downloader DownloadService
}

func NewHttpService(fileManager FileManager, cfg ...Configurator) (*Downloader, error) {
	var configurators configurators
	for _, c := range cfg {
		err := c(&configurators)
		if err != nil {
			return nil, err
		}
	}

	if configurators.url == "" {
		return nil, fmt.Errorf("url is required")
	}

	if configurators.concurrency == 0 {
		configurators.concurrency = 10
	}

	//if configurators.concurrency > 1 {
	//	supports, err := utils.SupportsResume(configurators.url)
	//	if err != nil {
	//		return nil, fmt.Errorf("failed to check if URL supports resume: %w", err)
	//	}
	//	if supports {
	//		return &Downloader{Downloader: NewResumeDownload(fileManager, configurators.timeout)}, nil
	//	}
	//	fmt.Println("URL does not support resume, falling back to StandardDownload")
	//}

	return &Downloader{Downloader: NewStandardDownload(fileManager, configurators.url, configurators.fileName, configurators.timeout)}, nil
}

// WithURL 配置 URL 並進行驗證
func WithURL(rawURL string) Configurator {
	return func(cfg *configurators) error {
		if rawURL == "" {
			return fmt.Errorf("url cannot be empty")
		}

		// 檢查 URL 是否有效
		parsedURL, err := url.Parse(rawURL)
		if err != nil || !strings.HasPrefix(parsedURL.Scheme, "http") {
			return fmt.Errorf("invalid URL: %s", rawURL)
		}

		cfg.url = rawURL
		return nil
	}
}

func WithFileName(fileName string) Configurator {
	return func(cfg *configurators) error {
		if fileName == "" {
			return fmt.Errorf("fileName cannot be empty")
		}

		cfg.fileName = fileName
		return nil
	}
}

// WithConcurrency 配置併發數量並初始化資源
func WithConcurrency(concurrency int) Configurator {
	return func(cfg *configurators) error {
		if concurrency <= 0 {
			return fmt.Errorf("concurrency must be greater than 0")
		}
		cfg.concurrency = concurrency

		// 在此處可以進行額外的併發設置邏輯，例如資源限制等
		fmt.Printf("Concurrency set to %d\n", concurrency)
		return nil
	}
}

func WithTimeout(timeout time.Duration) Configurator {
	return func(cfg *configurators) error {
		if timeout <= 0 {
			return fmt.Errorf("timeout must be greater than zero")
		}
		cfg.timeout = timeout
		return nil
	}
}
