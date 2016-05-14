package lib

import (
	"github.com/astaxie/beego/cache"
)

var (
	DataCache *_DataCache
)

type _DataCache struct {
	Cache cache.Cache
}

func init() {

	DataCache = &_DataCache{}
	m, err := cache.NewCache("memory", `{"interval":300}`)
	if err != nil {
		ELogger.Error("cache error:",err.Error())
	}

	DataCache.Cache = m
}


