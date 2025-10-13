package cache

import (
	"container/list"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

// CacheEntry represents a cached item with expiration
type CacheEntry struct {
	Key        string
	Value      []byte
	Expiration time.Time
	Size       int
}

// IsExpired checks if the cache entry has expired
func (e *CacheEntry) IsExpired() bool {
	return time.Now().After(e.Expiration)
}

// Cache provides an in-memory cache with TTL and LRU eviction
type Cache struct {
	mu          sync.RWMutex
	items       map[string]*list.Element
	lruList     *list.List
	ttl         time.Duration
	maxSize     int64
	currentSize int64
	enabled     bool
	hits        int64
	misses      int64
}

// NewCache creates a new cache with specified TTL and max size
func NewCache(ttl time.Duration, maxSizeMB int, enabled bool) *Cache {
	return &Cache{
		items:   make(map[string]*list.Element),
		lruList: list.New(),
		ttl:     ttl,
		maxSize: int64(maxSizeMB) * 1024 * 1024,
		enabled: enabled,
	}
}

// GenerateKey creates a cache key from endpoint and parameters
func GenerateKey(endpoint string, params map[string]string) string {
	paramsStr := endpoint
	for key, value := range params {
		paramsStr += fmt.Sprintf(":%s=%s", key, value)
	}

	hash := sha256.Sum256([]byte(paramsStr))
	return hex.EncodeToString(hash[:])
}

// Get retrieves a value from the cache
func (c *Cache) Get(key string) ([]byte, bool) {
	if !c.enabled {
		return nil, false
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	element, found := c.items[key]
	if !found {
		c.misses++
		return nil, false
	}

	entry := element.Value.(*CacheEntry)

	if entry.IsExpired() {
		c.misses++
		return nil, false
	}

	c.lruList.MoveToFront(element)
	c.hits++

	return entry.Value, true
}

// Set adds or updates a value in the cache
func (c *Cache) Set(key string, value []byte) error {
	if !c.enabled {
		return nil
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	valueSize := len(value)

	if element, found := c.items[key]; found {
		entry := element.Value.(*CacheEntry)
		oldSize := entry.Size
		entry.Value = value
		entry.Expiration = time.Now().Add(c.ttl)
		entry.Size = valueSize

		c.currentSize = c.currentSize - int64(oldSize) + int64(valueSize)
		c.lruList.MoveToFront(element)
		return nil
	}

	entry := &CacheEntry{
		Key:        key,
		Value:      value,
		Expiration: time.Now().Add(c.ttl),
		Size:       valueSize,
	}

	element := c.lruList.PushFront(entry)
	c.items[key] = element
	c.currentSize += int64(valueSize)

	c.evictIfNeeded()

	return nil
}

// evictIfNeeded removes least recently used items if cache is over max size
func (c *Cache) evictIfNeeded() {
	for c.currentSize > c.maxSize && c.lruList.Len() > 0 {
		element := c.lruList.Back()
		if element != nil {
			c.removeElement(element)
		}
	}
}

// removeElement removes an element from the cache
func (c *Cache) removeElement(element *list.Element) {
	c.lruList.Remove(element)
	entry := element.Value.(*CacheEntry)
	delete(c.items, entry.Key)
	c.currentSize -= int64(entry.Size)
}

// Clear removes all items from the cache
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]*list.Element)
	c.lruList.Init()
	c.currentSize = 0
	c.hits = 0
	c.misses = 0
}

// CleanupExpired removes all expired entries from the cache
func (c *Cache) CleanupExpired() int {
	if !c.enabled {
		return 0
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	removed := 0
	var toRemove []*list.Element

	for element := c.lruList.Front(); element != nil; element = element.Next() {
		entry := element.Value.(*CacheEntry)
		if entry.IsExpired() {
			toRemove = append(toRemove, element)
		}
	}

	for _, element := range toRemove {
		c.removeElement(element)
		removed++
	}

	return removed
}

// Stats returns cache statistics
func (c *Cache) Stats() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	hitRate := 0.0
	total := c.hits + c.misses
	if total > 0 {
		hitRate = float64(c.hits) / float64(total) * 100
	}

	return map[string]interface{}{
		"enabled":     c.enabled,
		"items":       len(c.items),
		"size_bytes":  c.currentSize,
		"size_mb":     float64(c.currentSize) / (1024 * 1024),
		"max_size_mb": float64(c.maxSize) / (1024 * 1024),
		"hits":        c.hits,
		"misses":      c.misses,
		"hit_rate":    fmt.Sprintf("%.2f%%", hitRate),
		"ttl_seconds": c.ttl.Seconds(),
	}
}

// StartCleanupRoutine starts a background goroutine to periodically clean up expired entries
func (c *Cache) StartCleanupRoutine(interval time.Duration) {
	if !c.enabled {
		return
	}

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			c.CleanupExpired()
		}
	}()
}
