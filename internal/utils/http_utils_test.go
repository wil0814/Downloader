package utils_test

import (
	"download/internal/utils"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestSupportsResume tests the SupportsResume function.
func TestSupportsResume(t *testing.T) {
	tests := []struct {
		name           string
		setupServer    func() *httptest.Server
		expectedResult bool
		expectError    bool
	}{
		{
			name: "Supports resume with Accept-Ranges header",
			setupServer: func() *httptest.Server {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.Method == http.MethodHead {
						w.Header().Set("Accept-Ranges", "bytes")
						w.WriteHeader(http.StatusOK)
					} else {
						w.WriteHeader(http.StatusMethodNotAllowed)
					}
				}))
				return server
			},
			expectedResult: true,
			expectError:    false,
		},
		{
			name: "Does not support resume without Accept-Ranges header",
			setupServer: func() *httptest.Server {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.Method == http.MethodHead {
						w.WriteHeader(http.StatusOK)
					} else {
						w.WriteHeader(http.StatusMethodNotAllowed)
					}
				}))
				return server
			},
			expectedResult: false,
			expectError:    false,
		},
		{
			name: "HTTP request fails",
			setupServer: func() *httptest.Server {
				return nil // No server setup to simulate network error
			},
			expectedResult: false,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var server *httptest.Server
			if tt.setupServer != nil {
				server = tt.setupServer()
				if server != nil {
					defer server.Close()
				}
			}

			url := ""
			if server != nil {
				url = server.URL
			}

			result, err := utils.SupportsResume(url)

			// Check if an error was expected
			if (err != nil) != tt.expectError {
				t.Fatalf("expected error: %v, got: %v", tt.expectError, err)
			}

			// Check the result
			if result != tt.expectedResult {
				t.Errorf("expected result: %v, got: %v", tt.expectedResult, result)
			}
		})
	}
}
