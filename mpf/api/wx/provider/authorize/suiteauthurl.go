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

// 获取套件授权引导地址
type suiteAuthUrl struct {
    wx.BaseWxProvider
    preAuthCode string // 预授权码
    redirectUri string // 授权回调地址
    state       string // 防跨域攻击标识
}

func (sau *suiteAuthUrl) SetPreAuthCode(preAuthCode string) {
    if len(preAuthCode) > 0 {
        sau.preAuthCode = preAuthCode
    } else {
        panic(mperr.NewWxProvider(errorcode.WxProviderParam, "预授权码不合法", nil))
    }
}

func (sau *suiteAuthUrl) SetState(state string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{1,128}$`, state)
    if match {
        sau.ReqData["state"] = state
    } else {
        panic(mperr.NewWxProvider(errorcode.WxProviderParam, "防跨域攻击标识不合法", nil))
    }
}

func (sau *suiteAuthUrl) checkData() {
    if len(sau.preAuthCode) == 0 {
        panic(mperr.NewWxProvider(errorcode.WxProviderParam, "预授权码不能为空", nil))
    }
    sau.ReqData["pre_auth_code"] = sau.preAuthCode
}

func (sau *suiteAuthUrl) GetResult() map[string]string {
    sau.checkData()

    result := make(map[string]string)
    result["url"] = "https://open.work.weixin.qq.com/3rdapp/install?" + mpf.HttpCreateParams(sau.ReqData, "none", 1)
    return result
}

func NewSuiteAuthUrl() *suiteAuthUrl {
    conf := wx.NewConfig().GetProvider()
    sau := &suiteAuthUrl{wx.NewBaseWxProvider(), "", "", ""}
    sau.ReqData["suite_id"] = conf.GetSuiteId()
    sau.ReqData["redirect_uri"] = conf.GetUrlAuthSuite()
    sau.ReqData["state"] = mpf.ToolCreateNonceStr(8, "numlower")
    return sau
}
