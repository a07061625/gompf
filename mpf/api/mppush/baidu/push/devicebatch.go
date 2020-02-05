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

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/mppush"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 推送消息到给定的一组设备
type deviceBatch struct {
    mppush.BaseBaiDu
    msgContent  map[string]interface{} // 消息内容
    channelList []string               // 设备列表
}

func (db *deviceBatch) SetMsgContent(msgContent map[string]interface{}) {
    if len(msgContent) > 0 {
        db.msgContent = msgContent
    } else {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "消息内容不合法", nil))
    }
}

func (db *deviceBatch) SetChannelList(channelList []string) {
    if len(channelList) == 0 {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "设备列表不能为空", nil))
    } else if len(channelList) > 10000 {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "设备列表不能超过10000个", nil))
    }
    db.channelList = make([]string, 0)
    for _, v := range channelList {
        if len(v) > 0 {
            db.channelList = append(db.channelList, v)
        }
    }
}

func (db *deviceBatch) SetMsgType(msgType int) {
    if (msgType == 0) || (msgType == 1) {
        db.ReqData["msg_type"] = strconv.Itoa(msgType)
    } else {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "消息类型不合法", nil))
    }
}

func (db *deviceBatch) SetMsgExpires(msgExpires int) {
    if (msgExpires >= 1) && (msgExpires <= 86400) {
        db.ReqData["msg_expires"] = strconv.Itoa(msgExpires)
    } else {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "消息过期时间不合法", nil))
    }
}

func (db *deviceBatch) SetTopicId(topicId string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{1,128}$`, topicId)
    if match {
        db.ReqData["topic_id"] = topicId
    } else {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "分类主题标识不合法", nil))
    }
}

func (db *deviceBatch) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(db.msgContent) == 0 {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "消息内容不能为空", nil))
    }
    if len(db.channelList) == 0 {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "设备列表不能为空", nil))
    }
    db.ReqData["msg"] = mpf.JsonMarshal(db.msgContent)
    db.ReqData["channel_ids"] = mpf.JsonMarshal(db.channelList)

    return db.GetRequest()
}

func NewDeviceBatch() *deviceBatch {
    db := &deviceBatch{mppush.NewBaseBaiDu(), make(map[string]interface{}), make([]string, 0)}
    db.ServiceUri = "/push/batch_device"
    db.ReqData["msg_type"] = "0"
    db.ReqData["msg_expires"] = "86400"
    return db
}
