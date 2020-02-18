/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/9 0009
 * Time: 12:58
 */
package user

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

// user id转openid
type userId2Openid struct {
    wx.BaseWxCorp
    corpId   string
    agentTag string
    userId   string // 用户ID
}

func (uo *userId2Openid) SetUserId(userId string) {
    match, _ := regexp.MatchString(project.RegexDigitLower, userId)
    if match {
        uo.userId = userId
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "用户ID不合法", nil))
    }
}

func (uo *userId2Openid) checkData() {
    if len(uo.userId) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "用户ID不能为空", nil))
    }
    uo.ReqData["userid"] = uo.userId
}

func (uo *userId2Openid) SendRequest(getType string) api.ApiResult {
    uo.checkData()

    reqBody := mpf.JsonMarshal(uo.ReqData)
    uo.ReqUrl = "https://qyapi.weixin.qq.com/cgi-bin/user/convert_to_openid?access_token=" + wx.NewUtilWx().GetCorpCache(uo.corpId, uo.agentTag, getType)
    client, req := uo.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := uo.SendInner(client, req, errorcode.WxCorpRequestPost)
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

func NewUserId2Openid(corpId, agentTag string) *userId2Openid {
    uo := &userId2Openid{wx.NewBaseWxCorp(), "", "", ""}
    uo.corpId = corpId
    uo.agentTag = agentTag
    uo.ReqContentType = project.HTTPContentTypeJSON
    uo.ReqMethod = fasthttp.MethodPost
    return uo
}
