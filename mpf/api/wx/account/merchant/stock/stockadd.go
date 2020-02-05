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

type stockAdd struct {
    wx.BaseWxAccount
    appId     string
    productId string // 商品ID
    skuInfo   string // sku信息,格式"id1:vid1;id2:vid2"
    quantity  int    // 库存数量
}

func (sa *stockAdd) SetProductId(productId string) {
    if len(productId) > 0 {
        sa.productId = productId
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商品ID不合法", nil))
    }
}

func (sa *stockAdd) SetSkuInfo(skuInfo string) {
    if len(skuInfo) > 0 {
        sa.skuInfo = skuInfo
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "sku信息不合法", nil))
    }
}

func (sa *stockAdd) SetQuantity(quantity int) {
    if quantity > 0 {
        sa.quantity = quantity
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "库存数量不合法", nil))
    }
}

func (sa *stockAdd) checkData() {
    if len(sa.productId) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商品ID不能为空", nil))
    }
    if len(sa.skuInfo) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "sku信息不能为空", nil))
    }
    if sa.quantity <= 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "库存数量不能为空", nil))
    }
}

func (sa *stockAdd) SendRequest() api.ApiResult {
    sa.checkData()

    reqData := make(map[string]interface{})
    reqData["product_id"] = sa.productId
    reqData["sku_info"] = sa.skuInfo
    reqData["quantity"] = sa.quantity
    reqBody := mpf.JsonMarshal(reqData)
    sa.ReqUrl = "https://api.weixin.qq.com/merchant/stock/add?access_token=" + wx.NewUtilWx().GetSingleAccessToken(sa.appId)
    client, req := sa.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := sa.SendInner(client, req, errorcode.WxAccountRequestPost)
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

func NewStockAdd(appId string) *stockAdd {
    sa := &stockAdd{wx.NewBaseWxAccount(), "", "", "", 0}
    sa.appId = appId
    sa.ReqContentType = project.HttpContentTypeJson
    sa.ReqMethod = fasthttp.MethodPost
    return sa
}
