package cache

import (
	"fmt"
	"sync"
	"time"
)

// Cache interface defines the methods for the cache
type Cache interface {
	Get(key string) ([]byte, bool)
	Set(key string, data []byte, ttl time.Duration)
	Del(key string)
	Update(key string, data []byte) error
}

type MemoryCache struct {
	data sync.Map
	ttl  sync.Map
}

func (m *MemoryCache) Get(key string) ([]byte, bool) {
	value, ok := m.data.Load(key)
	if !ok {
		return nil, false
	}
	if expire, ok := m.ttl.Load(key); ok {
		if time.Now().After(expire.(time.Time)) {
			m.Del(key)
			return nil, false
		}
	}
	return value.([]byte), true
}

func (m *MemoryCache) Set(key string, data []byte, ttl time.Duration) {
	m.data.Store(key, data)
	if ttl > 0 {
		m.ttl.Store(key, time.Now().Add(ttl))
	}
}

func (m *MemoryCache) Del(key string) {
	m.data.Delete(key)
	m.ttl.Delete(key)
}

func (m *MemoryCache) Update(key string, data []byte) error {
	if _, ok := m.data.Load(key); !ok {
		return fmt.Errorf("key not exists")
	}
	m.data.Store(key, data)
	return nil
}
