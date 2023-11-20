package cache

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var GoCache *cache.Cache
func InitGoCaches() {
	// 初始化 BigCache
	GoCache = cache.New(5*time.Minute, 10*time.Minute)
}