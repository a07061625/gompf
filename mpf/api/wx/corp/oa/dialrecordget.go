/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/8 0008
 * Time: 0:05
 */
package oa

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 获取公费电话拨打记录
type dialRecordGet struct {
    wx.BaseWxCorp
    corpId    string
    agentTag  string
    startTime int // 开始时间
    endTime   int // 结束时间
    offset    int // 偏移量
    limit     int // 每页记录数
}

func (drg *dialRecordGet) SetStartAndEndTime(startTime, endTime int) {
    if startTime <= 1000000000 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "开始时间不合法", nil))
    } else if endTime <= 1000000000 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "结束时间不合法", nil))
    } else if startTime > endTime {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "开始时间不能大于结束时间", nil))
    } else if (endTime - startTime) > 2592000 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "结束时间不能超过开始时间30天", nil))
    }
    drg.startTime = startTime
    drg.endTime = endTime
}

func (drg *dialRecordGet) SetOffset(offset int) {
    if offset >= 0 {
        drg.offset = offset
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "偏移量不合法", nil))
    }
}

func (drg *dialRecordGet) SetLimit(limit int) {
    if (limit > 0) && (limit <= 100) {
        drg.limit = limit
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "每页记录数不合法", nil))
    }
}

func (drg *dialRecordGet) SendRequest() api.APIResult {
    reqData := make(map[string]interface{})
    reqData["offset"] = drg.offset
    reqData["limit"] = drg.limit
    if drg.startTime > 0 {
        reqData["start_time"] = drg.startTime
        reqData["end_time"] = drg.endTime
    }
    reqBody := mpf.JSONMarshal(reqData)

    drg.ReqURI = "https://qyapi.weixin.qq.com/cgi-bin/dial/get_dial_record?access_token=" + wx.NewUtilWx().GetCorpAccessToken(drg.corpId, drg.agentTag)
    client, req := drg.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := drg.SendInner(client, req, errorcode.WxCorpRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxCorpRequestPost
        result.Msg = respData["errmsg"].(string)
    }
    return result
}

func NewDialRecordGet(corpId, agentTag string) *dialRecordGet {
    drg := &dialRecordGet{wx.NewBaseWxCorp(), "", "", 0, 0, 0, 0}
    drg.corpId = corpId
    drg.agentTag = agentTag
    drg.offset = 0
    drg.limit = 20
    drg.ReqContentType = project.HTTPContentTypeJSON
    drg.ReqMethod = fasthttp.MethodPost
    return drg
}
