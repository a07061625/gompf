/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/3 0003
 * Time: 18:59
 */
package push

import (
    "strconv"
    "time"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/mppush"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 推送广播消息
type pushAll struct {
    mppush.BaseBaiDu
    msgContent map[string]interface{} // 消息内容
}

func (pa *pushAll) SetMsgContent(msgContent map[string]interface{}) {
    if len(msgContent) > 0 {
        pa.msgContent = msgContent
    } else {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "消息内容不合法", nil))
    }
}

func (pa *pushAll) SetMsgType(msgType int) {
    if (msgType == 0) || (msgType == 1) {
        pa.ReqData["msg_type"] = strconv.Itoa(msgType)
    } else {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "消息类型不合法", nil))
    }
}

func (pa *pushAll) SetMsgExpires(msgExpires int) {
    if (msgExpires >= 0) && (msgExpires <= 604800) {
        pa.ReqData["msg_expires"] = strconv.Itoa(msgExpires)
    } else {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "消息过期时间不合法", nil))
    }
}

func (pa *pushAll) SetDeployStatus(deployStatus int) {
    if (deployStatus == 1) || (deployStatus == 2) {
        pa.ReqData["deploy_status"] = strconv.Itoa(deployStatus)
    } else {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "部署状态不合法", nil))
    }
}

func (pa *pushAll) SetSendTime(sendTime int) {
    if (sendTime - time.Now().Second()) > 60 {
        pa.ReqData["send_time"] = strconv.Itoa(sendTime)
    } else {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "发送时间不合法", nil))
    }
}

func (pa *pushAll) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(pa.msgContent) == 0 {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "消息内容不能为空", nil))
    }
    pa.ReqData["msg"] = mpf.JsonMarshal(pa.msgContent)

    return pa.GetRequest()
}

func NewPushAll() *pushAll {
    pa := &pushAll{mppush.NewBaseBaiDu(), make(map[string]interface{})}
    pa.ServiceUri = "/push/all"
    pa.ReqData["msg_type"] = "0"
    pa.ReqData["msg_expires"] = "18000"
    pa.ReqData["deploy_status"] = "1"
    return pa
}
