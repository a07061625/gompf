/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/5 0005
 * Time: 15:12
 */
package mpmq

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/spf13/viper"
    "github.com/streadway/amqp"
)

type producerRabbit struct {
    baseRabbit
    channel *amqp.Channel
}

func (p *producerRabbit) connect(conf *viper.Viper) {
    conn, err := amqp.Dial(conf.GetString("producer." + mpf.EnvProjectKey() + ".rabbit.uri"))
    if err != nil {
        panic(mperr.NewMQRabbit(errorcode.MQRabbitConnect, "rabbit连接出错", err))
    }
    defer conn.Close()

    channel, err := conn.Channel()
    if err != nil {
        panic(mperr.NewMQRabbit(errorcode.MQRabbitConnect, "rabbit出错", err))
    }

    if err := channel.ExchangeDeclare(
        p.ExchangeName,
        amqp.ExchangeTopic,
        true,
        false,
        false,
        false,
        nil,
    ); err != nil {
        panic(mperr.NewMQRabbit(errorcode.MQRabbitConnect, "rabbit出错", err))
    }

    p.channel = channel
}

func (p *producerRabbit) SendTopicData(topic string, data []string) int {
    num := 0
    routingKey := p.TopicPrefix + topic
    for _, v := range data {
        err := p.channel.Publish(
            p.ExchangeName,
            routingKey,
            false,
            false,
            amqp.Publishing{
                Headers:         amqp.Table{},
                ContentType:     "text/plain",
                ContentEncoding: "",
                Body:            []byte(v),
                DeliveryMode:    amqp.Transient,
                Priority:        0,
            },
        )
        if err != nil {
            num++
        }
    }

    return num
}

func newProducerRabbit() *producerRabbit {
    return &producerRabbit{newBaseRabbit(), nil}
}
