package fs

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type ItemType int

const (
	ItemTypeUnknown ItemType = iota - 1
	ItemTypeFile
	ItemTypeFolder

	itemTypeUnknownString = "unknown"
	itemTypeFileString    = "file"
	itemTypeFolderString  = "folder"
)

var itemTypeToString = map[ItemType]string{
	ItemTypeUnknown: itemTypeUnknownString,
	ItemTypeFile:    itemTypeFileString,
	ItemTypeFolder:  itemTypeFolderString,
}

var stringToItemType = map[string]ItemType{
	itemTypeUnknownString: ItemTypeUnknown,
	itemTypeFileString:    ItemTypeFile,
	itemTypeFolderString:  ItemTypeFolder,
}

func (i ItemType) String() string {
	if s, ok := itemTypeToString[i]; ok {
		return s
	}
	return "unknown"
}

//
// JSON Marshaling
//

func (i ItemType) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

func (i *ItemType) UnmarshalJSON(data []byte) error {
	// Try string first
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		s = strings.ToLower(s)
		if v, ok := stringToItemType[s]; ok {
			*i = v
			return nil
		}
		return fmt.Errorf("invalid ItemType string: %q", s)
	}

	// Try integer fallback
	var n int
	if err := json.Unmarshal(data, &n); err == nil {
		*i = ItemType(n)
		return nil
	}

	return fmt.Errorf("invalid ItemType: %s", string(data))
}

//
// YAML Marshaling (uses string form)
//

func (i ItemType) MarshalYAML() (any, error) {
	return i.String(), nil
}

func (i *ItemType) UnmarshalYAML(unmarshal func(any) error) error {
	// Try string
	var s string
	if err := unmarshal(&s); err == nil {
		s = strings.ToLower(s)
		if v, ok := stringToItemType[s]; ok {
			*i = v
			return nil
		}
		return fmt.Errorf("invalid ItemType string: %q", s)
	}

	// Try integer
	var n int
	if err := unmarshal(&n); err == nil {
		*i = ItemType(n)
		return nil
	}

	return errors.New("invalid ItemType YAML value")
}
