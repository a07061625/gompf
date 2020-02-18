// Package mpmq consumer
// User: 姜伟
// Time: 2020-02-19 06:41:32
package mpmq

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
)

// IConsumer IConsumer
type IConsumer interface {
    PullData(topic string) (interface{}, error)
    Shutdown() int
}

type consumerMQ struct {
    obj IConsumer
}

func (c *consumerMQ) PullData(topic string) (interface{}, error) {
    return c.obj.PullData(topic)
}

func (c *consumerMQ) Shutdown() int {
    return c.obj.Shutdown()
}

var (
    insConsumer *consumerMQ
)

func init() {
    insConsumer = &consumerMQ{}

    conf := mpf.NewConfig().GetConfig("mpmq")
    mqType := conf.GetString("common." + mpf.EnvProjectKey() + ".mqtype")
    switch mqType {
    case mqTypeRabbit:
        cr := newConsumerRabbit()
        cr.connect(conf)
        insConsumer.obj = cr
    case mqTypeRedis:
        cr := newConsumerRedis()
        cr.RefreshConfig(conf)
    default:
        panic(mperr.NewMQ(errorcode.MQParam, "消息队列类型不支持", nil))
    }
}

// NewConsumer NewConsumer
func NewConsumer() *consumerMQ {
    return insConsumer
}
