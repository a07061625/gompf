/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/14 0014
 * Time: 19:39
 */
package pay

import (
    "regexp"
    "strconv"
    "time"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
)

type payNativePre struct {
    wx.BaseWxAccount
    appId     string
    productId string // 商户定义的商品id
}

func (pnp *payNativePre) SetProductId(productId string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{1,32}$`, productId)
    if match {
        pnp.productId = productId
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商品ID不合法", nil))
    }
}

func (pnp *payNativePre) checkData() {
    if len(pnp.productId) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商品ID不能为空", nil))
    }
    pnp.ReqData["product_id"] = pnp.productId
}

func (pnp *payNativePre) GetResult() map[string]string {
    pnp.checkData()

    sign := wx.NewUtilWx().CreateSinglePaySign(pnp.ReqData, pnp.appId, "md5")
    pnp.ReqData["sign"] = sign

    result := make(map[string]string)
    result["url"] = "weixin://wxpay/bizpayurl?" + mpf.HttpCreateParams(pnp.ReqData, "none", 4)
    return result
}

func NewPayNativePre(appId string) *payNativePre {
    conf := wx.NewConfig().GetAccount(appId)
    pnp := &payNativePre{wx.NewBaseWxAccount(), "", ""}
    pnp.appId = appId
    pnp.ReqData["appid"] = conf.GetAppId()
    pnp.ReqData["mch_id"] = conf.GetPayMchId()
    pnp.ReqData["time_stamp"] = strconv.Itoa(time.Now().Second())
    pnp.ReqData["nonce_str"] = mpf.ToolCreateNonceStr(32, "numlower")
    return pnp
}
