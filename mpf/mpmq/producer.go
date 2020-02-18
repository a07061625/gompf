// Package mpmq producer
// User: 姜伟
// Time: 2020-02-19 06:42:40
package mpmq

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
)

// IProducer IProducer
type IProducer interface {
    SendTopicData(topic string, data []string) int
}

type producerMQ struct {
    obj IProducer
}

func (p *producerMQ) SendTopicData(topic string, data []string) int {
    return p.obj.SendTopicData(topic, data)
}

var (
    insProducer *producerMQ
)

func init() {
    insProducer = &producerMQ{}

    conf := mpf.NewConfig().GetConfig("mpmq")
    mqType := conf.GetString("common." + mpf.EnvProjectKey() + ".mqtype")
    switch mqType {
    case mqTypeRabbit:
        pr := newProducerRabbit()
        pr.connect(conf)
        insProducer.obj = pr
    case mqTypeRedis:
        pr := newProducerRedis()
        insProducer.obj = pr
    default:
        panic(mperr.NewMQ(errorcode.MQParam, "消息队列类型不支持", nil))
    }
}

// NewProducer NewProducer
func NewProducer() *producerMQ {
    return insProducer
}
