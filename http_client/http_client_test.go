package http_client

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIsHttpStatusOK(t *testing.T) {
	client := NewHTTPClient()
	resp := &http.Response{StatusCode: http.StatusOK}

	if !client.IsHttpStatusOK(resp) {
		t.Errorf("IsHttpStatusOK() = false, want true")
	}
}

func TestSupportRange(t *testing.T) {
	client := NewHTTPClient()
	resp := &http.Response{}
	resp.Header = http.Header{}
	resp.Header.Set("Accept-Ranges", "bytes")

	if !client.SupportRange(resp) {
		t.Errorf("SupportRange() = false, want true")
	}
}

func TestDownloadPartial(t *testing.T) {
	// 模擬一個虛擬的 HTTP 伺服器
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewHTTPClient()

	resp, err := client.DownloadPartial(server.URL, 0, 100)
	if err != nil {
		t.Errorf("DownloadPartial() error = %v, want nil", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("DownloadPartial() StatusCode = %v, want %v", resp.StatusCode, http.StatusOK)
	}

	resp, err = client.DownloadPartial(server.URL, 100, 50)
	if err == nil {
		t.Error("DownloadPartial() expected error, but got nil")
	}
}

func TestMakeRequest(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	client := NewHTTPClient()

	resp, err := client.MakeRequest(server.URL, 0, 100)
	if err != nil {
		t.Errorf("MakeRequest failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	resp, err = client.DownloadPartial(server.URL, 100, 50)
	if err == nil {
		t.Error("DownloadPartial() expected error, but got nil")
	}
}
