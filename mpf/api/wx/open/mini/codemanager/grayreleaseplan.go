/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/11 0011
 * Time: 12:15
 */
package codemanager

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
)

// 查询当前分阶段发布详情
type grayReleasePlan struct {
    wx.BaseWxOpen
    appId string // 应用ID
}

func (grp *grayReleasePlan) SendRequest() api.ApiResult {
    grp.ReqUrl = "https://api.weixin.qq.com/wxa/getgrayreleaseplan?access_token=" + wx.NewUtilWx().GetOpenAuthorizeAccessToken(grp.appId)
    client, req := grp.GetRequest()

    resp, result := grp.SendInner(client, req, errorcode.WxOpenRequestGet)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxOpenRequestGet
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewGrayReleasePlan(appId string) *grayReleasePlan {
    grp := &grayReleasePlan{wx.NewBaseWxOpen(), ""}
    grp.appId = appId
    return grp
}
