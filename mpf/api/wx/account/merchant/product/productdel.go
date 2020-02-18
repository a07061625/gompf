/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/13 0013
 * Time: 23:21
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

type productDel struct {
    wx.BaseWxAccount
    appId     string
    productId string // 商品ID
}

func (pd *productDel) SetProductId(productId string) {
    if len(productId) > 0 {
        pd.productId = productId
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商品ID不合法", nil))
    }
}

func (pd *productDel) checkData() {
    if len(pd.productId) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商品ID不能为空", nil))
    }
    pd.ReqData["product_id"] = pd.productId
}

func (pd *productDel) SendRequest() api.ApiResult {
    pd.checkData()

    reqBody := mpf.JSONMarshal(pd.ReqData)
    pd.ReqUrl = "https://api.weixin.qq.com/merchant/del?access_token=" + wx.NewUtilWx().GetSingleAccessToken(pd.appId)
    client, req := pd.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := pd.SendInner(client, req, errorcode.WxAccountRequestPost)
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

func NewProductDel(appId string) *productDel {
    pd := &productDel{wx.NewBaseWxAccount(), "", ""}
    pd.appId = appId
    pd.ReqContentType = project.HTTPContentTypeJSON
    pd.ReqMethod = fasthttp.MethodPost
    return pd
}
