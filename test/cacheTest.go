package test

import (
	rootConfig "crypto-bug/config"
	"errors"
	"fmt"
)

type CacheTest struct {
}

func (cacheTest CacheTest) Do() bool {
	err := cacheTest.TestTemporary()
	if err != nil {
		fmt.Println("Temporary cache testing error: " + err.Error())
		return false
	}

	err = cacheTest.TestData()
	if err != nil {
		fmt.Println("Cache testing error: " + err.Error())
		return false
	}

	return true
}

func (cacheTest CacheTest) TestData() error {
	cache := rootConfig.Cache
	cache.Set("first.data.test", 84)
	value, found := cache.Get("first.data.test")
	if !found || value != 84 {
		return errors.New("Not found number in one method")
	}

	cache.Set("second.data.test", 5832)
	err := cacheTest.ExtraGet("second.data.test", 5832)
	if err != nil {
		return err
	}

	return nil
}

func (cacheTest CacheTest) ExtraGet(key string, needed interface{}) error {
	cache := rootConfig.Cache
	value, found := cache.Get(key)
	if !found || value != needed {
		return errors.New("Not found in extra method")
	}
	return nil
}

func (cacheTest CacheTest) TestTemporary() error {
	cache := rootConfig.Cache
	cache.SetTemporary("first.temp.test", 84)
	value, found := cache.Get("first.temp.test")
	if !found || value != 84 {
		return errors.New("Not found number in one method")
	}

	cache.SetTemporary("second.temp.test", 5832)
	err := cacheTest.ExtraGet("second.temp.test", 5832)
	if err != nil {
		return err
	}

	cache.Clear()
	value, found = cache.Get("first.temp.test")
	if found {
		return errors.New("Found value after clearing")
	}

	return nil
}
