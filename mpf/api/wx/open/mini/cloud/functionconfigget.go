/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/11 0011
 * Time: 9:41
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

// 获取云函数配置
type functionConfigGet struct {
    wx.BaseWxOpen
    appId        string // 应用ID
    configType   int    // 配置类型
    env          string // 环境id
    functionName string // 云函数名
}

func (fcg *functionConfigGet) SetConfigType(configType int) {
    if (configType == 1) || (configType == 2) {
        fcg.configType = configType
    } else {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "配置类型不合法", nil))
    }
}

func (fcg *functionConfigGet) SetEnv(env string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, env)
    if match {
        fcg.env = env
    } else {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "环境id不合法", nil))
    }
}

func (fcg *functionConfigGet) SetFunctionName(functionName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, functionName)
    if match {
        fcg.functionName = functionName
    } else {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "云函数名不合法", nil))
    }
}

func (fcg *functionConfigGet) checkData() {
    if fcg.configType <= 0 {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "配置类型不能为空", nil))
    }
    if len(fcg.env) == 0 {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "环境id不能为空", nil))
    }
    if len(fcg.functionName) == 0 {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "云函数名不能为空", nil))
    }
}

func (fcg *functionConfigGet) SendRequest() api.ApiResult {
    fcg.checkData()

    reqData := make(map[string]interface{})
    reqData["type"] = fcg.configType
    reqData["env"] = fcg.env
    reqData["function_name"] = fcg.functionName
    reqBody := mpf.JSONMarshal(fcg.ReqData)
    fcg.ReqUrl = "https://api.weixin.qq.com/tcb/getfuncconfig?access_token=" + wx.NewUtilWx().GetOpenAuthorizeAccessToken(fcg.appId)
    client, req := fcg.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := fcg.SendInner(client, req, errorcode.WxOpenRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxOpenRequestPost
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewFunctionConfigGet(appId string) *functionConfigGet {
    fcg := &functionConfigGet{wx.NewBaseWxOpen(), "", 0, "", ""}
    fcg.appId = appId
    fcg.configType = 0
    fcg.ReqContentType = project.HTTPContentTypeJSON
    fcg.ReqMethod = fasthttp.MethodPost
    return fcg
}
