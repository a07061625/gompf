/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/12 0012
 * Time: 18:31
 */
package message

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
)

// 获取模板消息列表
type templateList struct {
    wx.BaseWxAccount
    appId string
}

func (tl *templateList) SendRequest() api.ApiResult {
    tl.ReqUrl = "https://api.weixin.qq.com/cgi-bin/template/get_all_private_template?access_token=" + wx.NewUtilWx().GetSingleAccessToken(tl.appId)
    client, req := tl.GetRequest()

    resp, result := tl.SendInner(client, req, errorcode.WxAccountRequestGet)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    _, ok := respData["template_list"]
    if ok {
        result.Data = respData
    } else {
        result.Code = errorcode.WxAccountRequestGet
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewTemplateList(appId string) *templateList {
    tl := &templateList{wx.NewBaseWxAccount(), ""}
    tl.appId = appId
    return tl
}
