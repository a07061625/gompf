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
    wx.BaseWxCorp
    corpId      string
    agentTag    string
    redirectUri string // 授权回调地址
    state       string // 防跨域攻击标识
}

func (lus *loginAuthUrlScan) SetState(state string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{1,128}$`, state)
    if match {
        lus.ReqData["state"] = state
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "防跨域攻击标识不合法", nil))
    }
}

func (lus *loginAuthUrlScan) GetResult() map[string]string {
    result := make(map[string]string)
    result["url"] = "https://open.work.weixin.qq.com/wwopen/sso/qrConnect?" + mpf.HTTPCreateParams(lus.ReqData, "none", 1)
    return result
}

func NewLoginAuthUrlScan(corpId, agentTag string) *loginAuthUrlScan {
    conf := wx.NewConfig().GetCorp(corpId)
    agentInfo := conf.GetAgentInfo(agentTag)
    lus := &loginAuthUrlScan{wx.NewBaseWxCorp(), "", "", "", ""}
    lus.corpId = corpId
    lus.agentTag = agentTag
    lus.ReqData["appid"] = corpId
    lus.ReqData["agentid"] = agentInfo["id"]
    lus.ReqData["redirect_uri"] = conf.GetUrlAuthLogin()
    lus.ReqData["state"] = mpf.ToolCreateNonceStr(8, "numlower")
    return lus
}
