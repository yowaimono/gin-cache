package cache

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"sort"
	"strings"
)

func generateCacheKey(req *http.Request) string {
	path := req.URL.Path
	query := req.URL.Query()

	var keys []string
	for k := range query {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var params []string
	for _, k := range keys {
		values := query[k]
		for _, v := range values {
			params = append(params, fmt.Sprintf("%s=%s", k, v))
		}
	}

	return fmt.Sprintf("%s?%s", path, strings.Join(params, "&"))
}

func hashString(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}
