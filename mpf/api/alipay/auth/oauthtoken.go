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

// 换取授权访问令牌
type oauthToken struct {
    alipay.BaseAliPay
    grantType    string // 准许类型
    authCode     string // 授权码
    refreshToken string // 刷新令牌
}

func (ot *oauthToken) SetGrantType(grantType string) {
    if (grantType == "authorization_code") || (grantType == "refresh_token") {
        ot.grantType = grantType
    } else {
        panic(mperr.NewAliPayAuth(errorcode.AliPayAuthParam, "准许类型不合法", nil))
    }
}

func (ot *oauthToken) SetAuthCode(authCode string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{1,128}$`, authCode)
    if match {
        ot.authCode = authCode
    } else {
        panic(mperr.NewAliPayAuth(errorcode.AliPayAuthParam, "授权码不合法", nil))
    }
}

func (ot *oauthToken) SetRefreshToken(refreshToken string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{1,40}$`, refreshToken)
    if match {
        ot.refreshToken = refreshToken
    } else {
        panic(mperr.NewAliPayAuth(errorcode.AliPayAuthParam, "刷新令牌不合法", nil))
    }
}

func (ot *oauthToken) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if ot.grantType == "authorization_code" {
        if len(ot.authCode) == 0 {
            panic(mperr.NewAliPayAuth(errorcode.AliPayAuthParam, "授权码不能为空", nil))
        }
        ot.BizContent["code"] = ot.authCode
    } else if ot.grantType == "refresh_token" {
        if len(ot.refreshToken) == 0 {
            panic(mperr.NewAliPayAuth(errorcode.AliPayAuthParam, "刷新令牌不能为空", nil))
        }
        ot.BizContent["refresh_token"] = ot.refreshToken
    } else if len(ot.grantType) == 0 {
        panic(mperr.NewAliPayAuth(errorcode.AliPayAuthParam, "准许类型不能为空", nil))
    }
    ot.BizContent["grant_type"] = ot.grantType

    return ot.GetRequest()
}

func NewOauthToken(appId string) *oauthToken {
    ot := &oauthToken{alipay.NewBase(appId), "", "", ""}
    ot.SetMethod("alipay.system.oauth.token")
    return ot
}
