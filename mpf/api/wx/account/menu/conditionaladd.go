/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/12 0012
 * Time: 12:24
 */
package menu

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 创建个性化菜单
type conditionalAdd struct {
    wx.BaseWxAccount
    appId      string
    buttonList []map[string]interface{} // 一级菜单列表
    matchRule  map[string]interface{}   // 菜单匹配规则
}

func (ca *conditionalAdd) AddMenu(menu *menu) {
    if len(ca.buttonList) >= 3 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "菜单数量不能超过3个", nil))
    }
    ca.buttonList = append(ca.buttonList, menu.GetResult())
}

func (ca *conditionalAdd) SetMatchRule(matchRule map[string]interface{}) {
    if len(matchRule) == 0 {
        ca.matchRule = matchRule
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "菜单匹配规则不合法", nil))
    }
}

func (ca *conditionalAdd) checkData() {
    if len(ca.buttonList) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "菜单不能为空", nil))
    }
    if len(ca.matchRule) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "菜单匹配规则不能为空", nil))
    }
}

func (ca *conditionalAdd) SendRequest() api.ApiResult {
    ca.checkData()

    reqData := make(map[string]interface{})
    reqData["button"] = ca.buttonList
    reqData["matchrule"] = ca.matchRule
    reqBody := mpf.JsonMarshal(reqData)
    ca.ReqUrl = "https://api.weixin.qq.com/cgi-bin/menu/addconditional?access_token=" + wx.NewUtilWx().GetSingleAccessToken(ca.appId)
    client, req := ca.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := ca.SendInner(client, req, errorcode.WxAccountRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    _, ok := respData["menuid"]
    if ok {
        result.Data = respData
    } else {
        result.Code = errorcode.WxAccountRequestPost
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewConditionalAdd(appId string) *conditionalAdd {
    ca := &conditionalAdd{wx.NewBaseWxAccount(), "", make([]map[string]interface{}, 0), make(map[string]interface{})}
    ca.appId = appId
    ca.ReqContentType = project.HttpContentTypeJson
    ca.ReqMethod = fasthttp.MethodPost
    return ca
}
