package lib

import (
	"github.com/astaxie/beego/cache"
	"strconv"
	"time"
)

var (
	DataCache *_DataCache
)

type _DataCache struct {
	Cache cache.Cache
}

type AdCountIndex struct {
	Uuid    string
	Ad_id   int
	Ad_type int
}

func init() {

	DataCache = &_DataCache{}
	m, err := cache.NewCache("memory", `{"interval":300}`)
	if err != nil {
		ELogger.Error("cache error:",err.Error())
	}

	DataCache.Cache = m
}

func GetAdCount(adInfo AdCountIndex, keyType string) int {
	key := ""
	ret := 0
	if keyType == "ad" {
		key = generateCacheKey(adInfo, "ad")
	} else {
		key = generateCacheKey(adInfo, "user")
	}
	max_show_num := DataCache.Cache.Get(key)
	if max_show_num == nil {
		DataCache.Cache.Put(key, 0, AppConfig.Cachetime)
		return ret
	}

	ret = max_show_num.(int)

	return ret
}

func incrementAdCountByOne(adInfo AdCountIndex) {
	key := generateCacheKey(adInfo, "ad")
	err := DataCache.Cache.Incr(key)
	if err != nil {
		DataCache.Cache.Put(key, 0, AppConfig.Cachetime)
	}
}

func IncrementUserAdCountByOne(adInfo AdCountIndex) {
	key := generateCacheKey(adInfo, "user")
	err := DataCache.Cache.Incr(key)
	if err != nil {
		DataCache.Cache.Put(key, 0, AppConfig.Cachetime)
	} else {
		incrementAdCountByOne(adInfo)
	}
}

func generateCacheKey(adInfo AdCountIndex, keyType string) string {
	ret := ""
	if keyType == "user" {
		ret = adInfo.Uuid + "-" + strconv.Itoa(adInfo.Ad_id) + "-" + strconv.Itoa(adInfo.Ad_type) + "-" + time.Now().Format("2006-01-02")
	} else {
		ret = strconv.Itoa(adInfo.Ad_id) + "-" + strconv.Itoa(adInfo.Ad_type) + "-" + time.Now().Format("2006-01-02")
	}

	return ret
}
