// Package mpcache util
// User: 姜伟
// Time: 2020-02-19 06:17:54
package mpcache

import (
    "math/rand"
    "strconv"
    "time"

    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
)

type utilCache struct {
}

// CreateUniqueID 生成唯一ID
func (util *utilCache) CreateUniqueID() string {
    num := insRedis.conn.Incr(keyUniqueID).Val()
    numStr := strconv.FormatInt(num, 10)

    return time.Now().Format("20060102150405") + numStr[len(numStr)-8:]
}

var (
    keyUniqueID string
    insUtil     *utilCache
)

func init() {
    keyUniqueID = project.RedisPrefix(project.RedisPrefixCommonUniqueID)
    keyUniqueID += "unique"
    insUtil = &utilCache{}

    num := insRedis.conn.Incr(keyUniqueID).Val()
    if num < 100000000 {
        rand.Seed(time.Now().UnixNano())
        randNum := rand.Int63n(50000000) + 100000000
        _, err := insRedis.conn.Set(keyUniqueID, randNum, 0).Result()
        if err != nil {
            panic(mperr.NewCacheRedis(errorcode.CacheRedisOperate, "设置唯一ID自增基数出错", err))
        }
    } else if num > 500000000 {
        reduceNum := num - 100000000 - num%100000000
        insRedis.conn.DecrBy(keyUniqueID, reduceNum)
    }
}

// NewUtilCache NewUtilCache
func NewUtilCache() *utilCache {
    return insUtil
}
