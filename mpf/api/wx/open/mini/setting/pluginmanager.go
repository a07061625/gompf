/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/10 0010
 * Time: 13:02
 */
package setting

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

// 小程序插件管理接口
type pluginManager struct {
    wx.BaseWxOpen
    appId string // 应用ID
    data  map[string]interface{}
}

func (pm *pluginManager) SetData(action string, data map[string]string) {
    switch action {
    case "apply":
        match, _ := regexp.MatchString(project.RegexDigitLower, data["plugin_appid"])
        if !match {
            panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "插件appid不合法", nil))
        }
        pm.data["plugin_appid"] = data["plugin_appid"]
    case "list":
    case "update":
        version, ok := data["user_version"]
        if !ok {
            panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "版本号不能为空", nil))
        } else if len(version) == 0 {
            panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "版本号不能为空", nil))
        }
        match, _ := regexp.MatchString(project.RegexDigitLower, data["plugin_appid"])
        if !match {
            panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "插件appid不合法", nil))
        }

        pm.data["user_version"] = version
        pm.data["plugin_appid"] = data["plugin_appid"]
    case "unbind":
        match, _ := regexp.MatchString(project.RegexDigitLower, data["plugin_appid"])
        if !match {
            panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "插件appid不合法", nil))
        }
        pm.data["plugin_appid"] = data["plugin_appid"]
    default:
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "操作类型不支持", nil))
    }

    pm.data["action"] = action
}

func (pm *pluginManager) checkData() {
    if len(pm.data) == 0 {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "操作类型不支持", nil))
    }
}

func (pm *pluginManager) SendRequest() api.ApiResult {
    pm.checkData()

    reqBody := mpf.JSONMarshal(pm.data)
    pm.ReqUrl = "https://api.weixin.qq.com/wxa/plugin?access_token=" + wx.NewUtilWx().GetOpenAuthorizeAccessToken(pm.appId)
    client, req := pm.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := pm.SendInner(client, req, errorcode.WxOpenRequestPost)
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

func NewPluginManager(appId string) *pluginManager {
    pm := &pluginManager{wx.NewBaseWxOpen(), "", make(map[string]interface{})}
    pm.appId = appId
    pm.ReqContentType = project.HTTPContentTypeJSON
    pm.ReqMethod = fasthttp.MethodPost
    return pm
}
