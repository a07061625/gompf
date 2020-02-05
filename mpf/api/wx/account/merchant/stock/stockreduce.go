/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/13 0013
 * Time: 16:48
 */
package stock

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type stockReduce struct {
    wx.BaseWxAccount
    appId     string
    productId string // 商品ID
    skuInfo   string // sku信息,格式"id1:vid1;id2:vid2"
    quantity  int    // 库存数量
}

func (sr *stockReduce) SetProductId(productId string) {
    if len(productId) > 0 {
        sr.productId = productId
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商品ID不合法", nil))
    }
}

func (sr *stockReduce) SetSkuInfo(skuInfo string) {
    if len(skuInfo) > 0 {
        sr.skuInfo = skuInfo
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "sku信息不合法", nil))
    }
}

func (sr *stockReduce) SetQuantity(quantity int) {
    if quantity > 0 {
        sr.quantity = quantity
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "库存数量不合法", nil))
    }
}

func (sr *stockReduce) checkData() {
    if len(sr.productId) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商品ID不能为空", nil))
    }
    if len(sr.skuInfo) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "sku信息不能为空", nil))
    }
    if sr.quantity <= 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "库存数量不能为空", nil))
    }
}

func (sr *stockReduce) SendRequest() api.ApiResult {
    sr.checkData()

    reqData := make(map[string]interface{})
    reqData["product_id"] = sr.productId
    reqData["sku_info"] = sr.skuInfo
    reqData["quantity"] = sr.quantity
    reqBody := mpf.JsonMarshal(reqData)
    sr.ReqUrl = "https://api.weixin.qq.com/merchant/stock/reduce?access_token=" + wx.NewUtilWx().GetSingleAccessToken(sr.appId)
    client, req := sr.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := sr.SendInner(client, req, errorcode.WxAccountRequestPost)
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

func NewStockReduce(appId string) *stockReduce {
    sr := &stockReduce{wx.NewBaseWxAccount(), "", "", "", 0}
    sr.appId = appId
    sr.ReqContentType = project.HttpContentTypeJson
    sr.ReqMethod = fasthttp.MethodPost
    return sr
}
