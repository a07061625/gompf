/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/13 0013
 * Time: 23:33
 */
package product

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type statusModify struct {
    wx.BaseWxAccount
    appId         string
    productId     string // 商品ID
    productStatus int    // 商品上下架状态(0-下架 1-上架)
}

func (sm *statusModify) SetProductId(productId string) {
    if len(productId) > 0 {
        sm.productId = productId
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商品ID不合法", nil))
    }
}

func (sm *statusModify) SetProductStatus(productStatus int) {
    if (productStatus == 0) || (productStatus == 1) {
        sm.productStatus = productStatus
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "上下架状态不合法", nil))
    }
}

func (sm *statusModify) checkData() {
    if len(sm.productId) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商品ID不能为空", nil))
    }
    if sm.productStatus < 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "上下架状态不能为空", nil))
    }
}

func (sm *statusModify) SendRequest() api.APIResult {
    sm.checkData()

    reqData := make(map[string]interface{})
    reqData["product_id"] = sm.productId
    reqData["status"] = sm.productStatus
    reqBody := mpf.JSONMarshal(reqData)
    sm.ReqURI = "https://api.weixin.qq.com/merchant/modproductstatus?access_token=" + wx.NewUtilWx().GetSingleAccessToken(sm.appId)
    client, req := sm.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := sm.SendInner(client, req, errorcode.WxAccountRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxAccountRequestPost
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewStatusModify(appId string) *statusModify {
    sm := &statusModify{wx.NewBaseWxAccount(), "", "", 0}
    sm.appId = appId
    sm.productStatus = -1
    sm.ReqContentType = project.HTTPContentTypeJSON
    sm.ReqMethod = fasthttp.MethodPost
    return sm
}
