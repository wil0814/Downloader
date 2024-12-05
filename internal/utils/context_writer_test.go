package utils_test

import (
	"bytes"
	"context"
	"download/internal/utils"
	"errors"
	"testing"
)

func TestContextWriter_Write(t *testing.T) {
	tests := []struct {
		name           string
		contextFunc    func() context.Context // func to generate test context
		input          []byte                 // data to be written
		expectedError  error                  // expected error
		expectedBytes  int                    // expected number of bytes written
		expectedOutput string                 // expected writer output
	}{
		{
			name: "Context not canceled",
			contextFunc: func() context.Context {
				return context.Background()
			},
			input:          []byte("hello world"),
			expectedError:  nil,
			expectedBytes:  11,
			expectedOutput: "hello world",
		},
		{
			name: "Context canceled",
			contextFunc: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			input:          []byte("hello world"),
			expectedError:  context.Canceled,
			expectedBytes:  0,
			expectedOutput: "",
		},
	}

	// for each test case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// init test
			buffer := &bytes.Buffer{}
			ctx := tt.contextFunc()
			writer := &utils.ContextWriter{
				Writer:  buffer,
				Context: ctx,
			}

			// run test
			n, err := writer.Write(tt.input)

			// validate test result
			if !errors.Is(err, tt.expectedError) {
				t.Errorf("Expected error %v, got %v", tt.expectedError, err)
			}

			// validate number of bytes written
			if n != tt.expectedBytes {
				t.Errorf("Expected written bytes %d, got %d", tt.expectedBytes, n)
			}

			// validate writer output
			if buffer.String() != tt.expectedOutput {
				t.Errorf("Expected output %q, got %q", tt.expectedOutput, buffer.String())
			}
		})
	}
}
