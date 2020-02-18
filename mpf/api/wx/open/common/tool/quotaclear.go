/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/10 0010
 * Time: 9:00
 */
package tool

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/valyala/fasthttp"
)

// 调用接口次数清零
type quotaClear struct {
    wx.BaseWxOpen
}

func (qc *quotaClear) SendRequest() api.APIResult {
    reqBody := mpf.JSONMarshal(qc.ReqData)
    qc.ReqURI = "https://api.weixin.qq.com/cgi-bin/component/clear_quota?component_access_token=" + wx.NewUtilWx().GetOpenAccessToken()
    client, req := qc.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := qc.SendInner(client, req, errorcode.WxOpenRequestPost)
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

func NewQuotaClear(appId string) *quotaClear {
    conf := wx.NewConfig().GetOpen()
    qc := &quotaClear{wx.NewBaseWxOpen()}
    qc.ReqData["component_appid"] = conf.GetAppId()
    qc.ReqContentType = project.HTTPContentTypeJSON
    qc.ReqMethod = fasthttp.MethodPost
    return qc
}
