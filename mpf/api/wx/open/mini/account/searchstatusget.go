/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/10 0010
 * Time: 22:39
 */
package account

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
)

// 查询小程序当前隐私设置（是否可被搜索）
type searchStatusGet struct {
    wx.BaseWxOpen
    appId string // 应用ID
}

func (ssg *searchStatusGet) SendRequest() api.ApiResult {
    ssg.ReqUrl = "https://api.weixin.qq.com/wxa/getwxasearchstatus?access_token=" + wx.NewUtilWx().GetOpenAuthorizeAccessToken(ssg.appId)
    client, req := ssg.GetRequest()

    resp, result := ssg.SendInner(client, req, errorcode.WxOpenRequestGet)
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

func NewSearchStatusGet(appId string) *searchStatusGet {
    ssg := &searchStatusGet{wx.NewBaseWxOpen(), ""}
    ssg.appId = appId
    return ssg
}
