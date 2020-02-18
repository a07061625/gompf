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

type optionSet struct {
    wx.BaseWxOpen
    optionName  string // 选项名称
    optionValue string // 选项值
}

func (os *optionSet) SetOptionName(optionName string) {
    if len(optionName) > 0 {
        os.optionName = optionName
    } else {
        panic(mperr.NewWxOpenAccount(errorcode.WxOpenParam, "选项名称不合法", nil))
    }
}

func (os *optionSet) SetOptionValue(optionValue string) {
    if len(optionValue) > 0 {
        os.optionValue = optionValue
    } else {
        panic(mperr.NewWxOpenAccount(errorcode.WxOpenParam, "选项值不合法", nil))
    }
}

func (os *optionSet) checkData() {
    if len(os.optionName) == 0 {
        panic(mperr.NewWxOpenAccount(errorcode.WxOpenParam, "选项名称不能为空", nil))
    }
    if len(os.optionValue) == 0 {
        panic(mperr.NewWxOpenAccount(errorcode.WxOpenParam, "选项值不能为空", nil))
    }
    os.ReqData["option_name"] = os.optionName
    os.ReqData["option_value"] = os.optionValue
}

func (os *optionSet) SendRequest() api.ApiResult {
    os.checkData()

    reqBody := mpf.JsonMarshal(os.ReqData)
    os.ReqUrl = "https://api.weixin.qq.com/cgi-bin/component/api_set_authorizer_option?component_access_token=" + wx.NewUtilWx().GetOpenAccessToken()
    client, req := os.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := os.SendInner(client, req, errorcode.WxOpenRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxOpenRequestPost
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewOptionSet(appId string) *optionSet {
    conf := wx.NewConfig().GetOpen()
    os := &optionSet{wx.NewBaseWxOpen(), "", ""}
    os.ReqData["component_appid"] = conf.GetAppId()
    os.ReqData["authorizer_appid"] = appId
    os.ReqContentType = project.HTTPContentTypeJSON
    os.ReqMethod = fasthttp.MethodPost
    return os
}
