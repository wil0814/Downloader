package progress_bar

import (
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
	"time"
)

type ProgressBar struct {
	Progress *mpb.Progress
}

func NewProgressBar() *ProgressBar {
	p := mpb.New(
		mpb.WithWidth(60),
		mpb.WithRefreshRate(time.Millisecond),
	)
	return &ProgressBar{Progress: p}
}

func (pb *ProgressBar) CreateBar() *mpb.Bar {
	bar := pb.Progress.AddBar(0,
		mpb.PrependDecorators(decor.Counters(decor.SizeB1024(0), "% .1f / % .1f")),
		mpb.AppendDecorators(decor.Percentage()),
	)
	return bar
}
