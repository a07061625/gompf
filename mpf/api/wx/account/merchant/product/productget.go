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

type productGet struct {
    wx.BaseWxAccount
    appId     string
    productId string // 商品ID
}

func (pg *productGet) SetProductId(productId string) {
    if len(productId) > 0 {
        pg.productId = productId
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商品ID不合法", nil))
    }
}

func (pg *productGet) checkData() {
    if len(pg.productId) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商品ID不能为空", nil))
    }
    pg.ReqData["product_id"] = pg.productId
}

func (pg *productGet) SendRequest() api.APIResult {
    pg.checkData()

    reqBody := mpf.JSONMarshal(pg.ReqData)
    pg.ReqURI = "https://api.weixin.qq.com/merchant/get?access_token=" + wx.NewUtilWx().GetSingleAccessToken(pg.appId)
    client, req := pg.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := pg.SendInner(client, req, errorcode.WxAccountRequestPost)
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

func NewProductGet(appId string) *productGet {
    pg := &productGet{wx.NewBaseWxAccount(), "", ""}
    pg.appId = appId
    pg.ReqContentType = project.HTTPContentTypeJSON
    pg.ReqMethod = fasthttp.MethodPost
    return pg
}
