/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/11 0011
 * Time: 8:48
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

// 获取小程序配置
type appConfigGet struct {
    wx.BaseWxOpen
    appId      string // 应用ID
    configType int    // 配置类型
}

func (acg *appConfigGet) SetConfigType(configType int) {
    if configType == 1 {
        acg.configType = configType
    } else {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "配置类型不合法", nil))
    }
}

func (acg *appConfigGet) checkData() {
    if acg.configType <= 0 {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "配置类型不能为空", nil))
    }
}

func (acg *appConfigGet) SendRequest() api.ApiResult {
    acg.checkData()

    reqData := make(map[string]interface{})
    reqData["type"] = acg.configType
    reqBody := mpf.JSONMarshal(reqData)
    acg.ReqUrl = "https://api.weixin.qq.com/tcb/getappconfig?access_token=" + wx.NewUtilWx().GetOpenAuthorizeAccessToken(acg.appId)
    client, req := acg.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := acg.SendInner(client, req, errorcode.WxOpenRequestPost)
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

func NewAppConfigGet(appId string) *appConfigGet {
    acg := &appConfigGet{wx.NewBaseWxOpen(), "", 0}
    acg.appId = appId
    acg.configType = 0
    acg.ReqContentType = project.HTTPContentTypeJSON
    acg.ReqMethod = fasthttp.MethodPost
    return acg
}
