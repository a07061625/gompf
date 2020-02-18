/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/12 0012
 * Time: 23:20
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

type quotaClear struct {
    wx.BaseWxAccount
    appId string
}

func (qc *quotaClear) SendRequest(getType string) api.ApiResult {
    reqBody := mpf.JsonMarshal(qc.ReqData)
    qc.ReqUrl = "https://api.weixin.qq.com/cgi-bin/clear_quota?access_token=" + wx.NewUtilWx().GetSingleCache(qc.appId, getType)
    client, req := qc.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := qc.SendInner(client, req, errorcode.WxAccountRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxAccountRequestPost
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewQuotaClear(appId string) *quotaClear {
    qc := &quotaClear{wx.NewBaseWxAccount(), ""}
    qc.appId = appId
    qc.ReqData["appid"] = appId
    qc.ReqContentType = project.HTTPContentTypeJSON
    qc.ReqMethod = fasthttp.MethodPost
    return qc
}
