package file

type ItemType int

const (
	ItemTypeUnknown ItemType = iota - 1
	ItemTypeFile
	ItemTypeFolder
)
