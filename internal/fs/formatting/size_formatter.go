package formatting

import (
	"fmt"
)

const (
	_ = 1 << (10 * iota)
	// KiB represents a kibibyte (1024 bytes).
	KiB
	// MiB represents a mebibyte (1024 KiB).
	MiB
	// GiB represents a gibibyte (1024 MiB).
	GiB
	// TiB represents a tebibyte (1024 GiB).
	TiB
	// PiB represents a pebibyte (1024 TiB).
	PiB
	// EiB represents an exbibyte (1024 PiB).
	EiB
)

// FormatSize converts a byte count into a human-readable string using binary suffixes (KiB, MiB, etc.).
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
