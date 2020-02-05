/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/11 0011
 * Time: 9:53
 */
package cloud

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

// 创建云函数
type functionCreate struct {
    wx.BaseWxOpen
    appId        string // 应用ID
    env          string // 环境id
    functionName string // 云函数名
}

func (fc *functionCreate) SetEnv(env string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, env)
    if match {
        fc.env = env
    } else {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "环境id不合法", nil))
    }
}

func (fc *functionCreate) SetFunctionName(functionName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, functionName)
    if match {
        fc.functionName = functionName
    } else {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "云函数名不合法", nil))
    }
}

func (fc *functionCreate) checkData() {
    if len(fc.env) == 0 {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "环境id不能为空", nil))
    }
    if len(fc.functionName) == 0 {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "云函数名不能为空", nil))
    }
    fc.ReqData["env"] = fc.env
    fc.ReqData["function_name"] = fc.functionName
}

func (fc *functionCreate) SendRequest() api.ApiResult {
    fc.checkData()

    reqBody := mpf.JsonMarshal(fc.ReqData)
    fc.ReqUrl = "https://api.weixin.qq.com/tcb/createfunction?access_token=" + wx.NewUtilWx().GetOpenAuthorizeAccessToken(fc.appId)
    client, req := fc.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := fc.SendInner(client, req, errorcode.WxOpenRequestPost)
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

func NewFunctionCreate(appId string) *functionCreate {
    fc := &functionCreate{wx.NewBaseWxOpen(), "", "", ""}
    fc.appId = appId
    fc.ReqContentType = project.HttpContentTypeJson
    fc.ReqMethod = fasthttp.MethodPost
    return fc
}
