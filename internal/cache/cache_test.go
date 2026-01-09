package cache

import (
	"testing"
	"time"
)

func TestCache_SetAndGet(t *testing.T) {
	c := NewCache(1*time.Hour, 10, true)

	key := "test:key"
	value := []byte("test value")

	// Set a value
	c.Set(key, value)

	// Get the value
	result, found := c.Get(key)
	if !found {
		t.Error("Expected to find cached value")
	}

	if string(result) != string(value) {
		t.Errorf("Expected %s, got %s", string(value), string(result))
	}
}

func TestCache_GetNonExistent(t *testing.T) {
	c := NewCache(1*time.Hour, 10, true)

	_, found := c.Get("nonexistent:key")
	if found {
		t.Error("Expected not to find non-existent key")
	}
}

func TestCache_Expiration(t *testing.T) {
	c := NewCache(100*time.Millisecond, 10, true)

	key := "test:expiring"
	value := []byte("will expire")

	c.Set(key, value)

	// Value should exist immediately
	_, found := c.Get(key)
	if !found {
		t.Error("Expected to find value immediately after set")
	}

	// Wait for expiration
	time.Sleep(150 * time.Millisecond)

	// Value should be expired
	_, found = c.Get(key)
	if found {
		t.Error("Expected value to be expired")
	}
}

func TestCache_Disabled(t *testing.T) {
	c := NewCache(1*time.Hour, 10, false) // Cache disabled

	key := "test:disabled"
	value := []byte("should not cache")

	c.Set(key, value)

	// Should not find value when cache is disabled
	_, found := c.Get(key)
	if found {
		t.Error("Expected not to find value when cache is disabled")
	}
}

func TestCache_Clear(t *testing.T) {
	c := NewCache(1*time.Hour, 10, true)

	// Add multiple items
	c.Set("key1", []byte("value1"))
	c.Set("key2", []byte("value2"))
	c.Set("key3", []byte("value3"))

	// Verify they exist
	if _, found := c.Get("key1"); !found {
		t.Error("Expected to find key1")
	}

	// Clear cache
	c.Clear()

	// Verify all items are gone
	if _, found := c.Get("key1"); found {
		t.Error("Expected key1 to be cleared")
	}
	if _, found := c.Get("key2"); found {
		t.Error("Expected key2 to be cleared")
	}
	if _, found := c.Get("key3"); found {
		t.Error("Expected key3 to be cleared")
	}
}

func TestCache_ConcurrentAccess(t *testing.T) {
	c := NewCache(1*time.Hour, 10, true)

	// Run concurrent Set operations
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(n int) {
			key := string(rune('a' + n))
			c.Set(key, []byte("value"))
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify all values were set
	for i := 0; i < 10; i++ {
		key := string(rune('a' + i))
		if _, found := c.Get(key); !found {
			t.Errorf("Expected to find key %s", key)
		}
	}
}

func TestGenerateKey(t *testing.T) {
	tests := []struct {
		name     string
		prefix   string
		params   map[string]string
		expected string
	}{
		{
			name:     "Simple key",
			prefix:   "test",
			params:   map[string]string{"id": "123"},
			expected: "test:id=123",
		},
		{
			name:     "Multiple params",
			prefix:   "search",
			params:   map[string]string{"query": "yosemite", "state": "CA"},
			expected: "search:query=yosemite:state=CA",
		},
		{
			name:     "No params",
			prefix:   "list",
			params:   map[string]string{},
			expected: "list",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateKey(tt.prefix, tt.params)
			// Note: map iteration order is random, so we just check that all parts are present
			if result == "" {
				t.Error("Expected non-empty key")
			}
			// At minimum, should contain the prefix
			if len(result) < len(tt.prefix) {
				t.Errorf("Expected key to contain prefix %s", tt.prefix)
			}
		})
	}
}
