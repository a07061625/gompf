/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/7 0007
 * Time: 23:42
 */
package oa

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 获取打卡数据
type checkInDataGet struct {
    wx.BaseWxCorp
    corpId      string
    agentTag    string
    startTime   int      // 开始时间
    endTime     int      // 结束时间
    userList    []string // 用户列表
    checkInType int      // 打卡类型
}

func (cdg *checkInDataGet) SetStartAndEndTime(startTime, endTime int) {
    if startTime <= 1000000000 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "开始时间不合法", nil))
    } else if endTime <= 1000000000 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "结束时间不合法", nil))
    } else if startTime > endTime {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "开始时间不能大于结束时间", nil))
    } else if (endTime - startTime) > 2592000 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "结束时间不能超过开始时间30天", nil))
    }
    cdg.startTime = startTime
    cdg.endTime = endTime
}

func (cdg *checkInDataGet) SetCheckInType(checkInType int) {
    if (checkInType >= 1) && (checkInType <= 3) {
        cdg.checkInType = checkInType
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "打卡类型不合法", nil))
    }
}

func (cdg *checkInDataGet) SetUserList(userList []string) {
    cdg.userList = make([]string, 0)
    for _, v := range userList {
        match, _ := regexp.MatchString(project.RegexDigitAlpha, v)
        if match {
            cdg.userList = append(cdg.userList, v)
        }
    }
    if len(cdg.userList) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "用户列表不能为空", nil))
    } else if len(cdg.userList) > 100 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "用户列表不能超过100个", nil))
    }
}

func (cdg *checkInDataGet) checkData() {
    if cdg.startTime <= 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "开始时间不能为空", nil))
    }
    if cdg.checkInType <= 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "打卡类型不能为空", nil))
    }
    if len(cdg.userList) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "用户列表不能为空", nil))
    }
}

func (cdg *checkInDataGet) SendRequest() api.APIResult {
    cdg.checkData()

    reqData := make(map[string]interface{})
    reqData["opencheckindatatype"] = cdg.checkInType
    reqData["starttime"] = cdg.startTime
    reqData["endtime"] = cdg.endTime
    reqData["useridlist"] = cdg.userList
    reqBody := mpf.JSONMarshal(reqData)

    cdg.ReqURI = "https://qyapi.weixin.qq.com/cgi-bin/checkin/getcheckindata?access_token=" + wx.NewUtilWx().GetCorpAccessToken(cdg.corpId, cdg.agentTag)
    client, req := cdg.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := cdg.SendInner(client, req, errorcode.WxCorpRequestPost)
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

func NewCheckInDataGet(corpId, agentTag string) *checkInDataGet {
    cdg := &checkInDataGet{wx.NewBaseWxCorp(), "", "", 0, 0, make([]string, 0), 0}
    cdg.corpId = corpId
    cdg.agentTag = agentTag
    cdg.ReqContentType = project.HTTPContentTypeJSON
    cdg.ReqMethod = fasthttp.MethodPost
    return cdg
}
