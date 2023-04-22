package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		c := NewCache(3)
		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)
		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)
		wasInCache = c.Set("ccc", 300)
		require.False(t, wasInCache)

		c.Clear()

		val, ok := c.Get("aaa")
		require.Nil(t, val)
		require.False(t, ok)

		val, ok = c.Get("bbb")
		require.Nil(t, val)
		require.False(t, ok)

		val, ok = c.Get("ccc")
		require.Nil(t, val)
		require.False(t, ok)
	})
}

func TestCacheMultithreading(t *testing.T) {
	t.Skip() // Remove me if task with asterisk completed.

	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}

func TestNewCache(t *testing.T) {
	if NewCache(5) == nil {
		t.Error("cache should not be nil")
	}
}

func TestSet(t *testing.T) {
	cache := NewCache(2)

	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key3", "value3")

	val, ok := cache.Get("key1")

	require.False(t, ok)
	require.Nil(t, val)
}

func TestSetExistingKey(t *testing.T) {
	cache := NewCache(2)

	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key1", "value3")

	val, _ := cache.Get("key1")

	require.Equal(t, "value3", val)
}

func TestGetNonExistingKey(t *testing.T) {
	cache := NewCache(2)

	val, _ := cache.Get("key1")

	require.Equal(t, val, nil)
}

func TestClear(t *testing.T) {
	cache := NewCache(2)

	cache.Set("key1", "value1")
	cache.Set("key2", "value2")

	cache.Clear()

	val, _ := cache.Get("key1")
	val1, _ := cache.Get("key2")

	require.Nil(t, val)
	require.Nil(t, val1)
}
