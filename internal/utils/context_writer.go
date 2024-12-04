package utils

import (
	"context"
	"io"
)

type ContextWriter struct {
	Writer  io.Writer
	Context context.Context
}

func (c *ContextWriter) Write(p []byte) (int, error) {
	select {
	case <-c.Context.Done():
		return 0, context.Canceled
	default:
		return c.Writer.Write(p)
	}
}
