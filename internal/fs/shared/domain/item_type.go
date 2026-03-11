package domain

import (
	"encoding/json"
	"fmt"
	"strings"
)

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

// MarshalJSON marshals the ItemType to a JSON string.
func (iT ItemType) MarshalJSON() ([]byte, error) {
	return json.Marshal(iT.String())
}

// UnmarshalJSON unmarshals an ItemType from a JSON string.
func (iT *ItemType) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("ItemType should be a string, got %s", string(b))
	}

	*iT = ParseItemType(s)
	return nil
}
