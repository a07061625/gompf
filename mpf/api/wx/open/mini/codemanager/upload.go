/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/11 0011
 * Time: 13:02
 */
package codemanager

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 为授权的小程序帐号上传小程序代码
type upload struct {
    wx.BaseWxOpen
    appId       string                 // 应用ID
    templateId  string                 // 代码模板ID
    userVersion string                 // 自定义代码版本号
    userDesc    string                 // 自定义代码描述
    extData     map[string]interface{} // 自定义配置
}

func (u *upload) SetTemplateId(templateId string) {
    if len(templateId) > 0 {
        u.templateId = templateId
    } else {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "代码模板ID不能为空", nil))
    }
}

func (u *upload) SetUserVersion(userVersion string) {
    if len(userVersion) > 0 {
        u.userVersion = userVersion
    } else {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "自定义代码版本号不能为空", nil))
    }
}

func (u *upload) SetUserDesc(userDesc string) {
    if len(userDesc) > 0 {
        u.userDesc = userDesc
    } else {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "自定义代码描述不能为空", nil))
    }
}

func (u *upload) SetExtData(extData map[string]interface{}) {
    if len(extData) > 0 {
        u.extData = extData
    } else {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "自定义配置不能为空", nil))
    }
}

func (u *upload) checkData() {
    if len(u.templateId) == 0 {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "模板ID不能为空", nil))
    }
    if len(u.extData) == 0 {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "自定义配置不能为空", nil))
    }
}

func (u *upload) SendRequest() api.ApiResult {
    u.checkData()

    reqData := make(map[string]interface{})
    reqData["template_id"] = u.templateId
    reqData["ext_json"] = mpf.JSONMarshal(u.extData)
    if len(u.userVersion) > 0 {
        reqData["user_version"] = u.userVersion
    }
    if len(u.userDesc) > 0 {
        reqData["user_desc"] = u.userDesc
    }
    reqBody := mpf.JSONMarshal(reqData)
    u.ReqUrl = "https://api.weixin.qq.com/wxa/commit?access_token=" + wx.NewUtilWx().GetOpenAuthorizeAccessToken(u.appId)
    client, req := u.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := u.SendInner(client, req, errorcode.WxOpenRequestPost)
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

func NewNicknameAuditQuery(appId string) *upload {
    u := &upload{wx.NewBaseWxOpen(), "", "", "", "", make(map[string]interface{})}
    u.appId = appId
    u.ReqContentType = project.HTTPContentTypeJSON
    u.ReqMethod = fasthttp.MethodPost
    return u
}
