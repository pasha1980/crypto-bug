package cacheTest2

import (
	rootConfig "crypto-bug/config"
	"errors"
)

func CacheTestInDifferentPackage(key string, needed interface{}) error {
	value, found := rootConfig.Cache.Get(key)
	if !found || value != needed {
		return errors.New("Not found in extra package")
	}
	return nil
}

func (c CacheTestStruct) CacheTestInDifferentPackageAndStruct(key string, needed interface{}) error {
	value, found := rootConfig.Cache.Get(key)
	if !found || value != needed {
		return errors.New("Not found in extra package and struct")
	}
	return nil
}

type CacheTestStruct struct {
}
