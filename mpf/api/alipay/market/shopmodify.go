/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/19 0019
 * Time: 18:45
 */
package market

import (
    "fmt"
    "regexp"
    "strconv"
    "strings"

    "github.com/a07061625/gompf/mpf/api/alipay"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 修改门店信息
type shopModify struct {
    alipay.BaseAliPay
    shopId string // 门店ID
}

func (sm *shopModify) SetShopId(shopId string) {
    match, _ := regexp.MatchString(`^[0-9]{1,64}$`, shopId)
    if match {
        sm.shopId = shopId
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "门店ID不合法", nil))
    }
}

func (sm *shopModify) SetStoreId(storeId string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{1,32}$`, storeId)
    if match {
        sm.BizContent["store_id"] = storeId
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "门店编号不合法", nil))
    }
}

func (sm *shopModify) SetBrandName(brandName string) {
    if (len(brandName) > 0) && (len(brandName) <= 512) {
        sm.BizContent["brand_name"] = brandName
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "品牌名不合法", nil))
    }
}

func (sm *shopModify) SetBrandLogo(brandLogo string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{1,512}$`, brandLogo)
    if match {
        sm.BizContent["brand_logo"] = brandLogo
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "品牌LOGO不合法", nil))
    }
}

func (sm *shopModify) SetMainName(mainName string) {
    if (len(mainName) > 0) && (len(mainName) <= 20) {
        sm.BizContent["main_shop_name"] = mainName
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "主门店名不合法", nil))
    }
}

func (sm *shopModify) SetBranchName(branchName string) {
    if (len(branchName) > 0) && (len(branchName) <= 20) {
        sm.BizContent["branch_shop_name"] = branchName
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "分店名称不合法", nil))
    }
}

func (sm *shopModify) SetProvinceCode(provinceCode string) {
    match, _ := regexp.MatchString(`^[0-9]{1,10}$`, provinceCode)
    if match {
        sm.BizContent["province_code"] = provinceCode
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "省份编码不合法", nil))
    }
}

func (sm *shopModify) SetCityCode(cityCode string) {
    match, _ := regexp.MatchString(`^[0-9]{1,10}$`, cityCode)
    if match {
        sm.BizContent["city_code"] = cityCode
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "城市编码不合法", nil))
    }
}

func (sm *shopModify) SetDistrictCode(districtCode string) {
    match, _ := regexp.MatchString(`^[0-9]{1,10}$`, districtCode)
    if match {
        sm.BizContent["district_code"] = districtCode
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "区县编码不合法", nil))
    }
}

func (sm *shopModify) SetAddress(address string) {
    if (len(address) >= 4) && (len(address) <= 50) {
        sm.BizContent["address"] = address
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "详细地址不合法", nil))
    }
}

func (sm *shopModify) SetLngAndLat(lng, lat float32) {
    if (lng < -180) || (lng > 180) {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "经度不合法", nil))
    }
    if (lat < -90) || (lat > 90) {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "纬度不合法", nil))
    }
    sm.BizContent["longitude"] = fmt.Sprintf("%f", lng)
    sm.BizContent["latitude"] = fmt.Sprintf("%f", lat)
}

func (sm *shopModify) SetContactNumber(contactNumber []string) {
    if len(contactNumber) > 0 {
        sm.BizContent["contact_number"] = strings.Join(contactNumber, ",")
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "门店电话号码不合法", nil))
    }
}

func (sm *shopModify) SetNotifyMobile(notifyMobile string) {
    match, _ := regexp.MatchString(project.RegexPhone, notifyMobile)
    if match {
        sm.BizContent["notify_mobile"] = notifyMobile
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "店长电话号码不合法", nil))
    }
}

func (sm *shopModify) SetMainImage(mainImage string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{1,512}$`, mainImage)
    if match {
        sm.BizContent["main_image"] = mainImage
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "门店首图不合法", nil))
    }
}

func (sm *shopModify) SetAuditImage(auditImage []string) {
    images := make([]string, 0)
    for _, v := range auditImage {
        match, _ := regexp.MatchString(`^[0-9a-zA-Z]{1,512}$`, v)
        if match {
            images = append(images, v)
        }
    }

    if len(images) > 0 {
        sm.BizContent["audit_images"] = strings.Join(images, ",")
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "审核图片不合法", nil))
    }
}

func (sm *shopModify) SetBusinessTime(businessTime string) {
    if (len(businessTime) > 0) && (len(businessTime) <= 256) {
        sm.BizContent["business_time"] = businessTime
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "营业时间不合法", nil))
    }
}

func (sm *shopModify) SetWifiFlag(wifiFlag string) {
    if (wifiFlag == "T") || (wifiFlag == "F") {
        sm.BizContent["wifi"] = wifiFlag
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "支持WIFI状态不合法", nil))
    }
}

func (sm *shopModify) SetParkingFlag(parkingFlag string) {
    if (parkingFlag == "T") || (parkingFlag == "F") {
        sm.BizContent["parking"] = parkingFlag
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "支持停车状态不合法", nil))
    }
}

func (sm *shopModify) SetValueAdded(valueAdded string) {
    if (len(valueAdded) > 0) && (len(valueAdded) <= 256) {
        sm.BizContent["value_added"] = valueAdded
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "其他的服务不合法", nil))
    }
}

func (sm *shopModify) SetAvgPrice(avgPrice float32) {
    if (avgPrice >= 1) && (avgPrice <= 99999) {
        nowPrice, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", avgPrice), 64)
        sm.BizContent["avg_price"] = float32(nowPrice)
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "人均消费价格不合法", nil))
    }
}

func (sm *shopModify) SetLicence(licence string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{1,512}$`, licence)
    if match {
        sm.BizContent["licence"] = licence
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "营业执照图片不合法", nil))
    }
}

