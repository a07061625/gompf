/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/19 0019
 * Time: 11:22
 */
package auth

import (
    "regexp"

    "github.com/a07061625/gompf/mpf/api/alipay"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 用户登陆授权
type userInfo struct {
    alipay.BaseAliPay
    scopes []string // 授权类型
    state  string   // 校验码
}

func (ui *userInfo) SetScopes(scope string) {
    if (scope == "auth_user") || (scope == "auth_base") {
        ui.scopes = make([]string, 0)
        ui.scopes = append(ui.scopes, scope)
    } else {
        panic(mperr.NewAliPayAuth(errorcode.AliPayAuthParam, "授权类型不合法", nil))
    }
}

func (ui *userInfo) SetState(state string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{1,100}$`, state)
    if match {
        ui.state = state
    } else {
        panic(mperr.NewAliPayAuth(errorcode.AliPayAuthParam, "校验码不合法", nil))
    }
}

func (ui *userInfo) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ui.scopes) == 0 {
        panic(mperr.NewAliPayAuth(errorcode.AliPayAuthParam, "授权类型不能为空", nil))
    }
    if len(ui.state) == 0 {
        panic(mperr.NewAliPayAuth(errorcode.AliPayAuthParam, "校验码不能为空", nil))
    }
    ui.BizContent["scopes"] = ui.scopes
    ui.BizContent["state"] = ui.state

    return ui.GetRequest()
}

func NewUserInfo(appId string) *userInfo {
    ui := &userInfo{alipay.NewBase(appId), make([]string, 0), ""}
    ui.SetMethod("alipay.user.info.auth")
    return ui
}
