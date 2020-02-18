/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/10 0010
 * Time: 11:13
 */
package authorize

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/valyala/fasthttp"
)

// 拉取当前所有已授权的帐号基本信息列表
type authorizeList struct {
    wx.BaseWxOpen
    componentAppId string
    offset         int // 偏移位置
    count          int // 拉取数量
}

func (al *authorizeList) SetRange(page, limit int) {
    truePage := 0
    if page > 0 {
        truePage = page
    } else {
        truePage = 1
    }
    if (limit > 0) && (limit <= 500) {
        al.count = limit
    } else {
        al.count = 100
    }
    al.offset = (truePage - 1) * al.count
}

func (al *authorizeList) SendRequest() api.ApiResult {
    reqData := make(map[string]interface{})
    reqData["component_appid"] = al.componentAppId
    reqData["offset"] = al.offset
    reqData["count"] = al.count
    reqBody := mpf.JSONMarshal(reqData)
    al.ReqUrl = "https://api.weixin.qq.com/cgi-bin/component/api_get_authorizer_list?component_access_token=" + wx.NewUtilWx().GetOpenAccessToken()
    client, req := al.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := al.SendInner(client, req, errorcode.WxOpenRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    _, ok := respData["errcode"]
    if ok {
        result.Code = errorcode.WxOpenRequestPost
        result.Msg = respData["errmsg"].(string)
    } else {
        result.Data = respData
    }

    return result
}

func NewAuthorizeList() *authorizeList {
    conf := wx.NewConfig().GetOpen()
    al := &authorizeList{wx.NewBaseWxOpen(), "", 0, 0}
    al.componentAppId = conf.GetAppId()
    al.offset = 0
    al.count = 20
    al.ReqContentType = project.HTTPContentTypeJSON
    al.ReqMethod = fasthttp.MethodPost
    return al
}
