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

type productCreate struct {
    wx.BaseWxAccount
    appId          string
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

func (pc *productCreate) SetName(name string) {
    if len(name) > 0 {
        pc.name = name
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商品名称不合法", nil))
    }
}

func (pc *productCreate) SetCategoryList(categoryList []int) {
    pc.categoryList = make([]int, 0)
    for _, v := range categoryList {
        if v > 0 {
            pc.categoryList = append(pc.categoryList, v)
        }
    }
}

func (pc *productCreate) SetMainImage(mainImage string) {
    match, _ := regexp.MatchString(project.RegexUrlHttp, mainImage)
    if match {
        pc.mainImage = mainImage
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商品主图不合法", nil))
    }
}

func (pc *productCreate) SetImageList(imageList []string) {
    pc.imageList = make([]string, 0)
    for _, v := range imageList {
        match, _ := regexp.MatchString(project.RegexUrlHttp, v)
        if match {
            pc.imageList = append(pc.imageList, v)
        }
    }
}

func (pc *productCreate) SetDetailList(detailList []map[string]interface{}) {
    pc.detailList = detailList
}

func (pc *productCreate) SetPropertyList(propertyList []map[string]interface{}) {
    pc.propertyList = propertyList
}

func (pc *productCreate) SetSkuList(skuList []map[string]interface{}) {
    pc.skuList = skuList
}

func (pc *productCreate) SetBuyLimit(buyLimit int) {
    if buyLimit > 0 {
        pc.buyLimit = buyLimit
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商品限购数量不合法", nil))
    }
}

func (pc *productCreate) SetProductSkuList(productSkuList []map[string]interface{}) {
    pc.productSkuList = productSkuList
}

func (pc *productCreate) SetExtAttr(extAttr map[string]interface{}) {
    if len(extAttr) > 0 {
        pc.extAttr = extAttr
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商品其他属性不合法", nil))
    }
}

func (pc *productCreate) SetDeliveryInfo(deliveryInfo map[string]interface{}) {
    if len(deliveryInfo) > 0 {
        pc.deliveryInfo = deliveryInfo
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商品运费信息不合法", nil))
    }
}

func (pc *productCreate) checkData() {
    if len(pc.name) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商品名称不能为空", nil))
    }
    if len(pc.categoryList) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商品分类ID不能为空", nil))
    }
    if len(pc.mainImage) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商品主图不能为空", nil))
    }
    if len(pc.imageList) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商品图片列表不能为空", nil))
    }
    if len(pc.detailList) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商品详情列表不能为空", nil))
    }
    if len(pc.skuList) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商品sku列表不能为空", nil))
    }
    if len(pc.deliveryInfo) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商品运费信息不能为空", nil))
    }
}

func (pc *productCreate) SendRequest() api.ApiResult {
    pc.checkData()

    baseInfo := make(map[string]interface{})
    baseInfo["name"] = pc.name
    baseInfo["main_img"] = pc.mainImage
    baseInfo["buy_limit"] = pc.buyLimit
    baseInfo["category_id"] = pc.categoryList
    baseInfo["img"] = pc.imageList
    baseInfo["detail"] = pc.detailList
    if len(pc.propertyList) > 0 {
        baseInfo["property"] = pc.propertyList
    }
    if len(pc.productSkuList) > 0 {
        baseInfo["sku_info"] = pc.productSkuList
    }
    reqData := make(map[string]interface{})
    reqData["product_base"] = baseInfo
    reqData["delivery_info"] = pc.deliveryInfo
    reqData["sku_list"] = pc.skuList
    if len(pc.extAttr) > 0 {
        reqData["attrext"] = pc.extAttr
    }
    reqBody := mpf.JsonMarshal(reqData)
    pc.ReqUrl = "https://api.weixin.qq.com/merchant/create?access_token=" + wx.NewUtilWx().GetSingleAccessToken(pc.appId)
    client, req := pc.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := pc.SendInner(client, req, errorcode.WxAccountRequestPost)
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

func NewProductCreate(appId string) *productCreate {
    pc := &productCreate{wx.NewBaseWxAccount(), "", "", make([]int, 0), "", make([]string, 0), make([]map[string]interface{}, 0), make([]map[string]interface{}, 0), make([]map[string]interface{}, 0), 0, make([]map[string]interface{}, 0), make(map[string]interface{}), make(map[string]interface{})}
    pc.appId = appId
    pc.ReqContentType = project.HttpContentTypeJson
    pc.ReqMethod = fasthttp.MethodPost
    return pc
}
