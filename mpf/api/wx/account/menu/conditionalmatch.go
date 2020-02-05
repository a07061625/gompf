/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/12 0012
 * Time: 13:10
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

// 测试个性化菜单匹配结果
type conditionalMatch struct {
    wx.BaseWxAccount
    appId  string
    userId string // 用户openid或粉丝的微信号
}

func (cm *conditionalMatch) SetUserId(userId string) {
    if len(userId) > 0 {
        cm.userId = userId
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "用户ID不合法", nil))
    }
}

func (cm *conditionalMatch) checkData() {
    if len(cm.userId) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "用户ID不能为空", nil))
    }
    cm.ReqData["user_id"] = cm.userId
}

func (cm *conditionalMatch) SendRequest() api.ApiResult {
    cm.checkData()

    reqBody := mpf.JsonMarshal(cm.ReqData)
    cm.ReqUrl = "https://api.weixin.qq.com/cgi-bin/menu/trymatch?access_token=" + wx.NewUtilWx().GetSingleAccessToken(cm.appId)
    client, req := cm.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := cm.SendInner(client, req, errorcode.WxAccountRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    _, ok := respData["button"]
    if ok {
        result.Data = respData
    } else {
        result.Code = errorcode.WxAccountRequestPost
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewConditionalMatch(appId string) *conditionalMatch {
    cm := &conditionalMatch{wx.NewBaseWxAccount(), "", ""}
    cm.appId = appId
    cm.ReqContentType = project.HttpContentTypeJson
    cm.ReqMethod = fasthttp.MethodPost
    return cm
}
