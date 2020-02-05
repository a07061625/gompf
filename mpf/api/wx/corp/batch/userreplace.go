/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/7 0007
 * Time: 8:55
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

// 全量覆盖成员
type userReplace struct {
    wx.BaseWxCorp
    corpId     string
    agentTag   string
    mediaId    string            // 媒体ID
    inviteFlag bool              // 邀请标识,默认值为true true:邀请使用企业微信 false:不邀请
    callback   map[string]string // 回调信息
}

func (ur *userReplace) SetMediaId(mediaId string) {
    if len(mediaId) > 0 {
        ur.mediaId = mediaId
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "媒体ID不合法", nil))
    }
}

func (ur *userReplace) SetInviteFlag(inviteFlag bool) {
    ur.inviteFlag = inviteFlag
}

func (ur *userReplace) SetCallback(callback map[string]string) {
    if len(callback) > 0 {
        ur.callback = callback
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "回调信息不合法", nil))
    }
}

func (ur *userReplace) checkData() {
    if len(ur.mediaId) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "媒体ID不能为空", nil))
    }
}

func (ur *userReplace) SendRequest() api.ApiResult {
    ur.checkData()
    reqData := make(map[string]interface{})
    reqData["media_id"] = ur.mediaId
    reqData["to_invite"] = ur.inviteFlag
    if len(ur.callback) > 0 {
        reqData["callback"] = ur.callback
    }
    reqBody := mpf.JsonMarshal(reqData)

    ur.ReqUrl = "https://qyapi.weixin.qq.com/cgi-bin/batch/replaceuser?access_token=" + wx.NewUtilWx().GetCorpAccessToken(ur.corpId, ur.agentTag)
    client, req := ur.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := ur.SendInner(client, req, errorcode.WxCorpRequestPost)
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

func NewUserReplace(corpId, agentTag string) *userReplace {
    ur := &userReplace{wx.NewBaseWxCorp(), "", "", "", false, make(map[string]string)}
    ur.corpId = corpId
    ur.agentTag = agentTag
    ur.inviteFlag = true
    ur.ReqContentType = project.HttpContentTypeJson
    ur.ReqMethod = fasthttp.MethodPost
    return ur
}
