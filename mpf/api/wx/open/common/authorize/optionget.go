/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/10 0010
 * Time: 12:02
 */
package authorize

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type optionGet struct {
    wx.BaseWxOpen
    optionName string // 选项名称
}

func (og *optionGet) SetOptionName(optionName string) {
    if len(optionName) > 0 {
        og.optionName = optionName
    } else {
        panic(mperr.NewWxOpenAccount(errorcode.WxOpenParam, "选项名称不合法", nil))
    }
}

func (og *optionGet) checkData() {
    if len(og.optionName) == 0 {
        panic(mperr.NewWxOpenAccount(errorcode.WxOpenParam, "选项名称不能为空", nil))
    }
    og.ReqData["option_name"] = og.optionName
}

func (og *optionGet) SendRequest() api.APIResult {
    og.checkData()

    reqBody := mpf.JSONMarshal(og.ReqData)
    og.ReqURI = "https://api.weixin.qq.com/cgi-bin/component/api_get_authorizer_option?component_access_token=" + wx.NewUtilWx().GetOpenAccessToken()
    client, req := og.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := og.SendInner(client, req, errorcode.WxOpenRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    _, ok := respData["authorizer_appid"]
    if ok {
        result.Data = respData
    } else {
        result.Code = errorcode.WxOpenRequestPost
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewOptionGet(appId string) *optionGet {
    conf := wx.NewConfig().GetOpen()
    og := &optionGet{wx.NewBaseWxOpen(), ""}
    og.ReqData["component_appid"] = conf.GetAppId()
    og.ReqData["authorizer_appid"] = appId
    og.ReqContentType = project.HTTPContentTypeJSON
    og.ReqMethod = fasthttp.MethodPost
    return og
}
