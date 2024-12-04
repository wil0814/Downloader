package utils

import (
	"github.com/schollz/progressbar/v3"
)

type ProgressBar struct {
}

func NewProgressBar() *ProgressBar {
	return &ProgressBar{}
}

func (pb *ProgressBar) CreateBar(total int64) *progressbar.ProgressBar {
	bar := progressbar.DefaultBytes(
		total,
		"downloading",
	)

	return bar
}
