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

// 获取网页登录授权引导地址
type loginAuthUrlWeb struct {
    wx.BaseWxCorp
    corpId       string
    redirectUri  string // 授权回调地址
    responseType string // 返回类型
    scope        string // 授权作用域
    state        string // 防跨域攻击标识
}

func (luw *loginAuthUrlWeb) SetState(state string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{1,128}$`, state)
    if match {
        luw.ReqData["state"] = state
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "防跨域攻击标识不合法", nil))
    }
}

func (luw *loginAuthUrlWeb) GetResult() map[string]string {
    result := make(map[string]string)
    result["url"] = "https://open.weixin.qq.com/connect/oauth2/authorize?" + mpf.HttpCreateParams(luw.ReqData, "none", 1) + "#wechat_redirect"
    return result
}

func NewLoginAuthUrlWeb(corpId string) *loginAuthUrlWeb {
    conf := wx.NewConfig().GetCorp(corpId)
    luw := &loginAuthUrlWeb{wx.NewBaseWxCorp(), "", "", "", "", ""}
    luw.corpId = corpId
    luw.ReqData["appid"] = conf.GetCorpId()
    luw.ReqData["redirect_uri"] = conf.GetUrlAuthLogin()
    luw.ReqData["response_type"] = "code"
    luw.ReqData["scope"] = "snsapi_base"
    luw.ReqData["state"] = mpf.ToolCreateNonceStr(8, "numlower")
    return luw
}
