/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/10 0010
 * Time: 22:33
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

// 设置小程序隐私设置（是否可被搜索）
type searchStatusChange struct {
    wx.BaseWxOpen
    appId        string // 应用ID
    searchStatus int    // 搜索状态
}

func (ssc *searchStatusChange) SetSearchStatus(searchStatus int) {
    if (searchStatus == 0) || (searchStatus == 1) {
        ssc.searchStatus = searchStatus
    } else {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "搜索状态不合法", nil))
    }
}

func (ssc *searchStatusChange) checkData() {
    if ssc.searchStatus < 0 {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "搜索状态不能为空", nil))
    }
}

func (ssc *searchStatusChange) SendRequest() api.APIResult {
    ssc.checkData()

    reqData := make(map[string]interface{})
    reqData["status"] = ssc.searchStatus
    reqBody := mpf.JSONMarshal(ssc.ReqData)
    ssc.ReqURI = "https://api.weixin.qq.com/wxa/changewxasearchstatus?access_token=" + wx.NewUtilWx().GetOpenAuthorizeAccessToken(ssc.appId)
    client, req := ssc.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := ssc.SendInner(client, req, errorcode.WxOpenRequestPost)
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

func NewSearchStatusChange(appId string) *searchStatusChange {
    ssc := &searchStatusChange{wx.NewBaseWxOpen(), "", 0}
    ssc.appId = appId
    ssc.searchStatus = -1
    ssc.ReqContentType = project.HTTPContentTypeJSON
    ssc.ReqMethod = fasthttp.MethodPost
    return ssc
}
