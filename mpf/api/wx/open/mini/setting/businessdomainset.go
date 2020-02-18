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

// 设置小程序业务域名
type businessDomainSet struct {
    wx.BaseWxOpen
    appId  string // 应用ID
    action string
    data   map[string][]string
}

func (bds *businessDomainSet) SetData(action string, data map[string][]string) {
    switch action {
    case "add":
        if len(data) == 0 {
            panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "域名不能为空", nil))
        }
        bds.data = data
    case "delete":
        if len(data) == 0 {
            panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "域名不能为空", nil))
        }
        bds.data = data
    case "set":
        if len(data) == 0 {
            panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "域名不能为空", nil))
        }
        bds.data = data
    case "get":
    default:
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "操作类型不支持", nil))
    }
    bds.action = action
}

func (bds *businessDomainSet) checkData() {
    if len(bds.action) == 0 {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "操作类型不能为空", nil))
    }
}

func (bds *businessDomainSet) SendRequest() api.APIResult {
    bds.checkData()

    reqData := make(map[string]interface{})
    if len(bds.data) > 0 {
        for k, v := range bds.data {
            reqData[k] = v
        }
    }
    reqData["action"] = bds.action
    reqBody := mpf.JSONMarshal(reqData)
    bds.ReqURI = "https://api.weixin.qq.com/wxa/setwebviewdomain?access_token=" + wx.NewUtilWx().GetOpenAuthorizeAccessToken(bds.appId)
    client, req := bds.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := bds.SendInner(client, req, errorcode.WxOpenRequestPost)
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

func NewBusinessDomainSet(appId string) *businessDomainSet {
    bds := &businessDomainSet{wx.NewBaseWxOpen(), "", "", make(map[string][]string)}
    bds.appId = appId
    bds.ReqContentType = project.HTTPContentTypeJSON
    bds.ReqMethod = fasthttp.MethodPost
    return bds
}
