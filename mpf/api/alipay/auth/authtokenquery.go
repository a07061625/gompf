/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/19 0019
 * Time: 11:08
 */
package auth

import (
    "regexp"

    "github.com/a07061625/gompf/mpf/api/alipay"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 查询某个应用授权AppAuthToken的授权信息
type authTokenQuery struct {
    alipay.BaseAliPay
    appAuthToken string // 应用授权令牌
}

func (atq *authTokenQuery) SetAppAuthToken(appAuthToken string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{1,128}$`, appAuthToken)
    if match {
        atq.appAuthToken = appAuthToken
    } else {
        panic(mperr.NewAliPayAuth(errorcode.AliPayAuthParam, "应用授权令牌不合法", nil))
    }
}

func (atq *authTokenQuery) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(atq.appAuthToken) == 0 {
        panic(mperr.NewAliPayAuth(errorcode.AliPayAuthParam, "应用授权令牌不能为空", nil))
    }
    atq.BizContent["app_auth_token"] = atq.appAuthToken

    return atq.GetRequest()
}

func NewAuthTokenQuery(appId string) *authTokenQuery {
    atq := &authTokenQuery{alipay.NewBase(appId), ""}
    atq.SetMethod("alipay.open.auth.token.app.query")
    return atq
}
