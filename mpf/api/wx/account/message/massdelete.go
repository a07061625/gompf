/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/12 0012
 * Time: 15:23
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

// 删除群发
type massDelete struct {
    wx.BaseWxAccount
    appId        string
    msgId        string // 消息ID
    articleIndex int    // 消息索引
}

func (md *massDelete) SetMsgId(msgId string) {
    match, _ := regexp.MatchString(project.RegexDigit, msgId)
    if match {
        md.msgId = msgId
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "消息ID不合法", nil))
    }
}

func (md *massDelete) SetArticleIndex(articleIndex int) {
    if articleIndex >= 0 {
        md.articleIndex = articleIndex
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "消息索引不合法", nil))
    }
}

func (md *massDelete) checkData() {
    if len(md.msgId) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "消息ID不能为空", nil))
    }
}

func (md *massDelete) SendRequest() api.ApiResult {
    md.checkData()

    reqData := make(map[string]interface{})
    reqData["msg_id"] = md.msgId
    reqData["article_idx"] = md.articleIndex
    reqBody := mpf.JSONMarshal(reqData)
    md.ReqUrl = "https://api.weixin.qq.com/cgi-bin/message/mass/delete?access_token=" + wx.NewUtilWx().GetSingleAccessToken(md.appId)
    client, req := md.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := md.SendInner(client, req, errorcode.WxAccountRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxAccountRequestPost
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewMassDelete(appId string) *massDelete {
    md := &massDelete{wx.NewBaseWxAccount(), "", "", 0}
    md.appId = appId
    md.articleIndex = 0
    md.ReqContentType = project.HTTPContentTypeJSON
    md.ReqMethod = fasthttp.MethodPost
    return md
}
