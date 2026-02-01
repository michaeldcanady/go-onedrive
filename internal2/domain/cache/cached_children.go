package cache

import (
	"fmt"

	"github.com/microsoft/kiota-abstractions-go/serialization"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

var _ serialization.Parsable = (*CachedChildren)(nil)

type CachedChildren struct {
	ETag  string
	Items models.DriveItemCollectionResponseable
}

func CreateCachedChildrenFromDiscriminatorValue(_ serialization.ParseNode) (serialization.Parsable, error) {
	return &CachedChildren{}, nil
}

// GetFieldDeserializers implements [serialization.Parsable].
func (c *CachedChildren) GetFieldDeserializers() map[string]func(serialization.ParseNode) error {
	return map[string]func(serialization.ParseNode) error{
		"etag": func(pn serialization.ParseNode) error {
			etag, err := pn.GetStringValue()
			if err != nil {
				return err
			}

			if etag == nil {
				c.ETag = ""
				return nil
			}
			c.ETag = *etag
			return nil
		},
		"items": func(pn serialization.ParseNode) error {
			obj, err := pn.GetObjectValue(models.CreateDriveItemCollectionResponseFromDiscriminatorValue)
			if err != nil {
				return err
			}

			typedObj, ok := obj.(models.DriveItemCollectionResponseable)
			if !ok {
				return fmt.Errorf("obj is not %T", typedObj)
			}

			c.Items = typedObj
			return nil
		},
	}
}

// Serialize implements [serialization.Parsable].
func (c *CachedChildren) Serialize(writer serialization.SerializationWriter) error {
	if err := writer.WriteStringValue("etag", &c.ETag); err != nil {
		return err
	}
	if err := writer.WriteObjectValue("items", c.Items); err != nil {
		return err
	}
	return nil
}
