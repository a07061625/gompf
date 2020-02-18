/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/3 0003
 * Time: 18:59
 */
package push

import (
    "strconv"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/mppush"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 推送消息到单台设备
type deviceSingle struct {
    mppush.BaseBaiDu
    msgContent map[string]interface{} // 消息内容
    channelId  string                 // 设备ID
}

func (ds *deviceSingle) SetMsgContent(msgContent map[string]interface{}) {
    if len(msgContent) > 0 {
        ds.msgContent = msgContent
    } else {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "消息内容不合法", nil))
    }
}

func (ds *deviceSingle) SetChannelId(channelId string) {
    if len(channelId) > 0 {
        ds.channelId = channelId
    } else {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "设备ID不合法", nil))
    }
}

func (ds *deviceSingle) SetMsgType(msgType int) {
    if (msgType == 0) || (msgType == 1) {
        ds.ReqData["msg_type"] = strconv.Itoa(msgType)
    } else {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "消息类型不合法", nil))
    }
}

func (ds *deviceSingle) SetMsgExpires(msgExpires int) {
    if (msgExpires >= 0) && (msgExpires <= 604800) {
        ds.ReqData["msg_expires"] = strconv.Itoa(msgExpires)
    } else {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "消息过期时间不合法", nil))
    }
}

func (ds *deviceSingle) SetDeployStatus(deployStatus int) {
    if (deployStatus == 1) || (deployStatus == 2) {
        ds.ReqData["deploy_status"] = strconv.Itoa(deployStatus)
    } else {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "部署状态不合法", nil))
    }
}

func (ds *deviceSingle) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ds.msgContent) == 0 {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "消息内容不能为空", nil))
    }
    if len(ds.channelId) == 0 {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "设备ID不能为空", nil))
    }
    ds.ReqData["msg"] = mpf.JSONMarshal(ds.msgContent)
    ds.ReqData["channel_id"] = ds.channelId

    return ds.GetRequest()
}

func NewDeviceSingle() *deviceSingle {
    ds := &deviceSingle{mppush.NewBaseBaiDu(), make(map[string]interface{}), ""}
    ds.ServiceUri = "/push/single_device"
    ds.ReqData["msg_type"] = "0"
    ds.ReqData["msg_expires"] = "18000"
    ds.ReqData["deploy_status"] = "1"
    return ds
}
