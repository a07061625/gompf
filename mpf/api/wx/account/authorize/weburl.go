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

// 网站应用授权地址
type webUrl struct {
    wx.BaseWxAccount
    appId       string
    redirectUrl string // 重定向链接
    state       string // 防csrf攻击标识
}

func (wu *webUrl) SetRedirectUrl(redirectUrl string) {
    match, _ := regexp.MatchString(project.RegexURLHTTP, redirectUrl)
    if match {
        wu.redirectUrl = redirectUrl
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "重定向链接不合法", nil))
    }
}

func (wu *webUrl) SetState(state string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{1,32}$`, state)
    if match {
        wu.ReqData["state"] = state
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "防csrf攻击标识不合法", nil))
    }
}

func (wu *webUrl) checkData() {
    if len(wu.redirectUrl) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "重定向链接不能为空", nil))
    }
    wu.ReqData["redirect_uri"] = wu.redirectUrl
}

func (wu *webUrl) GetResult() map[string]string {
    wu.checkData()

    result := make(map[string]string)
    result["url"] = "https://open.weixin.qq.com/connect/qrconnect?" + mpf.HTTPCreateParams(wu.ReqData, "key", 1) + "#wechat_redirect"
    return result
}

func NewWebUrl(appId string) *webUrl {
    wu := &webUrl{wx.NewBaseWxAccount(), "", "", ""}
    wu.appId = appId
    wu.ReqData["appid"] = appId
    wu.ReqData["scope"] = "snsapi_login"
    wu.ReqData["state"] = "STATE"
    wu.ReqData["response_type"] = "code"
    return wu
}
