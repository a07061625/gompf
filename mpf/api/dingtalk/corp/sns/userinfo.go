/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/1 0001
 * Time: 0:30
 */
package sns

import (
    "regexp"

    "github.com/a07061625/gompf/mpf/api/dingtalk"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 获取授权用户的个人信息
type userInfo struct {
    dingtalk.BaseCorp
    corpId         string
    openid         string // 用户openid
    persistentCode string // 持久授权码
}

func (ui *userInfo) SetOpenid(openid string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, openid)
    if match {
        ui.openid = openid
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "用户openid不合法", nil))
    }
}

func (ui *userInfo) SetPersistentCode(persistentCode string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, persistentCode)
    if match {
        ui.persistentCode = persistentCode
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "持久授权码不合法", nil))
    }
}

func (ui *userInfo) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ui.openid) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "用户openid不能为空", nil))
    }
    if len(ui.persistentCode) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "持久授权码不能为空", nil))
    }

    ui.ReqUrl = dingtalk.UrlService + "/sns/getuserinfo?sns_token="
    if len(ui.corpId) > 0 {
        ui.ReqUrl += dingtalk.NewUtil().GetCorpUserSnsToken(ui.corpId, ui.openid, ui.persistentCode)
    } else {
        ui.ReqUrl += dingtalk.NewUtil().GetProviderUserSnsToken(ui.openid, ui.persistentCode)
    }

    return ui.GetRequest()
}

func NewUserInfo(corpId string) *userInfo {
    ui := &userInfo{dingtalk.NewCorp(), "", "", ""}
    ui.corpId = corpId
    return ui
}
