/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/12 0012
 * Time: 9:56
 */
package material

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 删除永久素材
type materialDel struct {
    wx.BaseWxAccount
    appId   string
    mediaId string // 媒体文件ID
}

func (md *materialDel) SetMediaId(mediaId string) {
    if len(mediaId) > 0 {
        md.mediaId = mediaId
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "媒体文件ID不合法", nil))
    }
}

func (md *materialDel) checkData() {
    if len(md.mediaId) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "媒体文件ID不能为空", nil))
    }
    md.ReqData["media_id"] = md.mediaId
}

func (md *materialDel) SendRequest() api.ApiResult {
    md.checkData()

    reqBody := mpf.JsonMarshal(md.ReqData)
    md.ReqUrl = "https://api.weixin.qq.com/cgi-bin/material/del_material?access_token=" + wx.NewUtilWx().GetSingleAccessToken(md.appId)
    client, req := md.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := md.SendInner(client, req, errorcode.WxAccountRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxAccountRequestPost
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewMaterialDel(appId string) *materialDel {
    md := &materialDel{wx.NewBaseWxAccount(), "", ""}
    md.appId = appId
    md.ReqContentType = project.HttpContentTypeJson
    md.ReqMethod = fasthttp.MethodPost
    return md
}
