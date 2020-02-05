/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/5 0005
 * Time: 15:12
 */
package mpmq

import (
    "github.com/a07061625/gompf/mpf/cache"
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
    addNum := cache.NewRedis().GetConn().RPush(redisKey, pushData...).Val()
    num += int(addNum)

    return num
}

func newProducerRedis() *producerRedis {
    p := &producerRedis{""}
    p.topicPrefix = project.RedisPrefix(project.RedisPrefixMQRedis)
    return p
}
