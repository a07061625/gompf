/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/3 0003
 * Time: 18:59
 */
package push

import (
    "regexp"
    "strconv"
    "time"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/mppush"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 推送组播消息
type pushTags struct {
    mppush.BaseBaiDu
    msgContent map[string]interface{} // 消息内容
    tag        string                 // 标签名
}

func (pt *pushTags) SetMsgContent(msgContent map[string]interface{}) {
    if len(msgContent) > 0 {
        pt.msgContent = msgContent
    } else {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "消息内容不合法", nil))
    }
}

func (pt *pushTags) SetTag(tag string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, tag)
    if match {
        pt.tag = tag
    } else {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "标签名不合法", nil))
    }
}

func (pt *pushTags) SetMsgType(msgType int) {
    if (msgType == 0) || (msgType == 1) {
        pt.ReqData["msg_type"] = strconv.Itoa(msgType)
    } else {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "消息类型不合法", nil))
    }
}

func (pt *pushTags) SetMsgExpires(msgExpires int) {
    if (msgExpires >= 0) && (msgExpires <= 604800) {
        pt.ReqData["msg_expires"] = strconv.Itoa(msgExpires)
    } else {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "消息过期时间不合法", nil))
    }
}

func (pt *pushTags) SetDeployStatus(deployStatus int) {
    if (deployStatus == 1) || (deployStatus == 2) {
        pt.ReqData["deploy_status"] = strconv.Itoa(deployStatus)
    } else {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "部署状态不合法", nil))
    }
}

func (pt *pushTags) SetSendTime(sendTime int) {
    if (sendTime - time.Now().Second()) > 60 {
        pt.ReqData["send_time"] = strconv.Itoa(sendTime)
    } else {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "发送时间不合法", nil))
    }
}

func (pt *pushTags) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(pt.msgContent) == 0 {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "消息内容不能为空", nil))
    }
    if len(pt.tag) == 0 {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "标签名不能为空", nil))
    }
    pt.ReqData["msg"] = mpf.JsonMarshal(pt.msgContent)
    pt.ReqData["tag"] = pt.tag

    return pt.GetRequest()
}

func NewPushTags() *pushTags {
    pt := &pushTags{mppush.NewBaseBaiDu(), make(map[string]interface{}), ""}
    pt.ServiceUri = "/push/tags"
    pt.ReqData["type"] = "1"
    pt.ReqData["msg_type"] = "0"
    pt.ReqData["msg_expires"] = "18000"
    pt.ReqData["deploy_status"] = "1"
    return pt
}
