/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/10 0010
 * Time: 22:45
 */
package account

import (
    "strconv"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
)

// 获取可以用来设置的公众号列表
type showItemList struct {
    wx.BaseWxOpen
    appId string // 应用ID
    page  int    // 页数
    num   int    // 每页记录数
}

func (sil *showItemList) SetPage(page int) {
    if page >= 0 {
        sil.ReqData["page"] = strconv.Itoa(page)
    } else {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "页数不合法", nil))
    }
}

func (sil *showItemList) SetNum(num int) {
    if (num > 0) && (num <= 20) {
        sil.ReqData["num"] = strconv.Itoa(num)
    } else {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "每页记录数不合法", nil))
    }
}

func (sil *showItemList) SendRequest() api.APIResult {
    sil.ReqData["access_token"] = wx.NewUtilWx().GetOpenAuthorizeAccessToken(sil.appId)
    sil.ReqURI = "https://api.weixin.qq.com/wxa/getwxamplinkforshow?" + mpf.HTTPCreateParams(sil.ReqData, "none", 1)
    client, req := sil.GetRequest()

    resp, result := sil.SendInner(client, req, errorcode.WxOpenRequestGet)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxOpenRequestGet
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewShowItemList(appId string) *showItemList {
    sil := &showItemList{wx.NewBaseWxOpen(), "", 0, 0}
    sil.appId = appId
    sil.ReqData["page"] = "0"
    sil.ReqData["num"] = "10"
    return sil
}
