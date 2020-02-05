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

type templateList struct {
    wx.BaseWxMini
    appId  string // 应用ID
    offset int    // 位移
    count  int    // 记录数
}

func (tl *templateList) SetRange(page int, limit int) {
    truePage := 0
    if page > 0 {
        truePage = page
    } else {
        truePage = 1
    }
    if (limit > 0) && (limit <= 20) {
        tl.count = limit
    } else {
        tl.count = 20
    }
    tl.offset = (truePage - 1) * tl.count
    tl.ReqData["count"] = strconv.Itoa(tl.count)
    tl.ReqData["offset"] = strconv.Itoa(tl.offset)
}

func (tl *templateList) SendRequest(getType string) api.ApiResult {
    reqBody := mpf.JsonMarshal(tl.ReqData)

    tl.ReqUrl = "https://api.weixin.qq.com/cgi-bin/wxopen/template/list?access_token=" + wx.NewUtilWx().GetSingleCache(tl.appId, getType)
    client, req := tl.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := tl.SendInner(client, req, errorcode.WxMiniRequestPost)
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

func NewTemplateList(appId string) *templateList {
    tl := &templateList{wx.NewBaseWxMini(), "", 0, 0}
    tl.appId = appId
    tl.ReqData["offset"] = "0"
    tl.ReqData["count"] = "20"
    tl.ReqContentType = project.HttpContentTypeJson
    tl.ReqMethod = fasthttp.MethodPost
    return tl
}
