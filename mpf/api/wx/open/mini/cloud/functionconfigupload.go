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

// 上传云函数配置
type functionConfigUpload struct {
    wx.BaseWxOpen
    appId        string                 // 应用ID
    configType   int                    // 配置类型
    configValue  map[string]interface{} // 配置json
    env          string                 // 环境id
    functionName string                 // 云函数名
}

func (fcu *functionConfigUpload) SetConfigInfo(configType int, configValue map[string]interface{}) {
    if (configType != 1) && (configType != 2) {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "配置类型不合法", nil))
    }
    if len(configValue) == 0 {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "配置不能为空", nil))
    }
    fcu.configType = configType
    fcu.configValue = configValue
}

func (fcu *functionConfigUpload) SetEnv(env string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, env)
    if match {
        fcu.env = env
    } else {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "环境id不合法", nil))
    }
}

func (fcu *functionConfigUpload) SetFunctionName(functionName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, functionName)
    if match {
        fcu.functionName = functionName
    } else {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "云函数名不合法", nil))
    }
}

func (fcu *functionConfigUpload) checkData() {
    if fcu.configType <= 0 {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "配置类型不能为空", nil))
    }
    if len(fcu.env) == 0 {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "环境id不能为空", nil))
    }
    if len(fcu.functionName) == 0 {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "云函数名不能为空", nil))
    }
}

func (fcu *functionConfigUpload) SendRequest() api.ApiResult {
    fcu.checkData()

    reqData := make(map[string]interface{})
    reqData["type"] = fcu.configType
    reqData["config"] = mpf.JsonMarshal(fcu.configValue)
    reqData["env"] = fcu.env
    reqData["function_name"] = fcu.functionName
    reqBody := mpf.JsonMarshal(fcu.ReqData)
    fcu.ReqUrl = "https://api.weixin.qq.com/tcb/uploadfuncconfig?access_token=" + wx.NewUtilWx().GetOpenAuthorizeAccessToken(fcu.appId)
    client, req := fcu.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := fcu.SendInner(client, req, errorcode.WxOpenRequestPost)
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

func NewFunctionConfigUpload(appId string) *functionConfigUpload {
    fcu := &functionConfigUpload{wx.NewBaseWxOpen(), "", 0, make(map[string]interface{}), "", ""}
    fcu.appId = appId
    fcu.configType = 0
    fcu.ReqContentType = project.HttpContentTypeJson
    fcu.ReqMethod = fasthttp.MethodPost
    return fcu
}
