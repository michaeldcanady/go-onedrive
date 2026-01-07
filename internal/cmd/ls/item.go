package ls

import (
	"time"

	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

// Item represents a pseudo linux filesystem object.
type Item struct {
	Name         string    `json:"name" yaml:"name"`
	IsFolder     bool      `json:"is_folder" yaml:"is_folder"`
	Size         int64     `json:"size" yaml:"size"`
	ModifiedTime time.Time `json:"modified_time" yaml:"modified_time"`
}

// toItem converts models.DriveItemable to ls commands Item type.
func toItem(it models.DriveItemable) Item {
	return Item{
		Name:         safeName(it),
		IsFolder:     it.GetFolder() != nil,
		Size:         getSize(it),
		ModifiedTime: getModifiedTime(it),
	}
}
