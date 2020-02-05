/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/5 0005
 * Time: 15:11
 */
package mpmq

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/cache"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/spf13/viper"
)

type consumerRedis struct {
    pullNum     int64
    topicPrefix string
}

func (c *consumerRedis) RefreshConfig(conf *viper.Viper) {
    pullNum := conf.GetInt("consumer." + mpf.EnvProjectKey() + ".redis.pullnum")
    if pullNum <= 0 {
        panic(mperr.NewMQRedis(errorcode.MQRedisParam, "最大拉取数量不合法", nil))
    }
    c.pullNum = int64(pullNum)
}

func (c *consumerRedis) PullData(topic string) (interface{}, error) {
    redisKey := c.topicPrefix + topic
    dataList := cache.NewRedis().GetConn().LRange(redisKey, 0, c.pullNum).Val()
    dataNum := len(dataList)
    if dataNum > 0 {
        cache.NewRedis().GetConn().LTrim(redisKey, int64(dataNum), -1)
    }

    return dataList, nil
}

func (c *consumerRedis) Shutdown() int {
    return 0
}

func newConsumerRedis() *consumerRedis {
    c := &consumerRedis{0, ""}
    c.topicPrefix = project.RedisPrefix(project.RedisPrefixMQRedis)
    return c
}
