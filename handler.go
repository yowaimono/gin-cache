package cache

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/singleflight"
)

var (
	defaultCache     Cache
	defaultTTL       = 15 * time.Second
	defaultCacheOnce sync.Once

	flightGroup singleflight.Group // singleflight Group 实例
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

		val, err, _ := flightGroup.Do(hashedKey, func() (interface{}, error) {
			writer := newResponseCaptureWriter(c.Writer)
			c.Writer = writer

			c.Next()

			if c.Writer.Status() < 200 || c.Writer.Status() >= 300 {
				return nil, fmt.Errorf("status code error: %d", c.Writer.Status()) // 返回错误，singleflight 不会缓存错误结果
			}

			cacheData := writer.body.Bytes()
			defaultCache.Set(hashedKey, cacheData, defaultTTL) // 缓存数据
			return cacheData, nil
		})

		if err != nil {
			// singleflight 执行失败，可能是下游 Handler 错误，直接返回
			return
		}

		data := val.([]byte)
		c.Data(http.StatusOK, "application/json", data)
		c.Abort()
	}
}
