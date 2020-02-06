/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/5 0005
 * Time: 15:11
 */
package mpmq

import (
    "strconv"
    "strings"
    "time"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/spf13/viper"
    "github.com/streadway/amqp"
)

type consumerRabbit struct {
    baseRabbit
    conn         *amqp.Connection
    channel      *amqp.Channel
    queueName    string
    consumerName string
}

func (c *consumerRabbit) connect(conf *viper.Viper) {
    conn, err := amqp.Dial(conf.GetString("consumer." + mpf.EnvProjectKey() + ".rabbit.uri"))
    if err != nil {
        panic(mperr.NewMQRabbit(errorcode.MQRabbitConnect, "rabbit连接出错", err))
    }

    channel, err := conn.Channel()
    if err != nil {
        panic(mperr.NewMQRabbit(errorcode.MQRabbitConnect, "rabbit出错", err))
    }

    if err := channel.ExchangeDeclare(
        c.ExchangeName,
        amqp.ExchangeTopic,
        true,
        false,
        false,
        false,
        nil,
    ); err != nil {
        panic(mperr.NewMQRabbit(errorcode.MQRabbitConnect, "rabbit出错", err))
    }

    nonceStr := mpf.ToolCreateNonceStr(6, "numlower") + strconv.FormatInt(time.Now().Unix(), 10)
    queue, err := channel.QueueDeclare(
        "mpqueue"+mpf.EnvProjectKey()+nonceStr,
        true,
        false,
        false,
        false,
        nil,
    )
    if err != nil {
        panic(mperr.NewMQRabbit(errorcode.MQRabbitConnect, "rabbit出错", err))
    }

    if err = channel.QueueBind(
        queue.Name,
        c.TopicPrefix+"*",
        c.ExchangeName,
        false,
        nil,
    ); err != nil {
        panic(mperr.NewMQRabbit(errorcode.MQRabbitConnect, "rabbit出错", err))
    }

    c.conn = conn
    c.channel = channel
    c.queueName = queue.Name
    c.consumerName = "mpconsumer" + mpf.EnvProjectKey() + nonceStr
}

func (c *consumerRabbit) PullData(topic string) (interface{}, error) {
    deliveries, err := c.channel.Consume(
        c.queueName,
        c.consumerName,
        false,
        false,
        false,
        false,
        nil,
    )
    if err != nil {
        return nil, err
    }

    deliveryList := make(map[string][]amqp.Delivery)
    for delivery := range deliveries {
        topicName := strings.TrimPrefix(delivery.RoutingKey, c.TopicPrefix)
        _, ok := deliveryList[topicName]
        if !ok {
            deliveryList[topicName] = make([]amqp.Delivery, 0)
        }
        deliveryList[topicName] = append(deliveryList[topicName], delivery)
    }
    return deliveryList, nil
}

func (c *consumerRabbit) Shutdown() int {
    if err := c.channel.Cancel(c.consumerName, true); err != nil {
        mpf.NewLogger().Error("Consumer cancel failed: " + err.Error())
        return 1
    }
    if err := c.conn.Close(); err != nil {
        mpf.NewLogger().Error("AMQP connection close error: " + err.Error())
        return 1
    }
    return 0
}

func newConsumerRabbit() *consumerRabbit {
    return &consumerRabbit{newBaseRabbit(), nil, nil, "", ""}
}
