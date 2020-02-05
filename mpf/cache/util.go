/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2019/12/29 0029
 * Time: 11:51
 */
package cache

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

func (util *utilCache) CreateUniqueId() string {
    num := insRedis.conn.Incr(keyUniqueId).Val()
    numStr := strconv.FormatInt(num, 10)

    return time.Now().Format("20060102150405") + numStr[len(numStr)-8:]
}

var (
    keyUniqueId string
    insUtil     *utilCache
)

func init() {
    keyUniqueId = project.RedisPrefix(project.RedisPrefixCommonUniqueId)
    keyUniqueId += "unique"
    insUtil = &utilCache{}

    num := insRedis.conn.Incr(keyUniqueId).Val()
    if num < 100000000 {
        rand.Seed(time.Now().UnixNano())
        randNum := rand.Int63n(50000000) + 100000000
        _, err := insRedis.conn.Set(keyUniqueId, randNum, 0).Result()
        if err != nil {
            panic(mperr.NewCacheRedis(errorcode.CacheRedisOperate, "设置唯一ID自增基数出错", err))
        }
    } else if num > 500000000 {
        reduceNum := num - 100000000 - num%100000000
        insRedis.conn.DecrBy(keyUniqueId, reduceNum)
    }
}

func NewUtilCache() *utilCache {
    return insUtil
}
