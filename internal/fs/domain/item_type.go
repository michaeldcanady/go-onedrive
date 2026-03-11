package domain

import "strings"

type ItemType int

const (
	ItemTypeUnknown ItemType = iota - 1
	ItemTypeFile
	ItemTypeFolder
)

func ParseItemType(str string) ItemType {
	switch strings.ToLower(str) {
	case "unknown":
		return ItemTypeUnknown
	case "file":
		return ItemTypeFile
	case "folder":
		return ItemTypeFolder
	default:
		return ItemTypeUnknown
	}
}

func (iT ItemType) String() string {
	str, ok := map[ItemType]string{
		ItemTypeUnknown: "unknown",
		ItemTypeFile:    "file",
		ItemTypeFolder:  "folder",
	}[iT]

	if !ok {
		return ItemTypeUnknown.String()
	}
	return str
}
