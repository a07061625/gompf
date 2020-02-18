/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/13 0013
 * Time: 16:32
 */
package group

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type productModify struct {
    wx.BaseWxAccount
    appId       string
    groupId     int                      // 分组ID
    productList []map[string]interface{} // 商品列表
}

func (pm *productModify) SetGroupId(groupId int) {
    if groupId > 0 {
        pm.groupId = groupId
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "分组ID不合法", nil))
    }
}

func (pm *productModify) SeProductList(productList []map[string]interface{}) {
    pm.productList = productList
}

func (pm *productModify) checkData() {
    if pm.groupId <= 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "分组ID不能为空", nil))
    }
    if len(pm.productList) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商品列表不能为空", nil))
    }
}

func (pm *productModify) SendRequest() api.ApiResult {
    pm.checkData()

    reqData := make(map[string]interface{})
    reqData["group_id"] = pm.groupId
    reqData["product"] = pm.productList
    reqBody := mpf.JSONMarshal(reqData)
    pm.ReqUrl = "https://api.weixin.qq.com/merchant/group/productmod?access_token=" + wx.NewUtilWx().GetSingleAccessToken(pm.appId)
    client, req := pm.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := pm.SendInner(client, req, errorcode.WxAccountRequestPost)
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

func NewProductModify(appId string) *productModify {
    pm := &productModify{wx.NewBaseWxAccount(), "", 0, make([]map[string]interface{}, 0)}
    pm.appId = appId
    pm.ReqContentType = project.HTTPContentTypeJSON
    pm.ReqMethod = fasthttp.MethodPost
    return pm
}
