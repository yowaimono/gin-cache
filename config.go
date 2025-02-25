package cache

import "time"

func SetTTL(ttl time.Duration) {
	defaultTTL = ttl
}
