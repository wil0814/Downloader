package utils

import (
	"github.com/schollz/progressbar/v3"
)

type DefaultProgressBar struct {
	bar *progressbar.ProgressBar
}

func NewProgressBar(total int64) *DefaultProgressBar {
	return &DefaultProgressBar{
		bar: progressbar.DefaultBytes(total, "downloading"),
	}
}

func (pb *DefaultProgressBar) Write(p []byte) (n int, err error) {
	n = len(p)
	err = pb.bar.Add(n)
	return n, err
}
