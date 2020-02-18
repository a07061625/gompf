/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/6 0006
 * Time: 17:50
 */
package authorize

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
)

// 获取扫码登录授权引导地址
type loginAuthUrlScan struct {
    wx.BaseWxProvider
    redirectUri string // 授权回调地址
    state       string // 防跨域攻击标识
    userType    string // 登录类型 admin:管理员登录(使用微信扫码) member:成员登录(使用企业微信扫码),默认为admin
}

func (lus *loginAuthUrlScan) SetState(state string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{1,128}$`, state)
    if match {
        lus.ReqData["state"] = state
    } else {
        panic(mperr.NewWxProvider(errorcode.WxProviderParam, "防跨域攻击标识不合法", nil))
    }
}

func (lus *loginAuthUrlScan) SetUserType(userType string) {
    if (userType == "admin") || (userType == "member") {
        lus.ReqData["usertype"] = userType
    } else {
        panic(mperr.NewWxProvider(errorcode.WxProviderParam, "登录类型不合法", nil))
    }
}

func (lus *loginAuthUrlScan) GetResult() map[string]string {
    result := make(map[string]string)
    result["url"] = "https://open.work.weixin.qq.com/wwopen/sso/3rd_qrConnect?" + mpf.HTTPCreateParams(lus.ReqData, "none", 1)
    return result
}

func NewLoginAuthUrlScan() *loginAuthUrlScan {
    conf := wx.NewConfig().GetProvider()
    lus := &loginAuthUrlScan{wx.NewBaseWxProvider(), "", "", ""}
    lus.ReqData["appid"] = conf.GetCorpId()
    lus.ReqData["usertype"] = "admin"
    lus.ReqData["redirect_uri"] = conf.GetUrlAuthLogin()
    lus.ReqData["state"] = mpf.ToolCreateNonceStr(8, "numlower")
    return lus
}
