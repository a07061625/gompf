/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/12 0012
 * Time: 11:22
 */
package media

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type videoUpload struct {
    wx.BaseWxAccount
    appId       string
    mediaId     string // 媒体ID
    title       string // 视频标题
    description string // 视频描述
}

func (vu *videoUpload) SetMediaId(mediaId string) {
    if len(mediaId) > 0 {
        vu.mediaId = mediaId
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "媒体ID不合法", nil))
    }
}

func (vu *videoUpload) SetTitle(title string) {
    if len(title) > 0 {
        vu.title = title
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "视频标题不合法", nil))
    }
}

func (vu *videoUpload) SetDescription(description string) {
    vu.ReqData["description"] = description
}

func (vu *videoUpload) checkData() {
    if len(vu.mediaId) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "媒体ID不能为空", nil))
    }
    if len(vu.title) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "视频标题不能为空", nil))
    }
    vu.ReqData["media_id"] = vu.mediaId
    vu.ReqData["title"] = vu.title
}

func (vu *videoUpload) SendRequest() api.ApiResult {
    vu.checkData()

    reqBody := mpf.JSONMarshal(vu.ReqData)
    vu.ReqUrl = "https://api.weixin.qq.com/cgi-bin/media/uploadvideo?access_token=" + wx.NewUtilWx().GetSingleAccessToken(vu.appId)
    client, req := vu.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := vu.SendInner(client, req, errorcode.WxAccountRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    _, ok := respData["media_id"]
    if ok {
        result.Data = respData
    } else {
        result.Code = errorcode.WxAccountRequestPost
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewVideoUpload(appId string) *videoUpload {
    vu := &videoUpload{wx.NewBaseWxAccount(), "", "", "", ""}
    vu.appId = appId
    vu.ReqData["description"] = ""
    vu.ReqContentType = project.HTTPContentTypeJSON
    vu.ReqMethod = fasthttp.MethodPost
    return vu
}
