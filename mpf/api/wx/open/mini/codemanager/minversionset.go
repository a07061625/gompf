/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/11 0011
 * Time: 13:12
 */
package codemanager

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 设置最低基础库版本
type minVersionSet struct {
    wx.BaseWxOpen
    appId   string // 应用ID
    version string // 最低版本号
}

func (mvs *minVersionSet) SetVersion(version string) {
    if len(version) > 0 {
        mvs.version = version
    } else {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "最低版本号不合法", nil))
    }
}

func (mvs *minVersionSet) checkData() {
    if len(mvs.version) == 0 {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "最低版本号不能为空", nil))
    }
    mvs.ReqData["version"] = mvs.version
}

func (mvs *minVersionSet) SendRequest() api.ApiResult {
    mvs.checkData()

    reqBody := mpf.JSONMarshal(mvs.ReqData)
    mvs.ReqUrl = "https://api.weixin.qq.com/cgi-bin/wxopen/setweappsupportversion?access_token=" + wx.NewUtilWx().GetOpenAuthorizeAccessToken(mvs.appId)
    client, req := mvs.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := mvs.SendInner(client, req, errorcode.WxOpenRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxOpenRequestPost
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewMinVersionSet(appId string) *minVersionSet {
    mvs := &minVersionSet{wx.NewBaseWxOpen(), "", ""}
    mvs.appId = appId
    mvs.ReqContentType = project.HTTPContentTypeJSON
    mvs.ReqMethod = fasthttp.MethodPost
    return mvs
}
