package formatting

import (
	"fmt"
)

const (
	_ = 1 << (10 * iota)
	KiB
	MiB
	GiB
	TiB
	PiB
	EiB
)

// FormatSize returns a human-readable size string.
func FormatSize(bytes int64) string {
	if bytes < KiB {
		return fmt.Sprintf("%d B", bytes)
	}
	if bytes < MiB {
		return fmt.Sprintf("%.2f KiB", float64(bytes)/float64(KiB))
	}
	if bytes < GiB {
		return fmt.Sprintf("%.2f MiB", float64(bytes)/float64(MiB))
	}
	if bytes < TiB {
		return fmt.Sprintf("%.2f GiB", float64(bytes)/float64(GiB))
	}
	return fmt.Sprintf("%.2f TiB", float64(bytes)/float64(TiB))
}
