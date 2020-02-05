/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/1 0001
 * Time: 22:57
 */
package client

import (
    "strconv"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/mpiot"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 消息发布
type messagePublish struct {
    mpiot.BaseBaiDu
    topic      string // 主题名称
    msgContent string // 消息内容
    userName   string // 用户名
    password   string // 密码
}

func (mp *messagePublish) SetTopic(topic string) {
    if len(topic) > 0 {
        mp.topic = topic
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "主题名称不合法", nil))
    }
}

func (mp *messagePublish) SetMsgContent(msgContent string) {
    if len(msgContent) > 0 {
        mp.msgContent = msgContent
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "消息内容不合法", nil))
    }
}

func (mp *messagePublish) SetUserName(userName string) {
    if len(userName) > 0 {
        mp.userName = userName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "用户名不合法", nil))
    }
}

func (mp *messagePublish) SetPassword(password string) {
    if len(password) > 0 {
        mp.password = password
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "密码不合法", nil))
    }
}

func (mp *messagePublish) SetQos(qos int) {
    if (qos == 0) || (qos == 1) {
        mp.ReqData["qos"] = strconv.Itoa(qos)
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "消息QoS值不合法", nil))
    }
}

func (mp *messagePublish) SetRetain(retain string) {
    if (retain == "false") || (retain == "true") {
        mp.ReqData["retain"] = retain
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "保留消息标记不合法", nil))
    }
}

func (mp *messagePublish) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(mp.topic) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "主题名称不能为空", nil))
    }
    if len(mp.msgContent) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "消息内容不能为空", nil))
    }
    if len(mp.userName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "用户名不能为空", nil))
    }
    if len(mp.password) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "密码不能为空", nil))
    }
    mp.ReqData["topic"] = mp.topic
    mp.ReqHeader["auth.username"] = mp.userName
    mp.ReqHeader["auth.password"] = mp.password

    mp.ReqUrl = mp.GetServiceUrl() + "?" + mpf.HttpCreateParams(mp.ReqData, "none", 1)

    client, req := mp.GetRequest()
    req.SetBody([]byte(mp.msgContent))

    return client, req
}

func NewMessagePublish() *messagePublish {
    mp := &messagePublish{mpiot.NewBaseBaiDu(), "", "", "", ""}
    mp.ServiceUri = "/v1/proxy"
    mp.ReqData["qos"] = "0"
    mp.ReqData["retain"] = "false"
    mp.ReqHeader["Content-Type"] = "application/octet-stream"
    mp.ReqContentType = "application/octet-stream"
    mp.ReqMethod = fasthttp.MethodPost
    mp.SetServiceDomain(mpiot.BaiDuDomainGZMqtt)
    return mp
}
