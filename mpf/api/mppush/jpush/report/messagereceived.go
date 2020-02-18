/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/4 0004
 * Time: 16:13
 */
package report

import (
    "strconv"
    "strings"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/mppush"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 送达统计
type messageReceived struct {
    mppush.BaseJPush
    msgIdList []string // 消息ID列表
}

func (mr *messageReceived) SetMsgIdList(msgIdList []int64) {
    if len(msgIdList) == 0 {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "消息ID列表不能为空", nil))
    } else if len(msgIdList) > 100 {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "消息ID列表不能超过100个", nil))
    }
    mr.msgIdList = make([]string, 0)
    for _, v := range msgIdList {
        if v > 0 {
            mr.msgIdList = append(mr.msgIdList, strconv.FormatInt(v, 10))
        }
    }
}

func (mr *messageReceived) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(mr.msgIdList) == 0 {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "消息ID列表不能为空", nil))
    }
    mr.ReqData["msg_ids"] = strings.Join(mr.msgIdList, ",")

    mr.ReqURI = mr.GetServiceUrl() + "?" + mpf.HTTPCreateParams(mr.ReqData, "none", 1)

    return mr.GetRequest()
}

func NewMessageReceived(key string) *messageReceived {
    mr := &messageReceived{mppush.NewBaseJPush(mppush.JPushServiceDomainReport, key, "app"), make([]string, 0)}
    mr.ServiceUri = "/v3/received"
    return mr
}
