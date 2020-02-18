/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/6 0006
 * Time: 23:17
 */
package batch

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 全量覆盖部门
type partyReplace struct {
    wx.BaseWxCorp
    corpId   string
    agentTag string
    mediaId  string            // 媒体ID
    callback map[string]string // 回调信息
}

func (pr *partyReplace) SetMediaId(mediaId string) {
    if len(mediaId) > 0 {
        pr.mediaId = mediaId
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "媒体ID不合法", nil))
    }
}

func (pr *partyReplace) SetCallback(callback map[string]string) {
    if len(callback) > 0 {
        pr.callback = callback
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "回调信息不合法", nil))
    }
}

func (pr *partyReplace) checkData() {
    if len(pr.mediaId) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "媒体ID不能为空", nil))
    }
}

func (pr *partyReplace) SendRequest() api.ApiResult {
    pr.checkData()
    reqData := make(map[string]interface{})
    reqData["media_id"] = pr.mediaId
    if len(pr.callback) > 0 {
        reqData["callback"] = pr.callback
    }
    reqBody := mpf.JsonMarshal(reqData)

    pr.ReqUrl = "https://qyapi.weixin.qq.com/cgi-bin/batch/replaceparty?access_token=" + wx.NewUtilWx().GetCorpAccessToken(pr.corpId, pr.agentTag)
    client, req := pr.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := pr.SendInner(client, req, errorcode.WxCorpRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxCorpRequestPost
        result.Msg = respData["errmsg"].(string)
    }
    return result
}

func NewPartyReplace(corpId, agentTag string) *partyReplace {
    pr := &partyReplace{wx.NewBaseWxCorp(), "", "", "", make(map[string]string)}
    pr.corpId = corpId
    pr.agentTag = agentTag
    pr.ReqContentType = project.HTTPContentTypeJSON
    pr.ReqMethod = fasthttp.MethodPost
    return pr
}
