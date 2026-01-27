package driveservice

import (
	"fmt"

	"github.com/microsoft/kiota-abstractions-go/serialization"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

var _ serialization.Parsable = (*CachedItem)(nil)

type CachedItem struct {
	ETag string
	Item models.DriveItemable
}

func CreateCachedItemFromDiscriminatorValue(_ serialization.ParseNode) (serialization.Parsable, error) {
	return &CachedItem{}, nil
}

// GetFieldDeserializers implements [serialization.Parsable].
func (c *CachedItem) GetFieldDeserializers() map[string]func(serialization.ParseNode) error {
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
		"item": func(pn serialization.ParseNode) error {
			obj, err := pn.GetObjectValue(models.CreateDriveFromDiscriminatorValue)
			if err != nil {
				return err
			}

			typedObj, ok := obj.(models.DriveItemable)
			if !ok {
				return fmt.Errorf("obj is not %T", typedObj)
			}

			c.Item = typedObj
			return nil
		},
	}
}

// Serialize implements [serialization.Parsable].
func (c *CachedItem) Serialize(writer serialization.SerializationWriter) error {
	if err := writer.WriteStringValue("etag", &c.ETag); err != nil {
		return err
	}
	if err := writer.WriteObjectValue("item", c.Item); err != nil {
		return err
	}
	return nil
}
