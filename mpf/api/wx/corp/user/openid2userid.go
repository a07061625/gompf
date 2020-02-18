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

// openid转user id
type openid2UserId struct {
    wx.BaseWxCorp
    corpId   string
    agentTag string
    openid   string // 用户openid
}

func (ou *openid2UserId) SetOpenid(openid string) {
    match, _ := regexp.MatchString(project.RegexWxOpenid, openid)
    if match {
        ou.openid = openid
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "用户openid不合法", nil))
    }
}

func (ou *openid2UserId) checkData() {
    if len(ou.openid) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "用户openid不能为空", nil))
    }
    ou.ReqData["openid"] = ou.openid
}

func (ou *openid2UserId) SendRequest(getType string) api.ApiResult {
    ou.checkData()

    reqBody := mpf.JsonMarshal(ou.ReqData)
    ou.ReqUrl = "https://qyapi.weixin.qq.com/cgi-bin/user/convert_to_userid?access_token=" + wx.NewUtilWx().GetCorpCache(ou.corpId, ou.agentTag, getType)
    client, req := ou.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := ou.SendInner(client, req, errorcode.WxCorpRequestPost)
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

func NewOpenid2UserId(corpId, agentTag string) *openid2UserId {
    ou := &openid2UserId{wx.NewBaseWxCorp(), "", "", ""}
    ou.corpId = corpId
    ou.agentTag = agentTag
    ou.ReqContentType = project.HTTPContentTypeJSON
    ou.ReqMethod = fasthttp.MethodPost
    return ou
}
