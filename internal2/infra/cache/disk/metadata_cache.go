package disk

import (
	"context"
	"encoding/binary"
	"io"
	"os"
	"sync"

	"github.com/michaeldcanady/go-onedrive/internal2/infra/cache/abstractions"
)

// MetadataCache wraps any Cache[K, V] and adds typed metadata support.
// M is the metadata type, fully controlled by the application.
//
// This type does NOT modify the underlying cache’s behavior.
// It simply adds a parallel metadata storage layer.
type MetadataCache[K comparable, V any, M any] struct {
	// Embedded base cache
	abstractions.Cache[K, V]

	// Mutex protecting metadata operations
	cacheMx sync.RWMutex

	// Path to the metadata file
	metadataPath string

	// Index mapping serialized keys → file offsets for metadata
	metadataIndex map[string]int64

	// Serializer for metadata type M
	metadataSerializer abstractions.SerializerDeserializer[M]
}

func NewMetadataCache[K comparable, V, M any](cache abstractions.Cache[K, V], metadataPath string, metadataSerializer abstractions.SerializerDeserializer[M]) *MetadataCache[K, V, M] {
	return &MetadataCache[K, V, M]{
		Cache:              cache,
		metadataPath:       metadataPath,
		metadataIndex:      map[string]int64{},
		metadataSerializer: metadataSerializer,
	}
}

func (mc *MetadataCache[K, V, M]) GetMetadata(ctx context.Context, key K) (M, error) {
	var zero M

	// Serialize key
	serializedKey, err := mc.Cache.KeySerializer().Serialize(key)
	if err != nil {
		return zero, err
	}
	keyStr := string(serializedKey)

	mc.cacheMx.RLock()
	offset, ok := mc.metadataIndex[keyStr]
	mc.cacheMx.RUnlock()

	if !ok {
		// No metadata stored yet
		return zero, nil
	}

	// Open metadata file
	f, err := os.Open(mc.metadataPath)
	if err != nil {
		return zero, err
	}
	defer f.Close()

	// Seek to metadata offset
	if _, err := f.Seek(offset, io.SeekStart); err != nil {
		return zero, err
	}

	// Read length prefix
	var length uint32
	if err := binary.Read(f, binary.LittleEndian, &length); err != nil {
		return zero, err
	}

	// Read metadata bytes
	buf := make([]byte, length)
	if _, err := io.ReadFull(f, buf); err != nil {
		return zero, err
	}

	// Deserialize metadata
	md, err := mc.metadataSerializer.Deserialize(buf)
	if err != nil {
		return zero, err
	}

	return md, nil
}

func (mc *MetadataCache[K, V, M]) SetMetadata(ctx context.Context, key K, md M) error {
	// Serialize key
	serializedKey, err := mc.Cache.KeySerializer().Serialize(key)
	if err != nil {
		return err
	}
	keyStr := string(serializedKey)

	// Serialize metadata
	data, err := mc.metadataSerializer.Serialize(md)
	if err != nil {
		return err
	}

	mc.cacheMx.Lock()
	defer mc.cacheMx.Unlock()

	// Open metadata file for append
	f, err := os.OpenFile(mc.metadataPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	// Determine offset
	offset, err := f.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}

	// Write length prefix
	length := uint32(len(data))
	if err := binary.Write(f, binary.LittleEndian, length); err != nil {
		return err
	}

	// Write metadata bytes
	if _, err := f.Write(data); err != nil {
		return err
	}

	// Update index
	mc.metadataIndex[keyStr] = offset

	return nil
}
