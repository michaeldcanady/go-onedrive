package fs

import "time"

type Item struct {
	ID       string
	Name     string
	Path     string
	Type     ItemType
	Size     int64
	Modified time.Time
}
