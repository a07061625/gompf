/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/12 0012
 * Time: 16:58
 */
package message

import (
    "regexp"
    "strconv"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
)

// 生成订阅消息授权链接
type subscribeUrl struct {
    wx.BaseWxAccount
    appId        string
    action       string // 动作标识
    scene        int    // 订阅场景值
    templateId   string // 消息模板ID
    redirectUrl  string // 重定向地址
    reservedFlag string // 防止跨站请求伪造攻击标识
}

func (su *subscribeUrl) SetScene(scene int) {
    if (scene >= 0) && (scene <= 10000) {
        su.scene = scene
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "订阅场景值不合法", nil))
    }
}

func (su *subscribeUrl) SetTemplateId(templateId string) {
    if len(templateId) > 0 {
        su.templateId = templateId
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "消息模板ID不合法", nil))
    }
}

func (su *subscribeUrl) SetRedirectUrl(redirectUrl string) {
    match, _ := regexp.MatchString(project.RegexUrlHttp, redirectUrl)
    if match {
        su.redirectUrl = redirectUrl
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "重定向地址不合法", nil))
    }
}

func (su *subscribeUrl) SetReservedFlag(reservedFlag string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{1,128}$`, reservedFlag)
    if match {
        su.ReqData["reserved"] = reservedFlag
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "防止跨站请求伪造攻击标识不合法", nil))
    }
}

func (su *subscribeUrl) checkData() {
    if su.scene < 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "订阅场景值不能为空", nil))
    }
    if len(su.templateId) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "消息模板ID不能为空", nil))
    }
    if len(su.redirectUrl) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "重定向地址不能为空", nil))
    }
    su.ReqData["scene"] = strconv.Itoa(su.scene)
    su.ReqData["template_id"] = su.templateId
    su.ReqData["redirect_url"] = su.redirectUrl
}

func (su *subscribeUrl) GetResult() map[string]string {
    su.checkData()

    result := make(map[string]string)
    result["url"] = "https://mp.weixin.qq.com/mp/subscribemsg?" + mpf.HttpCreateParams(su.ReqData, "none", 1) + "#wechat_redirect"
    return result
}

func NewSubscribeUrl(appId string) *subscribeUrl {
    su := &subscribeUrl{wx.NewBaseWxAccount(), "", "", 0, "", "", ""}
    su.appId = appId
    su.scene = -1
    su.ReqData["appid"] = appId
    su.ReqData["action"] = "get_confirm"
    su.ReqData["reserved"] = mpf.ToolCreateNonceStr(8, "numlower")
    return su
}
