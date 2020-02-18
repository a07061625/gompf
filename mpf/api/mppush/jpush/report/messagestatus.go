/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/4 0004
 * Time: 12:19
 */
package report

import (
    "regexp"
    "time"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/mppush"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 送达状态查询
type messageStatus struct {
    mppush.BaseJPush
    msgId            string   // 消息ID
    registrationList []string // 设备列表
    dateTime         int64    // 日期时间
}

func (ms *messageStatus) SetMsgId(msgId string) {
    match, _ := regexp.MatchString(project.RegexDigit, msgId)
    if match {
        ms.msgId = msgId
    } else {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "消息ID不合法", nil))
    }
}

func (ms *messageStatus) SetRegistrationList(registrationList []string) {
    if len(registrationList) == 0 {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "设备列表不能为空", nil))
    } else if len(registrationList) > 1000 {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "设备列表不能超过1000个", nil))
    }
    ms.registrationList = make([]string, 0)
    for _, v := range registrationList {
        match, _ := regexp.MatchString(project.RegexDigitAlpha, v)
        if match {
            ms.registrationList = append(ms.registrationList, v)
        }
    }
}

func (ms *messageStatus) SetDateTime(dateTime int64) {
    if dateTime > 0 {
        ms.dateTime = dateTime
    } else {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "日期时间不合法", nil))
    }
}

func (ms *messageStatus) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ms.msgId) == 0 {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "消息ID不能为空", nil))
    }
    if len(ms.registrationList) == 0 {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "设备列表不能为空", nil))
    }
    ms.ExtendData["msg_id"] = ms.msgId
    ms.ExtendData["registration_ids"] = ms.registrationList
    dt := time.Unix(int64(ms.dateTime), 0)
    ms.ExtendData["date"] = dt.Format("2006-01-02")

    ms.ReqURI = ms.GetServiceUrl()

    reqBody := mpf.JSONMarshal(ms.ExtendData)
    client, req := ms.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewMessageStatus(key string) *messageStatus {
    ms := &messageStatus{mppush.NewBaseJPush(mppush.JPushServiceDomainReport, key, "app"), "", make([]string, 0), 0}
    ms.dateTime = time.Now().Unix()
    ms.ServiceUri = "/v3/status/message"
    ms.ReqContentType = project.HTTPContentTypeJSON
    ms.ReqMethod = fasthttp.MethodPost
    return ms
}
