/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/9 0009
 * Time: 23:31
 */
package service

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 查询注册状态
type registerInfoGet struct {
    wx.BaseWxProvider
    registerCode string // 注册码
}

func (rig *registerInfoGet) SetRegisterCode(registerCode string) {
    if len(registerCode) > 0 {
        rig.registerCode = registerCode
    } else {
        panic(mperr.NewWxProvider(errorcode.WxProviderParam, "注册码不合法", nil))
    }
}

func (rig *registerInfoGet) checkData() {
    if len(rig.registerCode) == 0 {
        panic(mperr.NewWxProvider(errorcode.WxProviderParam, "注册码不能为空", nil))
    }
}

func (rig *registerInfoGet) SendRequest() api.APIResult {
    rig.checkData()

    rig.ReqData["register_code"] = rig.registerCode
    reqBody := mpf.JSONMarshal(rig.ReqData)
    rig.ReqURI = "https://qyapi.weixin.qq.com/cgi-bin/service/get_register_info?provider_access_token=" + wx.NewUtilWx().GetProviderToken()
    client, req := rig.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := rig.SendInner(client, req, errorcode.WxProviderRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxProviderRequestPost
        result.Msg = respData["errmsg"].(string)
    }
    return result
}

func NewRegisterInfoGet() *registerInfoGet {
    rig := &registerInfoGet{wx.NewBaseWxProvider(), ""}
    rig.ReqContentType = project.HTTPContentTypeJSON
    rig.ReqMethod = fasthttp.MethodPost
    return rig
}
