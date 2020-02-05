/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/13 0013
 * Time: 22:41
 */
package product

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type productUpdate struct {
    wx.BaseWxAccount
    appId          string
    productId      string                   // 商品ID
    name           string                   // 商品名称
    categoryList   []int                    // 商品分类id列表
    mainImage      string                   // 商品主图
    imageList      []string                 // 商品图片列表
    detailList     []map[string]interface{} // 商品详情列表
    propertyList   []map[string]interface{} // 商品属性列表
    skuList        []map[string]interface{} // 商品sku信息列表
    buyLimit       int                      // 商品限购数量
    productSkuList []map[string]interface{} // 产品sku信息列表
    extAttr        map[string]interface{}   // 商品其他属性
    deliveryInfo   map[string]interface{}   // 运费信息
}

func (pu *productUpdate) SetProductId(productId string) {
    if len(productId) > 0 {
        pu.productId = productId
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商品ID不合法", nil))
    }
}

func (pu *productUpdate) SetName(name string) {
    if len(name) > 0 {
        pu.name = name
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商品名称不合法", nil))
    }
}

func (pu *productUpdate) SetCategoryList(categoryList []int) {
    pu.categoryList = make([]int, 0)
    for _, v := range categoryList {
        if v > 0 {
            pu.categoryList = append(pu.categoryList, v)
        }
    }
}

func (pu *productUpdate) SetMainImage(mainImage string) {
    match, _ := regexp.MatchString(project.RegexUrlHttp, mainImage)
    if match {
        pu.mainImage = mainImage
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商品主图不合法", nil))
    }
}

func (pu *productUpdate) SetImageList(imageList []string) {
    pu.imageList = make([]string, 0)
    for _, v := range imageList {
        match, _ := regexp.MatchString(project.RegexUrlHttp, v)
        if match {
            pu.imageList = append(pu.imageList, v)
        }
    }
}

func (pu *productUpdate) SetDetailList(detailList []map[string]interface{}) {
    pu.detailList = detailList
}

func (pu *productUpdate) SetPropertyList(propertyList []map[string]interface{}) {
    pu.propertyList = propertyList
}

func (pu *productUpdate) SetSkuList(skuList []map[string]interface{}) {
    pu.skuList = skuList
}

func (pu *productUpdate) SetBuyLimit(buyLimit int) {
    if buyLimit > 0 {
        pu.buyLimit = buyLimit
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商品限购数量不合法", nil))
    }
}

func (pu *productUpdate) SetProductSkuList(productSkuList []map[string]interface{}) {
    pu.productSkuList = productSkuList
}

func (pu *productUpdate) SetExtAttr(extAttr map[string]interface{}) {
    if len(extAttr) > 0 {
        pu.extAttr = extAttr
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商品其他属性不合法", nil))
    }
}

func (pu *productUpdate) SetDeliveryInfo(deliveryInfo map[string]interface{}) {
    if len(deliveryInfo) > 0 {
        pu.deliveryInfo = deliveryInfo
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商品运费信息不合法", nil))
    }
}

func (pu *productUpdate) checkData() {
    if len(pu.productId) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商品ID不能为空", nil))
    }
    if len(pu.name) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商品名称不能为空", nil))
    }
    if len(pu.categoryList) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商品分类ID不能为空", nil))
    }
    if len(pu.mainImage) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商品主图不能为空", nil))
    }
    if len(pu.imageList) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商品图片列表不能为空", nil))
    }
    if len(pu.detailList) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商品详情列表不能为空", nil))
    }
    if len(pu.skuList) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商品sku列表不能为空", nil))
    }
    if len(pu.deliveryInfo) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商品运费信息不能为空", nil))
    }
}

func (pu *productUpdate) SendRequest() api.ApiResult {
    pu.checkData()

    baseInfo := make(map[string]interface{})
    baseInfo["name"] = pu.name
    baseInfo["main_img"] = pu.mainImage
    baseInfo["buy_limit"] = pu.buyLimit
    baseInfo["category_id"] = pu.categoryList
    baseInfo["img"] = pu.imageList
    baseInfo["detail"] = pu.detailList
    if len(pu.propertyList) > 0 {
        baseInfo["property"] = pu.propertyList
    }
    if len(pu.productSkuList) > 0 {
        baseInfo["sku_info"] = pu.productSkuList
    }
    reqData := make(map[string]interface{})
    reqData["product_id"] = pu.productId
    reqData["product_base"] = baseInfo
    reqData["delivery_info"] = pu.deliveryInfo
    reqData["sku_list"] = pu.skuList
    if len(pu.extAttr) > 0 {
        reqData["attrext"] = pu.extAttr
    }
    reqBody := mpf.JsonMarshal(reqData)
    pu.ReqUrl = "https://api.weixin.qq.com/merchant/update?access_token=" + wx.NewUtilWx().GetSingleAccessToken(pu.appId)
    client, req := pu.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := pu.SendInner(client, req, errorcode.WxAccountRequestPost)
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

func NewProductUpdate(appId string) *productUpdate {
    pu := &productUpdate{wx.NewBaseWxAccount(), "", "", "", make([]int, 0), "", make([]string, 0), make([]map[string]interface{}, 0), make([]map[string]interface{}, 0), make([]map[string]interface{}, 0), 0, make([]map[string]interface{}, 0), make(map[string]interface{}), make(map[string]interface{})}
    pu.appId = appId
    pu.ReqContentType = project.HttpContentTypeJson
    pu.ReqMethod = fasthttp.MethodPost
    return pu
}
