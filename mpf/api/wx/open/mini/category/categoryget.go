/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/10 0010
 * Time: 23:34
 */
package category

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
)

// 获取已设置的二级类目及用于代码审核的可选三级类目
type categoryGet struct {
    wx.BaseWxOpen
    appId string // 应用ID
}

func (cg *categoryGet) SendRequest() api.APIResult {
    cg.ReqURI = "https://api.weixin.qq.com/wxa/get_category?access_token=" + wx.NewUtilWx().GetOpenAuthorizeAccessToken(cg.appId)
    client, req := cg.GetRequest()

    resp, result := cg.SendInner(client, req, errorcode.WxOpenRequestGet)
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

func NewCategoryGet(appId string) *categoryGet {
    cg := &categoryGet{wx.NewBaseWxOpen(), ""}
    cg.appId = appId
    return cg
}
