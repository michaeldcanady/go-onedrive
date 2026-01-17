package disk

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/michaeldcanady/go-onedrive/internal/cachev2/abstractions"
)

var _ abstractions.Cache[any, any] = (*Cache[any, any])(nil)

type Cache[K comparable, V any] struct {
	path            string
	keySerializer   abstractions.SerializerDeserializer[K]
	valueSerializer abstractions.SerializerDeserializer[V]

	mu    sync.RWMutex
	index map[string]int64 // serialized key → offset
}

// New creates a new disk cache and loads the index.
func New[K comparable, V any](
	path string,
	ks abstractions.SerializerDeserializer[K],
	vs abstractions.SerializerDeserializer[V],
) (*Cache[K, V], error) {

	c := &Cache[K, V]{
		path:            path,
		keySerializer:   ks,
		valueSerializer: vs,
		index:           make(map[string]int64),
	}

	if err := c.loadIndex(); err != nil {
		return nil, err
	}

	return c, nil
}

// loadIndex scans the file and builds the in-memory index.
func (c *Cache[K, V]) loadIndex() error {
	f, err := os.OpenFile(c.path, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	offset := int64(0)
	for {
		var keyLen uint32
		if err := binary.Read(f, binary.LittleEndian, &keyLen); err != nil {
			if err == io.EOF { // is expected
				break
			}
			return err
		}

		keyBytes := make([]byte, keyLen)
		if _, err := f.Read(keyBytes); err != nil {
			return err
		}

		var valLen uint32
		if err := binary.Read(f, binary.LittleEndian, &valLen); err != nil {
			return err
		}

		if _, err := f.Seek(int64(valLen), io.SeekCurrent); err != nil {
			return err
		}

		// Index the key
		c.index[string(keyBytes)] = offset

		// Move offset
		offset, _ = f.Seek(0, io.SeekCurrent)
	}

	return nil
}

// Clear removes all entries from the cache.
//
// This operation deletes the underlying file and resets the in‑memory index.
// After Clear returns, the cache behaves as if newly created.
func (c *Cache[K, V]) Clear(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := os.Remove(c.path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	c.index = make(map[string]int64)
	return nil
}

// GetEntry retrieves the most recent value associated with the given key.
//
// The key is serialized and looked up in the in‑memory index. If found,
// the method seeks to the corresponding file offset, reads the record,
// deserializes the value, and returns it.
//
// If the key does not exist, an error is returned.
func (c *Cache[K, V]) GetEntry(ctx context.Context, key K) (*abstractions.Entry[K, V], error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	serializedKey, err := c.keySerializer.Serialize(key)
	if err != nil {
		return nil, err
	}

	offset, ok := c.index[string(serializedKey)]
	if !ok {
		return nil, fmt.Errorf("key not found")
	}

	f, err := os.Open(c.path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if _, err := f.Seek(offset, io.SeekStart); err != nil {
		return nil, err
	}

	// Skip key
	var keyLen uint32
	binary.Read(f, binary.LittleEndian, &keyLen)
	f.Seek(int64(keyLen), io.SeekCurrent)

	// Read value
	var valLen uint32
	binary.Read(f, binary.LittleEndian, &valLen)

	valBytes := make([]byte, valLen)
	if _, err := f.Read(valBytes); err != nil {
		return nil, err
	}

	value, err := c.valueSerializer.Deserialize(valBytes)
	if err != nil {
		return nil, err
	}

	return abstractions.NewEntry(key, value), nil
}

// NewEntry creates a new cache entry for the given key with a zero‑value value.
//
// This method does not write anything to disk. It simply constructs a
// CacheEntry that the caller may later pass to SetEntry.
// NewEntry is useful when callers want to populate a value incrementally
// before committing it to the cache.
func (c *Cache[K, V]) NewEntry(_ context.Context, key K) (*abstractions.Entry[K, V], error) {
	var zero V
	return abstractions.NewEntry(key, zero), nil
}

// Remove deletes the given key from the cache.
//
// The key is removed from the in‑memory index, and the underlying file is
// compacted to remove stale entries. Compaction rewrites the file to contain
// only live entries and updates all index offsets accordingly.
func (c *Cache[K, V]) Remove(key K) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	serializedKey, err := c.keySerializer.Serialize(key)
	if err != nil {
		return err
	}

	delete(c.index, string(serializedKey))

	return c.rewriteFile()
}

// SetEntry writes or overwrites a cache entry.
//
// The key and value are serialized and appended to the end of the file.
// The in‑memory index is updated to point to the new record. Older records
// for the same key remain in the file until compaction occurs.
func (c *Cache[K, V]) SetEntry(ctx context.Context, entry *abstractions.Entry[K, V]) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	serializedKey, err := c.keySerializer.Serialize(entry.GetKey())
	if err != nil {
		return err
	}

	serializedValue, err := c.valueSerializer.Serialize(entry.GetValue())
	if err != nil {
		return err
	}

	f, err := os.OpenFile(c.path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	offset, _ := f.Seek(0, io.SeekEnd)

	// Write key
	if err := binary.Write(f, binary.LittleEndian, uint32(len(serializedKey))); err != nil {
		return err
	}
	if _, err := f.Write(serializedKey); err != nil {
		return err
	}

	// Write value
	if err := binary.Write(f, binary.LittleEndian, uint32(len(serializedValue))); err != nil {
		return err
	}

	if _, err := f.Write(serializedValue); err != nil {
		return err
	}

	// Update index
	c.index[string(serializedKey)] = offset

	return nil
}

// rewriteFile compacts the file by rewriting only live entries.
func (c *Cache[K, V]) rewriteFile() error {
	tmp := c.path + ".tmp"

	f, err := os.OpenFile(tmp, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	newIndex := make(map[string]int64)

	for keyStr := range c.index {
		offset, _ := f.Seek(0, io.SeekEnd)

		// Deserialize key
		key, err := c.keySerializer.Deserialize([]byte(keyStr))
		if err != nil {
			return err
		}

		// Get value
		entry, err := c.GetEntry(context.Background(), key)
		if err != nil {
			continue
		}

		keyBytes := []byte(keyStr)
		valBytes, _ := c.valueSerializer.Serialize(entry.GetValue())

		binary.Write(f, binary.LittleEndian, uint32(len(keyBytes)))
		f.Write(keyBytes)

		binary.Write(f, binary.LittleEndian, uint32(len(valBytes)))
		f.Write(valBytes)

		newIndex[keyStr] = offset
	}

	// Replace old file
	if err := os.Rename(tmp, c.path); err != nil {
		return err
	}

	c.index = newIndex
	return nil
}

func (c *Cache[K, V]) KeySerializer() abstractions.Serializer[K] {
	return c.keySerializer
}
