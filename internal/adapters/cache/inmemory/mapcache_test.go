package inmemory

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMapCache_SetAndGet(t *testing.T) {
	cache := NewCache[string, int](5*time.Minute, 10*time.Minute)

	// Test Set and Get
	cache.Set("key1", 42)
	value, exists := cache.Get("key1")

	assert.True(t, exists)
	assert.Equal(t, 42, value)
}

func TestMapCache_GetNonExistent(t *testing.T) {
	cache := NewCache[string, int](5*time.Minute, 10*time.Minute)

	// Test Get for non-existent key
	value, exists := cache.Get("nonexistent")

	assert.False(t, exists)
	assert.Equal(t, 0, value) // zero value for int
}

func TestMapCache_Delete(t *testing.T) {
	cache := NewCache[string, int](5*time.Minute, 10*time.Minute)

	// Set a value
	cache.Set("key1", 42)
	value, exists := cache.Get("key1")
	assert.True(t, exists)
	assert.Equal(t, 42, value)

	// Delete the value
	cache.Delete("key1")
	value, exists = cache.Get("key1")
	assert.False(t, exists)
	assert.Equal(t, 0, value)
}

func TestMapCache_GetAll(t *testing.T) {
	cache := NewCache[string, int](5*time.Minute, 10*time.Minute)

	// Test GetAll with empty cache
	values, exists := cache.GetAll()
	assert.False(t, exists)
	assert.Nil(t, values)

	// Add some values
	cache.Set("key1", 1)
	cache.Set("key2", 2)
	cache.Set("key3", 3)

	// Test GetAll with values
	values, exists = cache.GetAll()
	assert.True(t, exists)
	assert.Len(t, values, 3)

	// Check that all values are present (order may vary)
	valueMap := make(map[int]bool)
	for _, v := range values {
		valueMap[v] = true
	}
	assert.True(t, valueMap[1])
	assert.True(t, valueMap[2])
	assert.True(t, valueMap[3])
}

func TestMapCache_ConcurrentAccess(t *testing.T) {
	cache := NewCache[string, int](5*time.Minute, 10*time.Minute)

	// Test concurrent access
	done := make(chan bool, 10)

	// Start multiple goroutines
	for i := 0; i < 10; i++ {
		go func(i int) {
			key := "key" + string(rune(i))
			cache.Set(key, i)
			value, exists := cache.Get(key)
			assert.True(t, exists)
			assert.Equal(t, i, value)
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify all values are present
	values, exists := cache.GetAll()
	assert.True(t, exists)
	assert.Len(t, values, 10)
}

func TestMapCache_Overwrite(t *testing.T) {
	cache := NewCache[string, int](5*time.Minute, 10*time.Minute)

	// Set initial value
	cache.Set("key1", 42)
	value, exists := cache.Get("key1")
	assert.True(t, exists)
	assert.Equal(t, 42, value)

	// Overwrite with new value
	cache.Set("key1", 100)
	value, exists = cache.Get("key1")
	assert.True(t, exists)
	assert.Equal(t, 100, value)
}

func TestMapCache_DifferentTypes(t *testing.T) {
	// Test with string values
	stringCache := NewCache[string, string](5*time.Minute, 10*time.Minute)
	stringCache.Set("key1", "value1")
	value, exists := stringCache.Get("key1")
	assert.True(t, exists)
	assert.Equal(t, "value1", value)

	// Test with struct values
	type TestStruct struct {
		ID   int
		Name string
	}

	structCache := NewCache[string, TestStruct](5*time.Minute, 10*time.Minute)
	testStruct := TestStruct{ID: 1, Name: "test"}
	structCache.Set("key1", testStruct)
	valueStruct, exists := structCache.Get("key1")
	assert.True(t, exists)
	assert.Equal(t, testStruct, valueStruct)
}

func TestMapCache_Expiration(t *testing.T) {
	// Create a cache with a very short TTL
	cache := NewCache[string, int](100*time.Millisecond, 200*time.Millisecond)

	// Set a value
	cache.Set("key1", 42)
	value, exists := cache.Get("key1")
	assert.True(t, exists)
	assert.Equal(t, 42, value)

	// Wait for the item to expire
	time.Sleep(150 * time.Millisecond)

	// Try to get the expired item
	value, exists = cache.Get("key1")
	assert.False(t, exists)
	assert.Equal(t, 0, value)

	// Set another value and wait for cleanup
	cache.Set("key2", 100)
	time.Sleep(250 * time.Millisecond)

	// The cleanup goroutine should have removed the expired item
	// We can't directly check the internal map, but we can check that Get doesn't find it
	value, exists = cache.Get("key2")
	assert.False(t, exists)
	assert.Equal(t, 0, value)
}
