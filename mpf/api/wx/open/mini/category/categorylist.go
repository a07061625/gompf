/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/10 0010
 * Time: 23:37
 */
package category

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
)

// 获取账号已经设置的所有类目
type categoryList struct {
    wx.BaseWxOpen
    appId string // 应用ID
}

func (cl *categoryList) SendRequest() api.ApiResult {
    cl.ReqUrl = "https://api.weixin.qq.com/cgi-bin/wxopen/getcategory?access_token=" + wx.NewUtilWx().GetOpenAuthorizeAccessToken(cl.appId)
    client, req := cl.GetRequest()

    resp, result := cl.SendInner(client, req, errorcode.WxOpenRequestGet)
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

func NewCategoryList(appId string) *categoryList {
    cl := &categoryList{wx.NewBaseWxOpen(), ""}
    cl.appId = appId
    return cl
}
