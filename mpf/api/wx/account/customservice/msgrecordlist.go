/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/11 0011
 * Time: 22:56
 */
package customservice

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type msgRecordList struct {
    wx.BaseWxAccount
    appId     string
    startTime int // 起始时间
    endTime   int // 结束时间
    msgId     int // 消息id
    number    int // 条数
}

func (mrl *msgRecordList) SetTime(startTime, endTime int) {
    if startTime <= 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "起始时间不合法", nil))
    } else if endTime <= 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "结束时间不合法", nil))
    } else if startTime >= endTime {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "起始时间必须小于结束时间", nil))
    } else if (endTime - startTime) > 86400 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "结束时间不能超过起始时间24小时", nil))
    }
    mrl.startTime = startTime
    mrl.endTime = endTime
}

func (mrl *msgRecordList) SetMsgId(msgId int) {
    if msgId > 0 {
        mrl.msgId = msgId
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "消息id不合法", nil))
    }
}

func (mrl *msgRecordList) SetNumber(number int) {
    if (number > 0) && (number <= 10000) {
        mrl.number = number
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "条数不合法", nil))
    }
}

func (mrl *msgRecordList) checkData() {
    if mrl.startTime <= 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "起始时间不能为空", nil))
    }
}

func (mrl *msgRecordList) SendRequest() api.APIResult {
    mrl.checkData()

    reqData := make(map[string]interface{})
    reqData["starttime"] = mrl.startTime
    reqData["endtime"] = mrl.endTime
    reqData["msgid"] = mrl.msgId
    reqData["number"] = mrl.number
    reqBody := mpf.JSONMarshal(reqData)
    mrl.ReqURI = "https://api.weixin.qq.com/customservice/msgrecord/getmsglist?access_token=" + wx.NewUtilWx().GetSingleAccessToken(mrl.appId)
    client, req := mrl.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := mrl.SendInner(client, req, errorcode.WxAccountRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    _, ok := respData["recordlist"]
    if ok {
        result.Data = respData
    } else {
        result.Code = errorcode.WxAccountRequestPost
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewMsgRecordList(appId string) *msgRecordList {
    mrl := &msgRecordList{wx.NewBaseWxAccount(), "", 0, 0, 0, 0}
    mrl.appId = appId
    mrl.msgId = 1
    mrl.number = 100
    mrl.ReqContentType = project.HTTPContentTypeJSON
    mrl.ReqMethod = fasthttp.MethodPost
    return mrl
}
