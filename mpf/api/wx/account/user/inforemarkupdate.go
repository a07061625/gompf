/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/13 0013
 * Time: 0:08
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

type infoRemarkUpdate struct {
    wx.BaseWxAccount
    appId  string
    openid string // 用户openid
    remark string // 备注
}

func (iru *infoRemarkUpdate) SetOpenid(openid string) {
    match, _ := regexp.MatchString(project.RegexWxOpenid, openid)
    if match {
        iru.openid = openid
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "用户openid不合法", nil))
    }
}

func (iru *infoRemarkUpdate) SetRemark(remark string) {
    if (len(remark) > 0) && (len(remark) <= 30) {
        iru.remark = remark
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "备注不合法", nil))
    }
}

func (iru *infoRemarkUpdate) checkData() {
    if len(iru.openid) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "用户openid不能为空", nil))
    }
    if len(iru.remark) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "备注不能为空", nil))
    }
    iru.ReqData["openid"] = iru.openid
    iru.ReqData["remark"] = iru.remark
}

func (iru *infoRemarkUpdate) SendRequest() api.APIResult {
    iru.checkData()

    reqBody := mpf.JSONMarshal(iru.ReqData)
    iru.ReqURI = "https://api.weixin.qq.com/cgi-bin/user/info/updateremark?access_token=" + wx.NewUtilWx().GetSingleAccessToken(iru.appId)
    client, req := iru.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := iru.SendInner(client, req, errorcode.WxAccountRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxAccountRequestPost
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewInfoRemarkUpdate(appId string) *infoRemarkUpdate {
    iru := &infoRemarkUpdate{wx.NewBaseWxAccount(), "", "", ""}
    iru.appId = appId
    iru.ReqContentType = project.HTTPContentTypeJSON
    iru.ReqMethod = fasthttp.MethodPost
    return iru
}
