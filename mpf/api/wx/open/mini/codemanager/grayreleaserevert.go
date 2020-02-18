/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/11 0011
 * Time: 12:55
 */
package codemanager

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
)

// 取消分阶段发布
type grayReleaseRevert struct {
    wx.BaseWxOpen
    appId string // 应用ID
}

func (grr *grayReleaseRevert) SendRequest() api.APIResult {
    grr.ReqURI = "https://api.weixin.qq.com/wxa/revertgrayrelease?access_token=" + wx.NewUtilWx().GetOpenAuthorizeAccessToken(grr.appId)
    client, req := grr.GetRequest()

    resp, result := grr.SendInner(client, req, errorcode.WxOpenRequestGet)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxOpenRequestGet
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewGrayReleaseRevert(appId string) *grayReleaseRevert {
    grr := &grayReleaseRevert{wx.NewBaseWxOpen(), ""}
    grr.appId = appId
    return grr
}
