package cache_utils

import (
	"fmt"
	"time"

	"github.com/ReneKroon/ttlcache/v2"
)

var (
	Cache           *ttlcache.Cache
	NotFound        = ttlcache.ErrNotFound
	CredentialCache *ttlcache.Cache
)

func Initialize(ttlinSecond int) {
	Cache = ttlcache.NewCache()
	Cache.SetTTL(time.Duration(ttlinSecond) * time.Second)
}

func InitializeCredential(ttlinSecond int) {
	CredentialCache = ttlcache.NewCache()
	Cache.SetTTL(time.Duration(ttlinSecond) * time.Second)
}

func AddOrUpdateCache(key string, value interface{}) {
	fmt.Println("print cach", key, value)
	AddOrUpdateCacheWithTTL(key, value)
}

func AddOrUpdateCacheWithTTL(key string, value interface{}) {
	Cache.Set(key, value)
}

func AddOrUpdateCredentialCache(key string, value interface{}) {
	CredentialCache.Set(key, value)
}

func GetFromCredentialCache(key string) (interface{}, error) {
	return CredentialCache.Get(key)
}