func (sm *shopModify) SetLicenceCode(licenceCode string) {
    if (len(licenceCode) > 0) && (len(licenceCode) <= 255) {
        sm.BizContent["licence_code"] = licenceCode
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "营业执照编号不合法", nil))
    }
}

func (sm *shopModify) SetLicenceName(licenceName string) {
    if (len(licenceName) > 0) && (len(licenceName) <= 255) {
        sm.BizContent["licence_name"] = licenceName
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "营业执照名称不合法", nil))
    }
}

func (sm *shopModify) SetLicenceExpire(licenceExpire string) {
    if (len(licenceExpire) > 0) && (len(licenceExpire) <= 64) {
        sm.BizContent["licence_expires"] = licenceExpire
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "营业执照过期时间不合法", nil))
    }
}

func (sm *shopModify) SetBusinessCertificate(businessCertificate string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{1,512}$`, businessCertificate)
    if match {
        sm.BizContent["business_certificate"] = businessCertificate
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "许可证不合法", nil))
    }
}

func (sm *shopModify) SetBusinessCertificateExpire(businessCertificateExpire string) {
    if (len(businessCertificateExpire) > 0) && (len(businessCertificateExpire) <= 64) {
        sm.BizContent["business_certificate_expires"] = businessCertificateExpire
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "许可证有效期不合法", nil))
    }
}

func (sm *shopModify) SetAuthLetter(authLetter string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{1,512}$`, authLetter)
    if match {
        sm.BizContent["auth_letter"] = authLetter
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "授权函不合法", nil))
    }
}

func (sm *shopModify) SetOtherOnlineFlag(otherOnlineFlag string) {
    if (otherOnlineFlag == "T") || (otherOnlineFlag == "F") {
        sm.BizContent["is_operating_online"] = otherOnlineFlag
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "其他平台开店状态不合法", nil))
    }
}

func (sm *shopModify) SetOnlineUrl(onlineUrl []string) {
    urls := make([]string, 0)
    for _, v := range onlineUrl {
        match, _ := regexp.MatchString(project.RegexURLHTTP, v)
        if match {
            urls = append(urls, v)
        }
    }

    if len(urls) > 0 {
        sm.BizContent["online_url"] = strings.Join(urls, ",")
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "其他平台店铺链接url不合法", nil))
    }
}

func (sm *shopModify) SetOperateNotifyUrl(operateNotifyUrl string) {
    match, _ := regexp.MatchString(project.RegexURLHTTP, operateNotifyUrl)
    if match {
        sm.BizContent["operate_notify_url"] = operateNotifyUrl
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "审核状态消息推送地址不合法", nil))
    }
}

func (sm *shopModify) SetImplementId(implementId []string) {
    ids := make([]string, 0)
    for _, v := range implementId {
        match, _ := regexp.MatchString(project.RegexDigitAlpha, v)
        if match {
            ids = append(ids, v)
        }
    }

    if len(ids) > 0 {
        sm.BizContent["implement_id"] = strings.Join(ids, ",")
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "机具号不合法", nil))
    }
}

func (sm *shopModify) SetNoSmokingFlag(noSmokingFlag string) {
    if (noSmokingFlag == "T") || (noSmokingFlag == "F") {
        sm.BizContent["no_smoking"] = noSmokingFlag
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "无烟区状态不合法", nil))
    }
}

func (sm *shopModify) SetBoxFlag(boxFlag string) {
    if (boxFlag == "T") || (boxFlag == "F") {
        sm.BizContent["box"] = boxFlag
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "包厢状态不合法", nil))
    }
}

func (sm *shopModify) SetRequestId(requestId string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{1,64}$`, requestId)
    if match {
        sm.BizContent["request_id"] = requestId
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "请求ID不合法", nil))
    }
}

func (sm *shopModify) SetOtherAuthorization(otherAuthorization string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{1,500}$`, otherAuthorization)
    if match {
        sm.BizContent["other_authorization"] = otherAuthorization
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "其他资质不合法", nil))
    }
}

func (sm *shopModify) SetOpRole(opRole string) {
    if (opRole == "ISV") || (opRole == "PROVIDER") {
        sm.BizContent["op_role"] = opRole
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "操作人角色不合法", nil))
    }
}

func (sm *shopModify) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(sm.shopId) == 0 {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "门店ID不能为空", nil))
    }
    sm.BizContent["shop_id"] = sm.shopId

    return sm.GetRequest()
}

func NewShopModify(appId string) *shopModify {
    sm := &shopModify{alipay.NewBase(appId), ""}
    sm.BizContent["biz_version"] = "2.0"
    sm.SetMethod("alipay.offline.market.shop.modify")
    sm.SetUrlNotify(true)
    return sm
}
