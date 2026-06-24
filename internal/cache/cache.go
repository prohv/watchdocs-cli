package cache

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

type Cache struct {
	mu       sync.RWMutex
	filePath string
	data     map[string]string
}

// NewCache returns a new Cache instance.
// If the OS cache directory cannot be located, it returns an error.
func NewCache() (*Cache, error) {
	userCache, err := os.UserCacheDir()
	if err != nil {
		return nil, err
	}

	dir := filepath.Join(userCache, "watchdocs")
	filePath := filepath.Join(dir, "cache.json")

	c := &Cache{
		filePath: filePath,
		data:     make(map[string]string),
	}

	if err := c.load(); err != nil {
		// If the file doesn't exist, start with an empty cache.
		if !os.IsNotExist(err) {
			return nil, err
		}
	}

	return c, nil
}

func (c *Cache) load() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	file, err := os.Open(c.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&c.data); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// Get retrieves a URL from the cache for a given ecosystem and package name.
func (c *Cache) Get(ecosystem, name string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	key := fmt.Sprintf("%s:%s", ecosystem, name)
	url, found := c.data[key]
	return url, found
}

// Set stores a URL in the in-memory cache.
// Call Save() to persist the changes to disk.
func (c *Cache) Set(ecosystem, name, url string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := fmt.Sprintf("%s:%s", ecosystem, name)
	c.data[key] = url
}

// Save persists the in-memory cache to disk atomically.
func (c *Cache) Save() error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	dir := filepath.Dir(c.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Atomic write: write to temp file then rename
	tmpFile, err := os.CreateTemp(dir, "cache-*.json")
	if err != nil {
		return err
	}
	defer func() {
		tmpFile.Close()
		os.Remove(tmpFile.Name()) // clean up if rename didn't happen
	}()

	enc := json.NewEncoder(tmpFile)
	enc.SetIndent("", "  ")
	if err := enc.Encode(c.data); err != nil {
		return err
	}

	if err := tmpFile.Close(); err != nil {
		return err
	}

	return os.Rename(tmpFile.Name(), c.filePath)
}

// Clear clears the cache file from disk and resets the in-memory cache.
func (c *Cache) Clear() error {
	c.mu.Lock()
	c.data = make(map[string]string)
	c.mu.Unlock()

	err := os.Remove(c.filePath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
