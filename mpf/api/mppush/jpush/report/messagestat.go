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

// 消息统计
type messageStat struct {
    mppush.BaseJPush
    msgIdList []string // 消息ID列表
}

func (ms *messageStat) SetMsgIdList(msgIdList []int64) {
    if len(msgIdList) == 0 {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "消息ID列表不能为空", nil))
    } else if len(msgIdList) > 100 {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "消息ID列表不能超过100个", nil))
    }
    ms.msgIdList = make([]string, 0)
    for _, v := range msgIdList {
        if v > 0 {
            ms.msgIdList = append(ms.msgIdList, strconv.FormatInt(v, 10))
        }
    }
}

func (ms *messageStat) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ms.msgIdList) == 0 {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "消息ID列表不能为空", nil))
    }
    ms.ReqData["msg_ids"] = strings.Join(ms.msgIdList, ",")

    ms.ReqUrl = ms.GetServiceUrl() + "?" + mpf.HttpCreateParams(ms.ReqData, "none", 1)

    return ms.GetRequest()
}

func NewMessageStat(key string) *messageStat {
    ms := &messageStat{mppush.NewBaseJPush(mppush.JPushServiceDomainReport, key, "app"), make([]string, 0)}
    ms.ServiceUri = "/v3/messages"
    return ms
}
