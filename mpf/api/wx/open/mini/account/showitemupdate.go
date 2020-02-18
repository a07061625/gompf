/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/10 0010
 * Time: 22:52
 */
package account

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 设置展示的公众号
type showItemUpdate struct {
    wx.BaseWxOpen
    appId   string // 应用ID
    bizFlag int    // 展示公众号开启标识 0:关闭 1:开启
}

func (siu *showItemUpdate) SetBizFlag(bizFlag int) {
    if (bizFlag == 0) || (bizFlag == 1) {
        siu.bizFlag = bizFlag
    } else {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "展示公众号开启标识不合法", nil))
    }
}

func (siu *showItemUpdate) SendRequest() api.ApiResult {
    reqData := make(map[string]interface{})
    reqData["wxa_subscribe_biz_flag"] = siu.bizFlag
    if siu.bizFlag == 1 {
        reqData["appid"] = siu.appId
    }
    reqBody := mpf.JsonMarshal(reqData)
    siu.ReqUrl = "https://api.weixin.qq.com/wxa/updateshowwxaitem?access_token=" + wx.NewUtilWx().GetOpenAuthorizeAccessToken(siu.appId)
    client, req := siu.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := siu.SendInner(client, req, errorcode.WxOpenRequestPost)
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

func NewShowItemUpdate(appId string) *showItemUpdate {
    siu := &showItemUpdate{wx.NewBaseWxOpen(), "", 0}
    siu.appId = appId
    siu.bizFlag = 0
    siu.ReqContentType = project.HTTPContentTypeJSON
    siu.ReqMethod = fasthttp.MethodPost
    return siu
}
