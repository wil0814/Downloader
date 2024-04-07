package progress_bar

import (
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
	"time"
)

type ProgressBar struct {
	p *mpb.Progress
}

func NewProgressBar() *ProgressBar {
	p := mpb.New(
		mpb.WithWidth(60),
		mpb.WithRefreshRate(time.Millisecond),
	)
	return &ProgressBar{p: p}
}

func (pb *ProgressBar) CreateBar(total int64) *mpb.Bar {
	//bar := pb.p.New(total,
	//	mpb.BarStyle().Rbound("|"),
	//	mpb.PrependDecorators(
	//		decor.Counters(decor.SizeB1024(0), "% .2f / % .2f"),
	//	),
	//	mpb.AppendDecorators(
	//		decor.EwmaETA(decor.ET_STYLE_GO, 30),
	//		decor.Name(" ] "),
	//		decor.EwmaSpeed(decor.SizeB1024(0), "% .2f", 60),
	//	),
	//)
	bar := pb.p.AddBar(total,
		mpb.PrependDecorators(decor.Counters(decor.SizeB1024(0), "% .1f / % .1f")),
		mpb.AppendDecorators(decor.Percentage()),
	)
	return bar
}
