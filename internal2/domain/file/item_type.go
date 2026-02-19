package file

import (
	"encoding/json"
	"strings"
)

type ItemType int

const (
	ItemTypeUnknown ItemType = iota - 1
	ItemTypeFile
	ItemTypeFolder
)

func ParseItemType(str string) ItemType {
	ype, ok := map[string]ItemType{
		"file":    ItemTypeFile,
		"folder":  ItemTypeFolder,
		"unknown": ItemTypeUnknown,
	}[strings.ToLower(str)]
	if !ok {
		return ItemTypeUnknown
	}
	return ype
}

func (i ItemType) String() string {
	str, ok := map[ItemType]string{
		ItemTypeFile:    "file",
		ItemTypeFolder:  "folder",
		ItemTypeUnknown: "unknown",
	}[i]
	if !ok {
		return "unknown"
	}
	return str
}

func (i ItemType) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

func (i *ItemType) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	*i = ParseItemType(s)
	return nil
}
