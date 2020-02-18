/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/1 0001
 * Time: 21:57
 */
package auth

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/mpiot"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 鉴权
type authorize struct {
    mpiot.BaseBaiDu
    principalUuid string // 用户uuid
    action        string // 操作
    topic         string // 主题名
}

func (a *authorize) SetPrincipalUuid(principalUuid string) {
    if len(principalUuid) > 0 {
        a.principalUuid = principalUuid
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "用户uuid不合法", nil))
    }
}

func (a *authorize) SetActionAndTopic(action, topic string) {
    switch action {
    case "CONNECT":
        a.action = action
        a.topic = topic
    case "CREATE":
        a.action = action
        a.topic = topic
    case "SEND":
        if len(topic) == 0 {
            panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "主题名不合法", nil))
        }
        a.action = action
        a.topic = topic
    case "RECEIVE":
        if len(topic) == 0 {
            panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "主题名不合法", nil))
        }
        a.action = action
        a.topic = topic
    case "CONSUME":
        if len(topic) == 0 {
            panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "主题名不合法", nil))
        }
        a.action = action
        a.topic = topic
    default:
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "操作不合法", nil))
    }
}

func (a *authorize) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(a.principalUuid) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "用户uuid不能为空", nil))
    }
    if len(a.action) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "操作不能为空", nil))
    }
    a.ExtendData["principalUuid"] = a.principalUuid
    a.ExtendData["action"] = a.action
    a.ExtendData["topic"] = a.topic

    a.ReqURI = a.GetServiceUrl()

    reqBody := mpf.JSONMarshal(a.ExtendData)
    client, req := a.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewAuthorize() *authorize {
    a := &authorize{mpiot.NewBaseBaiDu(), "", "", ""}
    a.ServiceUri = "/v1/auth/authorize"
    a.ReqContentType = project.HTTPContentTypeJSON
    a.ReqMethod = fasthttp.MethodPost
    return a
}
