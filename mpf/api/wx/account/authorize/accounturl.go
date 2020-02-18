/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/11 0011
 * Time: 15:19
 */
package authorize

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
)

// 公众号授权地址
type accountUrl struct {
    wx.BaseWxAccount
    appId       string
    redirectUrl string // 重定向链接
    authType    string // 授权类型
    state       string // 防csrf攻击标识
}

func (au *accountUrl) SetRedirectUrl(redirectUrl string) {
    match, _ := regexp.MatchString(project.RegexURLHTTP, redirectUrl)
    if match {
        au.redirectUrl = redirectUrl
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "重定向链接不合法", nil))
    }
}

func (au *accountUrl) SetAuthType(authType string) {
    if (authType == "snsapi_base") || (authType == "snsapi_userinfo") {
        au.authType = authType
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "授权类型不合法", nil))
    }
}

func (au *accountUrl) SetState(state string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{1,32}$`, state)
    if match {
        au.ReqData["state"] = state
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "防csrf攻击标识不合法", nil))
    }
}

func (au *accountUrl) checkData() {
    if len(au.authType) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "授权类型不能为空", nil))
    }
    if len(au.redirectUrl) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "重定向链接不能为空", nil))
    }
    au.ReqData["scope"] = au.authType
    au.ReqData["redirect_uri"] = au.redirectUrl
}

func (au *accountUrl) GetResult() map[string]string {
    au.checkData()

    result := make(map[string]string)
    result["url"] = "https://open.weixin.qq.com/connect/oauth2/authorize?" + mpf.HTTPCreateParams(au.ReqData, "key", 1) + "#wechat_redirect"
    return result
}

func NewAuthorizeUrl(appId string) *accountUrl {
    au := &accountUrl{wx.NewBaseWxAccount(), "", "", "", ""}
    au.appId = appId
    au.ReqData["state"] = "STATE"
    au.ReqData["response_type"] = "code"
    return au
}
