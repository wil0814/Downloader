package progress_bar

import (
	"testing"
)

func TestProgressBar_CreateBar(t *testing.T) {
	progress := NewProgressBar()
	total := int64(100)

	bar := progress.CreateBar(total)
	if bar.Current() != 0 {
		t.Errorf("Expected current progress: 0, got: %d", bar.Current())
	}

	progressValue := int64(50)
	bar.SetCurrent(progressValue)

	if bar.Current() != progressValue {
		t.Errorf("Expected current progress: %d, got: %d", progressValue, bar.Current())
	}
}
