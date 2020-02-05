/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/9 0009
 * Time: 19:31
 */
package authorize

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 获取访问用户敏感信息
type userDetailGet struct {
    wx.BaseWxProvider
    userTicket string // 成员票据
}

func (udg *userDetailGet) SetUserTicket(userTicket string) {
    if len(userTicket) > 0 {
        udg.userTicket = userTicket
    } else {
        panic(mperr.NewWxProvider(errorcode.WxProviderParam, "成员票据不合法", nil))
    }
}

func (udg *userDetailGet) checkData() {
    if len(udg.userTicket) == 0 {
        panic(mperr.NewWxProvider(errorcode.WxProviderParam, "成员票据不能为空", nil))
    }
    udg.ReqData["user_ticket"] = udg.userTicket
}

func (udg *userDetailGet) SendRequest() api.ApiResult {
    udg.checkData()

    reqBody := mpf.JsonMarshal(udg.ReqData)
    udg.ReqUrl = "https://qyapi.weixin.qq.com/cgi-bin/service/getuserdetail3rd?access_token=" + wx.NewUtilWx().GetProviderSuiteToken()
    client, req := udg.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := udg.SendInner(client, req, errorcode.WxProviderRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxProviderRequestPost
        result.Msg = respData["errmsg"].(string)
    }
    return result
}

func NewUserDetailGet() *userDetailGet {
    udg := &userDetailGet{wx.NewBaseWxProvider(), ""}
    udg.ReqContentType = project.HttpContentTypeJson
    udg.ReqMethod = fasthttp.MethodPost
    return udg
}
