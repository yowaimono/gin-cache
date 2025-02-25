package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/yowaimono/cache"
)

// Use default cache middleware (no extra settings required)
func runDefaultCacheExample() {
	r := gin.Default()

	// Use default cache middleware
	r.GET("/default", cache.CacheMiddleware(), func(c *gin.Context) {
		// 模拟数据库查询
		time.Sleep(500 * time.Millisecond)
		c.JSON(http.StatusOK, gin.H{
			"message":   "default cache example",
			"timestamp": time.Now().Unix(),
		})
	})

	fmt.Println("Running default cache example on :8080")
	r.Run(":8080")
}

// Use Redis cache middleware (requires Redis server running on localhost:6379)
func runRedisCacheExample() {
	// Init Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis地址
		Password: "",               // 密码
		DB:       0,                // 数据库
	})

	// Test Redis connection
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
	}

	// Set global cache store to Redis cache store
	redisStore := cache.NewRedisStore(redisClient, "myapp")
	cache.SetGlobalCache(redisStore)

	r := gin.Default()
	r.GET("/redis", cache.CacheMiddleware(), func(c *gin.Context) {
		// Database query simulation
		time.Sleep(500 * time.Millisecond)
		c.JSON(http.StatusOK, gin.H{
			"message":   "redis cache example",
			"timestamp": time.Now().Unix(),
		})
	})

	fmt.Println("Running Redis cache example on :8081")
	r.Run(":8081")
}

// Use custom cache middleware (requires custom cache implementation)
type CustomCache struct {
	data map[string][]byte
	mu   sync.RWMutex
}

func (c *CustomCache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, ok := c.data[key]
	return val, ok
}

func (c *CustomCache) Set(key string, data []byte, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = data

	// Simple implementation of TTL expiration
	go func() {
		time.Sleep(ttl)
		c.Del(key)
	}()
}

func (c *CustomCache) Del(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
}

func (c *CustomCache) Update(key string, data []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.data[key]; !ok {
		return fmt.Errorf("key not exists")
	}
	c.data[key] = data
	return nil
}

func runCustomCacheExample() {
	// Set global cache store to custom cache store
	customCache := &CustomCache{
		data: make(map[string][]byte),
	}
	cache.SetGlobalCache(customCache)

	r := gin.Default()
	r.GET("/custom", cache.CacheMiddleware(), func(c *gin.Context) {
		// Database query simulation
		time.Sleep(500 * time.Millisecond)
		c.JSON(http.StatusOK, gin.H{
			"message":   "custom cache example",
			"timestamp": time.Now().Unix(),
		})
	})

	fmt.Println("Running custom cache example on :8082")
	r.Run(":8082")
}

func main() {

	// runDefaultCacheExample()
	// runRedisCacheExample()
	runCustomCacheExample()
}
