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

// 增量更新成员
type userSync struct {
    wx.BaseWxCorp
    corpId     string
    agentTag   string
    mediaId    string            // 媒体ID
    inviteFlag bool              // 邀请标识,默认值为true true:邀请使用企业微信 false:不邀请
    callback   map[string]string // 回调信息
}

func (us *userSync) SetMediaId(mediaId string) {
    if len(mediaId) > 0 {
        us.mediaId = mediaId
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "媒体ID不合法", nil))
    }
}

func (us *userSync) SetInviteFlag(inviteFlag bool) {
    us.inviteFlag = inviteFlag
}

func (us *userSync) SetCallback(callback map[string]string) {
    if len(callback) > 0 {
        us.callback = callback
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "回调信息不合法", nil))
    }
}

func (us *userSync) checkData() {
    if len(us.mediaId) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "媒体ID不能为空", nil))
    }
}

func (us *userSync) SendRequest(getType string) api.ApiResult {
    us.checkData()
    reqData := make(map[string]interface{})
    reqData["media_id"] = us.mediaId
    reqData["to_invite"] = us.inviteFlag
    if len(us.callback) > 0 {
        reqData["callback"] = us.callback
    }
    reqBody := mpf.JsonMarshal(reqData)

    us.ReqUrl = "https://qyapi.weixin.qq.com/cgi-bin/batch/syncuser?access_token=" + wx.NewUtilWx().GetCorpCache(us.corpId, us.agentTag, getType)
    client, req := us.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := us.SendInner(client, req, errorcode.WxCorpRequestPost)
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

func NewUserSync(corpId, agentTag string) *userSync {
    us := &userSync{wx.NewBaseWxCorp(), "", "", "", false, make(map[string]string)}
    us.corpId = corpId
    us.agentTag = agentTag
    us.inviteFlag = true
    us.ReqContentType = project.HttpContentTypeJson
    us.ReqMethod = fasthttp.MethodPost
    return us
}
