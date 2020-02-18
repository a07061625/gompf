/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/12 0012
 * Time: 15:31
 */
package message

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 查询群发消息发送状态
type massGet struct {
    wx.BaseWxAccount
    appId string
    msgId string // 消息ID
}

func (mg *massGet) SetMsgId(msgId string) {
    match, _ := regexp.MatchString(project.RegexDigit, msgId)
    if match {
        mg.msgId = msgId
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "消息ID不合法", nil))
    }
}

func (mg *massGet) checkData() {
    if len(mg.msgId) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "消息ID不能为空", nil))
    }
    mg.ReqData["msg_id"] = mg.msgId
}

func (mg *massGet) SendRequest() api.ApiResult {
    mg.checkData()

    reqBody := mpf.JsonMarshal(mg.ReqData)
    mg.ReqUrl = "https://api.weixin.qq.com/cgi-bin/message/mass/get?access_token=" + wx.NewUtilWx().GetSingleAccessToken(mg.appId)
    client, req := mg.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := mg.SendInner(client, req, errorcode.WxAccountRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    _, ok := respData["msg_id"]
    if ok {
        result.Data = respData
    } else {
        result.Code = errorcode.WxAccountRequestPost
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewMassGet(appId string) *massGet {
    mg := &massGet{wx.NewBaseWxAccount(), "", ""}
    mg.appId = appId
    mg.ReqContentType = project.HTTPContentTypeJSON
    mg.ReqMethod = fasthttp.MethodPost
    return mg
}
