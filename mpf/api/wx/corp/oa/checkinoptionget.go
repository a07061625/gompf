/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/7 0007
 * Time: 23:57
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

// 获取打卡规则
type checkInOptionGet struct {
    wx.BaseWxCorp
    corpId   string
    agentTag string
    dateTime int      // 日期
    userList []string // 用户列表
}

func (cog *checkInOptionGet) SetDateTime(dateTime int) {
    if (dateTime > 0) && ((dateTime % 86400) == 0) {
        cog.dateTime = dateTime
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "日期不合法", nil))
    }
}

func (cog *checkInOptionGet) SetUserList(userList []string) {
    cog.userList = make([]string, 0)
    for _, v := range userList {
        match, _ := regexp.MatchString(project.RegexDigitAlpha, v)
        if match {
            cog.userList = append(cog.userList, v)
        }
    }
    if len(cog.userList) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "用户列表不能为空", nil))
    } else if len(cog.userList) > 100 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "用户列表不能超过100个", nil))
    }
}

func (cog *checkInOptionGet) checkData() {
    if cog.dateTime <= 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "日期不能为空", nil))
    }
    if len(cog.userList) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "用户列表不能为空", nil))
    }
}

func (cog *checkInOptionGet) SendRequest() api.ApiResult {
    cog.checkData()

    reqData := make(map[string]interface{})
    reqData["datetime"] = cog.dateTime
    reqData["useridlist"] = cog.userList
    reqBody := mpf.JsonMarshal(reqData)

    cog.ReqUrl = "https://qyapi.weixin.qq.com/cgi-bin/checkin/getcheckinoption?access_token=" + wx.NewUtilWx().GetCorpAccessToken(cog.corpId, cog.agentTag)
    client, req := cog.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := cog.SendInner(client, req, errorcode.WxCorpRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxCorpRequestPost
        result.Msg = respData["errmsg"].(string)
    }
    return result
}

func NewCheckInOptionGet(corpId, agentTag string) *checkInOptionGet {
    cog := &checkInOptionGet{wx.NewBaseWxCorp(), "", "", 0, make([]string, 0)}
    cog.corpId = corpId
    cog.agentTag = agentTag
    cog.ReqContentType = project.HTTPContentTypeJSON
    cog.ReqMethod = fasthttp.MethodPost
    return cog
}
