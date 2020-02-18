/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/6 0006
 * Time: 22:13
 */
package authorize

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
)

// 获取网页登录服务商授权引导地址
type loginAuthUrlWeb struct {
    wx.BaseWxProvider
    redirectUri  string // 授权回调地址
    responseType string // 返回类型
    scope        string // 授权作用域 snsapi_base:静默授权,可获取成员的基础信息(UserId与DeviceId) snsapi_userinfo:静默授权,可获取成员的详细信息,但不包含手机、邮箱等敏感信息 snsapi_privateinfo:手动授权,可获取成员的详细信息,包含手机、邮箱等敏感信息
    state        string // 防跨域攻击标识
}

func (luw *loginAuthUrlWeb) SetState(state string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{1,128}$`, state)
    if match {
        luw.ReqData["state"] = state
    } else {
        panic(mperr.NewWxProvider(errorcode.WxProviderParam, "防跨域攻击标识不合法", nil))
    }
}

func (luw *loginAuthUrlWeb) SetScope(scope string) {
    if (scope == "snsapi_base") || (scope == "snsapi_userinfo") || (scope == "snsapi_privateinfo") {
        luw.scope = scope
    } else {
        panic(mperr.NewWxProvider(errorcode.WxProviderParam, "授权作用域不合法", nil))
    }
}

func (luw *loginAuthUrlWeb) checkData() {
    if len(luw.scope) == 0 {
        panic(mperr.NewWxProvider(errorcode.WxProviderParam, "授权作用域不能为空", nil))
    }
    luw.ReqData["scope"] = luw.scope
}

func (luw *loginAuthUrlWeb) GetResult() map[string]string {
    luw.checkData()

    result := make(map[string]string)
    result["url"] = "https://open.weixin.qq.com/connect/oauth2/authorize?" + mpf.HTTPCreateParams(luw.ReqData, "none", 1) + "#wechat_redirect"
    return result
}

func NewLoginAuthUrlWeb() *loginAuthUrlWeb {
    conf := wx.NewConfig().GetProvider()
    luw := &loginAuthUrlWeb{wx.NewBaseWxProvider(), "", "", "", ""}
    luw.ReqData["appid"] = conf.GetSuiteId()
    luw.ReqData["redirect_uri"] = conf.GetUrlAuthLogin()
    luw.ReqData["response_type"] = "code"
    luw.ReqData["state"] = mpf.ToolCreateNonceStr(8, "numlower")
    return luw
}
