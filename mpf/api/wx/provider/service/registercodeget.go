/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/9 0009
 * Time: 23:16
 */
package service

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

// 获取注册码
type registerCodeGet struct {
    wx.BaseWxProvider
    templateId  string // 推广包ID
    corpName    string // 企业名称
    adminName   string // 管理员姓名
    adminMobile string // 管理员手机号
    state       string // 自定义状态值
    followUser  string // 跟进人用户ID
}

func (rcg *registerCodeGet) SetTemplateId(templateId string) {
    if (len(templateId) > 0) && (len(templateId) <= 128) {
        rcg.templateId = templateId
    } else {
        panic(mperr.NewWxProvider(errorcode.WxProviderParam, "推广包ID不合法", nil))
    }
}

func (rcg *registerCodeGet) SetCorpName(corpName string) {
    rcg.corpName = corpName
}

func (rcg *registerCodeGet) SetAdminName(adminName string) {
    rcg.adminName = adminName
}

func (rcg *registerCodeGet) SetAdminMobile(adminMobile string) {
    match, _ := regexp.MatchString(`^[0-9]{11}$`, adminMobile)
    if match && (adminMobile[:1] == "1") {
        rcg.adminMobile = adminMobile
    } else {
        panic(mperr.NewWxProvider(errorcode.WxProviderParam, "管理员手机号不合法", nil))
    }
}

func (rcg *registerCodeGet) SetState(state string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{1,128}$`, state)
    if match {
        rcg.state = state
    } else {
        panic(mperr.NewWxProvider(errorcode.WxProviderParam, "自定义状态值不合法", nil))
    }
}

func (rcg *registerCodeGet) SetFollowUser(followUser string) {
    rcg.followUser = followUser
}

func (rcg *registerCodeGet) checkData() {
    if len(rcg.templateId) == 0 {
        panic(mperr.NewWxProvider(errorcode.WxProviderParam, "推广包ID不能为空", nil))
    }
}

func (rcg *registerCodeGet) SendRequest() api.ApiResult {
    rcg.checkData()

    rcg.ReqData["template_id"] = rcg.templateId
    rcg.ReqData["state"] = rcg.state
    if len(rcg.corpName) > 0 {
        rcg.ReqData["corp_name"] = rcg.corpName
    }
    if len(rcg.adminName) > 0 {
        rcg.ReqData["admin_name"] = rcg.adminName
    }
    if len(rcg.adminMobile) > 0 {
        rcg.ReqData["admin_mobile"] = rcg.adminMobile
    }
    if len(rcg.followUser) > 0 {
        rcg.ReqData["follow_user"] = rcg.followUser
    }
    reqBody := mpf.JsonMarshal(rcg.ReqData)
    rcg.ReqUrl = "https://qyapi.weixin.qq.com/cgi-bin/service/get_register_code?provider_access_token=" + wx.NewUtilWx().GetProviderToken()
    client, req := rcg.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := rcg.SendInner(client, req, errorcode.WxProviderRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxProviderRequestPost
        result.Msg = respData["errmsg"].(string)
    }
    return result
}

func NewRegisterCodeGet() *registerCodeGet {
    rcg := &registerCodeGet{wx.NewBaseWxProvider(), "", "", "", "", "", ""}
    rcg.state = mpf.ToolCreateNonceStr(8, "numlower")
    rcg.ReqContentType = project.HTTPContentTypeJSON
    rcg.ReqMethod = fasthttp.MethodPost
    return rcg
}
