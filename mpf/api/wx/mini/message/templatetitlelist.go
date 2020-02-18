/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/6 0006
 * Time: 9:20
 */
package message

import (
    "strconv"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/valyala/fasthttp"
)

type templateTitleList struct {
    wx.BaseWxMini
    appId  string // 应用ID
    offset int    // 位移
    count  int    // 记录数
}

func (ttl *templateTitleList) SetRange(page int, limit int) {
    truePage := 0
    if page > 0 {
        truePage = page
    } else {
        truePage = 1
    }
    if (limit > 0) && (limit <= 20) {
        ttl.count = limit
    } else {
        ttl.count = 20
    }
    ttl.offset = (truePage - 1) * ttl.count
    ttl.ReqData["count"] = strconv.Itoa(ttl.count)
    ttl.ReqData["offset"] = strconv.Itoa(ttl.offset)
}

func (ttl *templateTitleList) SendRequest(getType string) api.ApiResult {
    reqBody := mpf.JsonMarshal(ttl.ReqData)

    ttl.ReqUrl = "https://api.weixin.qq.com/cgi-bin/wxopen/template/library/list?access_token=" + wx.NewUtilWx().GetSingleCache(ttl.appId, getType)
    client, req := ttl.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := ttl.SendInner(client, req, errorcode.WxMiniRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    _, ok := respData["list"]
    if ok {
        result.Data = respData
    } else {
        result.Code = errorcode.WxMiniRequestPost
        result.Msg = respData["errmsg"].(string)
    }
    return result
}

func NewTemplateTitleList(appId string) *templateTitleList {
    ttl := &templateTitleList{wx.NewBaseWxMini(), "", 0, 0}
    ttl.appId = appId
    ttl.ReqData["offset"] = "0"
    ttl.ReqData["count"] = "20"
    ttl.ReqContentType = project.HTTPContentTypeJSON
    ttl.ReqMethod = fasthttp.MethodPost
    return ttl
}
