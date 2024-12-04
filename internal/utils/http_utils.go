package utils

import (
	"fmt"
	"net/http"
)

// SupportsResume 檢查給定 URL 是否支持斷點續傳
func SupportsResume(url string) (bool, error) {
	// 發送 HTTP HEAD 請求
	resp, err := http.Head(url)
	if err != nil {
		return false, fmt.Errorf("failed to check URL: %w", err)
	}
	defer resp.Body.Close()

	// 檢查 "Accept-Ranges" 是否包含 "bytes"
	if resp.Header.Get("Accept-Ranges") == "bytes" {
		return true, nil
	}
	return false, nil
}
