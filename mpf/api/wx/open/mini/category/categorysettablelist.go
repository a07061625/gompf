/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/10 0010
 * Time: 23:47
 */
package category

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
)

// 可设置的所有类目
type categorySettableList struct {
    wx.BaseWxOpen
    appId string // 应用ID
}

func (csl *categorySettableList) SendRequest() api.ApiResult {
    csl.ReqUrl = "https://api.weixin.qq.com/cgi-bin/wxopen/getallcategories?access_token=" + wx.NewUtilWx().GetOpenAuthorizeAccessToken(csl.appId)
    client, req := csl.GetRequest()

    resp, result := csl.SendInner(client, req, errorcode.WxOpenRequestGet)
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

func NewCategorySettableList(appId string) *categorySettableList {
    csl := &categorySettableList{wx.NewBaseWxOpen(), ""}
    csl.appId = appId
    return csl
}
