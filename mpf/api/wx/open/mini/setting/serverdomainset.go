/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/10 0010
 * Time: 14:30
 */
package setting

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 设置小程序服务器域名
type serverDomainSet struct {
    wx.BaseWxOpen
    appId  string // 应用ID
    action string
    data   map[string][]string
}

func (sds *serverDomainSet) SetData(action string, data map[string][]string) {
    switch action {
    case "add":
        if len(data) == 0 {
            panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "域名不能为空", nil))
        }
        sds.data = data
    case "delete":
        if len(data) == 0 {
            panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "域名不能为空", nil))
        }
        sds.data = data
    case "set":
        if len(data) == 0 {
            panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "域名不能为空", nil))
        }
        sds.data = data
    case "get":
    default:
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "操作类型不支持", nil))
    }
    sds.action = action
}

func (sds *serverDomainSet) checkData() {
    if len(sds.action) == 0 {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "操作类型不能为空", nil))
    }
}

func (sds *serverDomainSet) SendRequest() api.ApiResult {
    sds.checkData()

    reqData := make(map[string]interface{})
    if len(sds.data) > 0 {
        for k, v := range sds.data {
            reqData[k] = v
        }
    }
    reqData["action"] = sds.action
    reqBody := mpf.JsonMarshal(reqData)
    sds.ReqUrl = "https://api.weixin.qq.com/wxa/modify_domain?access_token=" + wx.NewUtilWx().GetOpenAuthorizeAccessToken(sds.appId)
    client, req := sds.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := sds.SendInner(client, req, errorcode.WxOpenRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxOpenRequestPost
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewServerDomainSet(appId string) *serverDomainSet {
    sds := &serverDomainSet{wx.NewBaseWxOpen(), "", "", make(map[string][]string)}
    sds.appId = appId
    sds.ReqContentType = project.HTTPContentTypeJSON
    sds.ReqMethod = fasthttp.MethodPost
    return sds
}
