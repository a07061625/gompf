/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/13 0013
 * Time: 23:28
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

type listByStatus struct {
    wx.BaseWxAccount
    appId         string
    productStatus int // 商品状态(0-全部 1-上架 2-下架)
}

func (lbs *listByStatus) SetProductStatus(productStatus int) {
    if (productStatus >= 0) && (productStatus <= 2) {
        lbs.productStatus = productStatus
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商品状态不合法", nil))
    }
}

func (lbs *listByStatus) SendRequest() api.ApiResult {
    reqData := make(map[string]interface{})
    reqData["status"] = lbs.productStatus
    reqBody := mpf.JSONMarshal(reqData)
    lbs.ReqUrl = "https://api.weixin.qq.com/merchant/getbystatus?access_token=" + wx.NewUtilWx().GetSingleAccessToken(lbs.appId)
    client, req := lbs.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := lbs.SendInner(client, req, errorcode.WxAccountRequestPost)
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

func NewListByStatus(appId string) *listByStatus {
    lbs := &listByStatus{wx.NewBaseWxAccount(), "", 0}
    lbs.appId = appId
    lbs.productStatus = 0
    lbs.ReqContentType = project.HTTPContentTypeJSON
    lbs.ReqMethod = fasthttp.MethodPost
    return lbs
}
