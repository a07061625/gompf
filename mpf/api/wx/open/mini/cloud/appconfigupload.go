/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/11 0011
 * Time: 8:52
 */
package cloud

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 上传小程序配置
type appConfigUpload struct {
    wx.BaseWxOpen
    appId       string                 // 应用ID
    configType  int                    // 配置类型
    configValue map[string]interface{} // 配置json
}

func (acu *appConfigUpload) SetConfigInfo(configType int, configValue map[string]interface{}) {
    if configType != 1 {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "配置类型不合法", nil))
    }
    if len(configValue) == 0 {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "配置不能为空", nil))
    }
    acu.configType = configType
    acu.configValue = configValue
}

func (acu *appConfigUpload) checkData() {
    if acu.configType <= 0 {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "配置类型不能为空", nil))
    }
}

func (acu *appConfigUpload) SendRequest() api.ApiResult {
    acu.checkData()

    reqData := make(map[string]interface{})
    reqData["type"] = acu.configType
    reqData["config"] = mpf.JsonMarshal(acu.configValue)
    reqBody := mpf.JsonMarshal(reqData)
    acu.ReqUrl = "https://api.weixin.qq.com/tcb/uploadappconfig?access_token=" + wx.NewUtilWx().GetOpenAuthorizeAccessToken(acu.appId)
    client, req := acu.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := acu.SendInner(client, req, errorcode.WxOpenRequestPost)
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

func NewAppConfigUpload(appId string) *appConfigUpload {
    acu := &appConfigUpload{wx.NewBaseWxOpen(), "", 0, make(map[string]interface{})}
    acu.appId = appId
    acu.configType = 0
    acu.ReqContentType = project.HttpContentTypeJson
    acu.ReqMethod = fasthttp.MethodPost
    return acu
}
