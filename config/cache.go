package config

import "crypto-bug/service/cache"

var Cache *cache.Cache

func cacheSetup() {
	data := make(map[string]interface{})
	Cache = &cache.Cache{
		Temp: data,
		Data: data,
	}
}
