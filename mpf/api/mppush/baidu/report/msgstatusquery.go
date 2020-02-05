/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/3 0003
 * Time: 22:30
 */
package report

import (
    "regexp"

    "github.com/a07061625/gompf/mpf/api/mppush"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 查询消息的发送状态
type msgStatusQuery struct {
    mppush.BaseBaiDu
    msgId string // 消息ID
}

func (msq *msgStatusQuery) SetMsgId(msgId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, msgId)
    if match {
        msq.msgId = msgId
    } else {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "消息ID不合法", nil))
    }
}

func (msq *msgStatusQuery) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(msq.msgId) == 0 {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "消息ID不能为空", nil))
    }
    msq.ReqData["msg_id"] = msq.msgId

    return msq.GetRequest()
}

func NewMsgStatusQuery() *msgStatusQuery {
    msq := &msgStatusQuery{mppush.NewBaseBaiDu(), ""}
    msq.ServiceUri = "/report/query_msg_status"
    return msq
}
