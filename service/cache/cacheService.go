package cache

type Cache struct {
	Temp map[string]interface{}
	Data map[string]interface{}
}

func (cache *Cache) SetTemporary(key string, value interface{}) {
	cache.Temp[key] = value
}

func (cache *Cache) Set(key string, value interface{}) {
	cache.Data[key] = value
}

func (cache *Cache) Get(key string) (interface{}, bool) {
	item, found := cache.Temp[key]
	if !found {
		item, found = cache.Data[key]
		if !found {
			return nil, false
		}
	}

	return item, true
}

func (cache Cache) Clear() {
	for key, _ := range cache.Temp {
		delete(cache.Temp, key)
	}
}
