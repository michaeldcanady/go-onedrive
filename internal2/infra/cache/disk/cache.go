package disk

import (
	"context"
	"encoding/binary"
	"errors"
	"io"
	"os"
	"sync"

	"github.com/gofrs/flock"
	domaincache "github.com/michaeldcanady/go-onedrive/internal2/domain/cache"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/cache/bolt"
)

var _ domaincache.Cache[any] = (*Cache[string, any])(nil)

type Cache[K comparable, V any] struct {
	path            string
	lock            *flock.Flock
	keySerializer   domaincache.SerializerDeserializer[K]
	valueSerializer domaincache.SerializerDeserializer[V]

	mu    sync.RWMutex
	index map[string]int64 // serialized key → offset
}

// New creates a new disk cache and loads the index.
func New[K comparable, V any](
	path string,
	ks domaincache.SerializerDeserializer[K],
	vs domaincache.SerializerDeserializer[V],
) (*Cache[K, V], error) {

	c := &Cache[K, V]{
		path:            path,
		lock:            flock.New(path + ".lock"),
		keySerializer:   ks,
		valueSerializer: vs,
		index:           make(map[string]int64),
	}

	if err := c.loadIndex(); err != nil {
		return nil, err
	}

	return c, nil
}

// Get implements [domaincache.Cache].
func (c *Cache[K, V]) Get(ctx context.Context, key string) (V, error) {
	var zero V
	// Shared cross-process lock
	if err := c.lock.RLock(); err != nil {
		return zero, err
	}
	defer c.lock.Unlock()

	c.mu.RLock()
	defer c.mu.RUnlock()

	offset, ok := c.index[key]
	if !ok {
		return zero, bolt.ErrKeyNotFound
	}

	k, err := c.keySerializer.Deserialize([]byte(key))
	if err != nil {
		return zero, err
	}

	entry, err := c.readEntryAtOffset(offset, k)
	if err != nil {
		return zero, err
	}
	return entry.GetValue(), nil
}

// Set implements [domaincache.Cache].
func (c *Cache[K, V]) Set(ctx context.Context, key string, value V) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Exclusive cross-process lock
	if err := c.lock.Lock(); err != nil {
		return err
	}
	defer c.lock.Unlock()

	k, err := c.keySerializer.Deserialize([]byte(key))
	if err != nil {
		return err
	}

	return c.setEntry(ctx, domaincache.NewEntry(k, value))
}

// Delete implements [domaincache.Cache].
func (c *Cache[K, V]) Delete(ctx context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.index, key)

	return c.rewriteFile()
}

// List implements [domaincache.Cache].
func (c *Cache[K, V]) List(ctx context.Context, callback func(key string, value V) error) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for kStr, offset := range c.index {
		k, err := c.keySerializer.Deserialize([]byte(kStr))
		if err != nil {
			return err
		}
		entry, err := c.readEntryAtOffset(offset, k)
		if err != nil {
			return err
		}
		if err := callback(kStr, entry.GetValue()); err != nil {
			return err
		}
	}
	return nil
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
func (c *Cache[K, V]) Clear(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.lock.Lock(); err != nil {
		return err
	}
	defer c.lock.Unlock()

	if err := os.Remove(c.path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	c.index = make(map[string]int64)
	return nil
}

// setEntry writes or overwrites a cache entry.
func (c *Cache[K, V]) setEntry(ctx context.Context, entry *domaincache.Entry[K, V]) error {
	serializedKey, err := c.keySerializer.Serialize(entry.GetKey())
	if err != nil {
		return err
	}

	serializedValue, err := c.valueSerializer.Serialize(entry.GetValue())
	if err != nil {
		return err
	}

	// Append new record
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

	// Compact file
	return c.rewriteFile()
}

func (c *Cache[K, V]) readEntryAtOffset(offset int64, key K) (*domaincache.Entry[K, V], error) {
	f, err := os.Open(c.path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if _, err := f.Seek(offset, io.SeekStart); err != nil {
		return nil, err
	}

	// Read key length
	var keyLen uint32
	if err := binary.Read(f, binary.LittleEndian, &keyLen); err != nil {
		return nil, err
	}

	// Skip key bytes
	if _, err := f.Seek(int64(keyLen), io.SeekCurrent); err != nil {
		return nil, err
	}

	// Read value length
	var valLen uint32
	if err := binary.Read(f, binary.LittleEndian, &valLen); err != nil {
		return nil, err
	}

	valBytes := make([]byte, valLen)
	if _, err := f.Read(valBytes); err != nil {
		return nil, err
	}

	value, err := c.valueSerializer.Deserialize(valBytes)
	if err != nil {
		return nil, err
	}

	return domaincache.NewEntry(key, value), nil
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

	for keyStr, oldOffset := range c.index {
		offset, _ := f.Seek(0, io.SeekEnd)

		// Deserialize key
		key, err := c.keySerializer.Deserialize([]byte(keyStr))
		if err != nil {
			return err
		}

		// Read value directly from disk
		entry, err := c.readEntryAtOffset(oldOffset, key)
		if err != nil {
			return err
		}

		keyBytes := []byte(keyStr)
		valBytes, err := c.valueSerializer.Serialize(entry.GetValue())
		if err != nil {
			return err
		}

		// Write key
		if err := binary.Write(f, binary.LittleEndian, uint32(len(keyBytes))); err != nil {
			return err
		}
		if _, err := f.Write(keyBytes); err != nil {
			return err
		}

		// Write value
		if err := binary.Write(f, binary.LittleEndian, uint32(len(valBytes))); err != nil {
			return err
		}
		if _, err := f.Write(valBytes); err != nil {
			return err
		}

		newIndex[keyStr] = offset
	}

	// Atomic replace
	if err := os.Rename(tmp, c.path); err != nil {
		return err
	}

	c.index = newIndex
	return nil
}

func (c *Cache[K, V]) KeySerializer() domaincache.Serializer[K] {
	return c.keySerializer
}
