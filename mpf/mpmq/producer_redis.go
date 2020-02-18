// Package mpmq producer_redis
// User: 姜伟
// Time: 2020-02-19 06:42:18
package mpmq

import (
    "github.com/a07061625/gompf/mpf/mpcache"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
)

type producerRedis struct {
    topicPrefix string
}

func (p *producerRedis) SendTopicData(topic string, data []string) int {
    num := 0
    pushData := make([]interface{}, 0)
    for _, v := range data {
        pushData = append(pushData, v)
    }
    redisKey := p.topicPrefix + topic
    addNum := mpcache.NewRedis().GetConn().RPush(redisKey, pushData...).Val()
    num += int(addNum)

    return num
}

func newProducerRedis() *producerRedis {
    p := &producerRedis{""}
    p.topicPrefix = project.RedisPrefix(project.RedisPrefixMQRedis)
    return p
}
