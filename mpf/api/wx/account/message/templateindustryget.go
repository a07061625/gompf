/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/12 0012
 * Time: 18:35
 */
package message

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
)

// 获取模板消息行业信息
type templateIndustryGet struct {
    wx.BaseWxAccount
    appId string
}

func (tig *templateIndustryGet) SendRequest() api.ApiResult {
    tig.ReqUrl = "https://api.weixin.qq.com/cgi-bin/template/get_industry?access_token=" + wx.NewUtilWx().GetSingleAccessToken(tig.appId)
    client, req := tig.GetRequest()

    resp, result := tig.SendInner(client, req, errorcode.WxAccountRequestGet)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    _, ok := respData["errcode"]
    if ok {
        result.Code = errorcode.WxAccountRequestGet
        result.Msg = respData["errmsg"].(string)
    } else {
        result.Data = respData
    }

    return result
}

func NewTemplateIndustryGet(appId string) *templateIndustryGet {
    tig := &templateIndustryGet{wx.NewBaseWxAccount(), ""}
    tig.appId = appId
    return tig
}
