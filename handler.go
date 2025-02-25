package cache

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	defaultCache     Cache
	defaultTTL       = 15 * time.Second
	defaultCacheOnce sync.Once
)

func SetGlobalCache(c Cache) {
	defaultCache = c
}

func initDefaultCache() {
	defaultCacheOnce.Do(func() {
		if defaultCache == nil {
			defaultCache = &MemoryCache{
				data: sync.Map{},
				ttl:  sync.Map{},
				
			}
		}
	})
}

func CacheMiddleware() gin.HandlerFunc {
	// initialize default cache if not set
	// just in case if default cache is not set by user
	initDefaultCache()
	return func(c *gin.Context) {
		if c.Request.Method != http.MethodGet {
			c.Next()
			return
		}

		key := generateCacheKey(c.Request)
		hashedKey := hashString(key)

		if data, ok := defaultCache.Get(hashedKey); ok {
			c.Data(http.StatusOK, "application/json", data)
			c.Abort()
			return
		}

		writer := newResponseCaptureWriter(c.Writer)
		c.Writer = writer
		c.Next()

		if c.Writer.Status() < 200 || c.Writer.Status() >= 300 {
			return
		}

		defaultCache.Set(hashedKey, writer.body.Bytes(), defaultTTL)
	}
}
