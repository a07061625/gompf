/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/19 0019
 * Time: 10:51
 */
package auth

import (
    "regexp"

    "github.com/a07061625/gompf/mpf/api/alipay"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 换取应用授权令牌
type authToken struct {
    alipay.BaseAliPay
    grantType    string // 准许类型
    authCode     string // 授权码
    refreshToken string // 刷新令牌
}

func (at *authToken) SetGrantType(grantType string) {
    if (grantType == "authorization_code") || (grantType == "refresh_token") {
        at.grantType = grantType
    } else {
        panic(mperr.NewAliPayAuth(errorcode.AliPayAuthParam, "准许类型不合法", nil))
    }
}

func (at *authToken) SetAuthCode(authCode string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{1,40}$`, authCode)
    if match {
        at.authCode = authCode
    } else {
        panic(mperr.NewAliPayAuth(errorcode.AliPayAuthParam, "授权码不合法", nil))
    }
}

func (at *authToken) SetRefreshToken(refreshToken string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{1,40}$`, refreshToken)
    if match {
        at.refreshToken = refreshToken
    } else {
        panic(mperr.NewAliPayAuth(errorcode.AliPayAuthParam, "刷新令牌不合法", nil))
    }
}

func (at *authToken) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if at.grantType == "authorization_code" {
        if len(at.authCode) == 0 {
            panic(mperr.NewAliPayAuth(errorcode.AliPayAuthParam, "授权码不能为空", nil))
        }
        at.BizContent["code"] = at.authCode
    } else if at.grantType == "refresh_token" {
        if len(at.refreshToken) == 0 {
            panic(mperr.NewAliPayAuth(errorcode.AliPayAuthParam, "刷新令牌不能为空", nil))
        }
        at.BizContent["refresh_token"] = at.refreshToken
    } else if len(at.grantType) == 0 {
        panic(mperr.NewAliPayAuth(errorcode.AliPayAuthParam, "准许类型不能为空", nil))
    }
    at.BizContent["grant_type"] = at.grantType

    return at.GetRequest()
}

func NewAuthToken(appId string) *authToken {
    at := &authToken{alipay.NewBase(appId), "", "", ""}
    at.SetMethod("alipay.open.auth.token.app")
    return at
}
