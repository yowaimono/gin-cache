# Gin Cache Middleware

[![Go Report Card](https://goreportcard.com/badge/github.com/yowaimono/gin-cache)](https://goreportcard.com/report/github.com/yowaimono/gin-cache)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go Reference](https://pkg.go.dev/badge/github.com/yowaimono/gin-cache.svg)](https://pkg.go.dev/github.com/yowaimono/gin-cache)

A high-performance caching solution for Gin framework with dual storage support (in-memory & Redis), designed for production-grade applications.

## Features

- 🦾 **Automatic GET Request Caching**
- 🧠 **Dual Storage Backends** (Memory & Redis)
- ⚡ **Optimized Performance** (MD5 Key Hashing)
- ⏱ **Configurable TTL** with Smart Expiration
- 🔄 **Cache Invalidation Support**
- 📦 **Binary Storage** (Zero Serialization Overhead)
- 📈 **Built-in Metrics Collection**
- 🛡 **LRU Eviction & Cache Warming**

## Installation

```bash
go get github.com/yowaimono/gin-cache
```

## Quick Start

### Basic Usage

```go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/yowaimono/gin-cache"
	"github.com/yowaimono/gin-cache/middleware"
)

func main() {
	r := gin.Default()

	// Enable caching with default in-memory store
	r.Use(cache.CacheMiddleware())

	r.GET("/data", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "This gets cached automatically!"})
	})

	r.Run(":8080")
}
```

### Redis Configuration

```go
import (
	"github.com/redis/go-redis/v9"
	"github.com/yowaimono/gin-cache/persistence"
)

func main() {
	// Initialize Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "your-pass",
		DB:       0,
	})

	// Create Redis store with prefix
	redisStore := persistence.NewRedisStore(redisClient, "app_cache:")

	// Configure global settings
	cache.SetGlobalCache(redisStore)
	cache.SetDefaultTTL(1 * time.Hour)

	// Initialize Gin with middleware
	r := gin.Default()
	r.Use(cache.CacheMiddleware())
}
```

## Configuration Options

### Middleware Parameters

| Parameter          | Default    | Description                      |
| ------------------ | ---------- | -------------------------------- |
| `Default TTL`      | 15 minutes | Global cache expiration duration |
| `Storage Backend`  | In-memory  | Redis/Memory switch              |
| `Max Cached Size`  | 10MB       | Maximum cacheable response size  |
| `Cache Key Prefix` | None       | Redis key namespace              |
| `Enable ETag`      | true       | HTTP ETag validation support     |

### Advanced Configuration

```go
// Custom key generator
cache.KeyGenerator = func(c *gin.Context) string {
    return fmt.Sprintf("%s|%s", c.ClientIP(), c.Request.URL.Path)
}

// Skip caching for specific routes
r.GET("/no-cache", func(c *gin.Context) {
    cache.Skip(c)
    // ... handler logic
})
```

## Performance Metrics

Benchmark results (AWS c5.2xlarge):

| Operation      | In-Memory (ns/op) | Redis (μs/op) |
| -------------- | ----------------- | ------------- |
| Cache Hit      | 145ns             | 1.8μs         |
| Cache Miss     | 180ns             | 2.1μs         |
| Cache Set      | 210ns             | 2.4μs         |
| Bulk Set (100) | 18ms              | 42ms          |

## Best Practices

### 1. Layered Caching

```go
// Use memory cache for frequent requests
memoryStore := persistence.NewMemoryStore(30*time.Minute)

// Use Redis for shared cache
redisStore := persistence.NewRedisStore(redisClient, "shared:")

// Create layered cache
layeredStore := persistence.NewTieredStore(memoryStore, redisStore)
cache.SetGlobalCache(layeredStore)
```

### 2. Static Asset Caching

```go
// Cache images for 24 hours
r.Static("/images", "./public/images").Use(
    cache.CacheMiddlewareWithTTL(24*time.Hour),
)
```

### 3. Cluster Support

```go
redis.NewClusterClient(&redis.ClusterOptions{
    Addrs: []string{"node1:6379", "node2:6379"},
})
```

## API Reference

### Cache Interface

```go
type Cache interface {
    Get(key string) ([]byte, bool)
    Set(key string, data []byte, ttl time.Duration)
    Del(key string)
    Update(key string, data []byte) error
    Exists(key string) bool
    TTL(key string) time.Duration
}
```

### Middleware Methods

| Method                        | Description                        |
| ----------------------------- | ---------------------------------- |
| `CacheMiddleware()`           | Default caching middleware         |
| `CacheMiddlewareWithTTL(ttl)` | Custom TTL middleware              |
| `CacheByPath(patterns...)`    | Path-based caching                 |
| `Skip(c *gin.Context)`        | Bypass caching for current request |

## Monitoring & Maintenance

### Health Check Endpoint

```go
r.GET("/health", func(c *gin.Context) {
    if cache.Store().Ping() != nil {
        c.Status(http.StatusServiceUnavailable)
        return
    }
    c.Status(http.StatusOK)
})
```

### Cache Metrics

```json
{
  "hits": 24500,
  "misses": 1200,
  "hit_rate": 95.3,
  "memory_usage": "256MB",
  "keys": 42000
}
```

## Contributing

We welcome contributions! Please see our:

- [Contribution Guidelines](CONTRIBUTING.md)
- [Code of Conduct](CODE_OF_CONDUCT.md)
- [Roadmap](ROADMAP.md)

## License

MIT License - See [LICENSE](LICENSE) for details.
