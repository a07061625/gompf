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

// 创建门店信息
type shopCreate struct {
    alipay.BaseAliPay
    storeId       string   // 门店编号
    categoryId    string   // 类目id
    mainName      string   // 主门店名
    provinceCode  string   // 省份编码
    cityCode      string   // 城市编码
    districtCode  string   // 区县编码
    address       string   // 详细地址
    longitude     string   // 坐标系经度
    latitude      string   // 坐标系纬度
    contactNumber []string // 门店电话号码
    mainImage     string   // 门店首图
    isvUid        string   // ISV返佣id
    requestId     string   // 请求ID
}

func (sc *shopCreate) SetStoreId(storeId string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{1,32}$`, storeId)
    if match {
        sc.storeId = storeId
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "门店编号不合法", nil))
    }
}

func (sc *shopCreate) SetCategoryId(categoryId string) {
    match, _ := regexp.MatchString(`^[0-9]{1,32}$`, categoryId)
    if match {
        sc.categoryId = categoryId
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "类目id不合法", nil))
    }
}

func (sc *shopCreate) SetBrandName(brandName string) {
    if (len(brandName) > 0) && (len(brandName) <= 512) {
        sc.BizContent["brand_name"] = brandName
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "品牌名不合法", nil))
    }
}

func (sc *shopCreate) SetBrandLogo(brandLogo string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{1,512}$`, brandLogo)
    if match {
        sc.BizContent["brand_logo"] = brandLogo
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "品牌LOGO不合法", nil))
    }
}

func (sc *shopCreate) SetMainName(mainName string) {
    if (len(mainName) > 0) && (len(mainName) <= 20) {
        sc.mainName = mainName
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "主门店名不合法", nil))
    }
}

func (sc *shopCreate) SetBranchName(branchName string) {
    if (len(branchName) > 0) && (len(branchName) <= 20) {
        sc.BizContent["branch_shop_name"] = branchName
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "分店名称不合法", nil))
    }
}

func (sc *shopCreate) SetProvinceCode(provinceCode string) {
    match, _ := regexp.MatchString(`^[0-9]{1,10}$`, provinceCode)
    if match {
        sc.provinceCode = provinceCode
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "省份编码不合法", nil))
    }
}

func (sc *shopCreate) SetCityCode(cityCode string) {
    match, _ := regexp.MatchString(`^[0-9]{1,10}$`, cityCode)
    if match {
        sc.cityCode = cityCode
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "城市编码不合法", nil))
    }
}

func (sc *shopCreate) SetDistrictCode(districtCode string) {
    match, _ := regexp.MatchString(`^[0-9]{1,10}$`, districtCode)
    if match {
        sc.districtCode = districtCode
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "区县编码不合法", nil))
    }
}

func (sc *shopCreate) SetAddress(address string) {
    if (len(address) >= 4) && (len(address) <= 50) {
        sc.address = address
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "详细地址不合法", nil))
    }
}

func (sc *shopCreate) SetLngAndLat(lng, lat float32) {
    if (lng < -180) || (lng > 180) {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "经度不合法", nil))
    }
    if (lat < -90) || (lat > 90) {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "纬度不合法", nil))
    }
    sc.longitude = fmt.Sprintf("%f", lng)
    sc.latitude = fmt.Sprintf("%f", lat)
}

func (sc *shopCreate) SetContactNumber(contactNumber []string) {
    if len(contactNumber) > 0 {
        sc.contactNumber = contactNumber
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "门店电话号码不合法", nil))
    }
}

func (sc *shopCreate) SetNotifyMobile(notifyMobile string) {
    match, _ := regexp.MatchString(project.RegexPhone, notifyMobile)
    if match {
        sc.BizContent["notify_mobile"] = notifyMobile
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "店长电话号码不合法", nil))
    }
}

func (sc *shopCreate) SetMainImage(mainImage string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{1,512}$`, mainImage)
    if match {
        sc.mainImage = mainImage
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "门店首图不合法", nil))
    }
}

func (sc *shopCreate) SetAuditImage(auditImage []string) {
    images := make([]string, 0)
    for _, v := range auditImage {
        match, _ := regexp.MatchString(`^[0-9a-zA-Z]{1,512}$`, v)
        if match {
            images = append(images, v)
        }
    }

    if len(images) > 0 {
        sc.BizContent["audit_images"] = strings.Join(images, ",")
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "审核图片不合法", nil))
    }
}

func (sc *shopCreate) SetBusinessTime(businessTime string) {
    if (len(businessTime) > 0) && (len(businessTime) <= 256) {
        sc.BizContent["business_time"] = businessTime
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "营业时间不合法", nil))
    }
}

func (sc *shopCreate) SetWifiFlag(wifiFlag string) {
    if (wifiFlag == "T") || (wifiFlag == "F") {
        sc.BizContent["wifi"] = wifiFlag
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "支持WIFI状态不合法", nil))
    }
}

func (sc *shopCreate) SetParkingFlag(parkingFlag string) {
    if (parkingFlag == "T") || (parkingFlag == "F") {
        sc.BizContent["parking"] = parkingFlag
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "支持停车状态不合法", nil))
    }
}

func (sc *shopCreate) SetValueAdded(valueAdded string) {
    if (len(valueAdded) > 0) && (len(valueAdded) <= 256) {
        sc.BizContent["value_added"] = valueAdded
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "其他的服务不合法", nil))
    }
}

func (sc *shopCreate) SetAvgPrice(avgPrice float32) {
    if (avgPrice >= 1) && (avgPrice <= 99999) {
        nowPrice, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", avgPrice), 64)
        sc.BizContent["avg_price"] = float32(nowPrice)
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "人均消费价格不合法", nil))
    }
}

func (sc *shopCreate) SetIsvUid(isvUid string) {
    match, _ := regexp.MatchString(`^[0-9]{1,16}$`, isvUid)
    if match {
        sc.isvUid = isvUid
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "ISV返佣id不合法", nil))
    }
}

func (sc *shopCreate) SetLicence(licence string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{1,512}$`, licence)
    if match {
        sc.BizContent["licence"] = licence
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "营业执照图片不合法", nil))
    }
}

func (sc *shopCreate) SetLicenceCode(licenceCode string) {
    if (len(licenceCode) > 0) && (len(licenceCode) <= 255) {
        sc.BizContent["licence_code"] = licenceCode
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "营业执照编号不合法", nil))
    }
}

func (sc *shopCreate) SetLicenceName(licenceName string) {
    if (len(licenceName) > 0) && (len(licenceName) <= 255) {
        sc.BizContent["licence_name"] = licenceName
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "营业执照名称不合法", nil))
    }
}

func (sc *shopCreate) SetLicenceExpire(licenceExpire string) {
    if (len(licenceExpire) > 0) && (len(licenceExpire) <= 64) {
        sc.BizContent["licence_expires"] = licenceExpire
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "营业执照过期时间不合法", nil))
    }
}

func (sc *shopCreate) SetBusinessCertificate(businessCertificate string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{1,512}$`, businessCertificate)
    if match {
        sc.BizContent["business_certificate"] = businessCertificate
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "许可证不合法", nil))
    }
}

func (sc *shopCreate) SetBusinessCertificateExpire(businessCertificateExpire string) {
    if (len(businessCertificateExpire) > 0) && (len(businessCertificateExpire) <= 64) {
        sc.BizContent["business_certificate_expires"] = businessCertificateExpire
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "许可证有效期不合法", nil))
    }
}

func (sc *shopCreate) SetAuthLetter(authLetter string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{1,512}$`, authLetter)
    if match {
        sc.BizContent["auth_letter"] = authLetter
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "授权函不合法", nil))
    }
}

func (sc *shopCreate) SetOtherOnlineFlag(otherOnlineFlag string) {
    if (otherOnlineFlag == "T") || (otherOnlineFlag == "F") {
        sc.BizContent["is_operating_online"] = otherOnlineFlag
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "其他平台开店状态不合法", nil))
    }
}

func (sc *shopCreate) SetOnlineUrl(onlineUrl []string) {
    urls := make([]string, 0)
    for _, v := range onlineUrl {
        match, _ := regexp.MatchString(project.RegexURLHTTP, v)
        if match {
            urls = append(urls, v)
        }
    }

    if len(urls) > 0 {
        sc.BizContent["online_url"] = strings.Join(urls, ",")
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "其他平台店铺链接url不合法", nil))
    }
}

func (sc *shopCreate) SetOperateNotifyUrl(operateNotifyUrl string) {
    match, _ := regexp.MatchString(project.RegexURLHTTP, operateNotifyUrl)
    if match {
        sc.BizContent["operate_notify_url"] = operateNotifyUrl
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "审核状态消息推送地址不合法", nil))
    }
}

func (sc *shopCreate) SetImplementId(implementId []string) {
    ids := make([]string, 0)
    for _, v := range implementId {
        match, _ := regexp.MatchString(project.RegexDigitAlpha, v)
        if match {
            ids = append(ids, v)
        }
    }

    if len(ids) > 0 {
        sc.BizContent["implement_id"] = strings.Join(ids, ",")
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "机具号不合法", nil))
    }
}

func (sc *shopCreate) SetNoSmokingFlag(noSmokingFlag string) {
    if (noSmokingFlag == "T") || (noSmokingFlag == "F") {
        sc.BizContent["no_smoking"] = noSmokingFlag
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "无烟区状态不合法", nil))
    }
}

func (sc *shopCreate) SetBoxFlag(boxFlag string) {
    if (boxFlag == "T") || (boxFlag == "F") {
        sc.BizContent["box"] = boxFlag
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "包厢状态不合法", nil))
    }
}

func (sc *shopCreate) SetRequestId(requestId string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{1,64}$`, requestId)
    if match {
        sc.requestId = requestId
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "请求ID不合法", nil))
    }
}

func (sc *shopCreate) SetOtherAuthorization(otherAuthorization string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{1,500}$`, otherAuthorization)
    if match {
        sc.BizContent["other_authorization"] = otherAuthorization
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "其他资质不合法", nil))
    }
}

func (sc *shopCreate) SetOpRole(opRole string) {
    if (opRole == "ISV") || (opRole == "PROVIDER") {
        sc.BizContent["op_role"] = opRole
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "操作人角色不合法", nil))
    }
}

func (sc *shopCreate) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(sc.storeId) == 0 {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "门店编号不能为空", nil))
    }
    if len(sc.categoryId) == 0 {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "类目id不能为空", nil))
    }
    if len(sc.mainName) == 0 {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "主门店名不能为空", nil))
    }
    if len(sc.provinceCode) == 0 {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "省份编码不能为空", nil))
    }
    if len(sc.cityCode) == 0 {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "城市编码不能为空", nil))
    }
    if len(sc.districtCode) == 0 {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "区县编码不能为空", nil))
    }
    if len(sc.address) == 0 {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "详细地址不能为空", nil))
    }
    if len(sc.longitude) == 0 {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "经度不能为空", nil))
    }
    if len(sc.latitude) == 0 {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "纬度不能为空", nil))
    }
    if len(sc.contactNumber) == 0 {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "门店电话号码不能为空", nil))
    }
    if len(sc.mainImage) == 0 {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "门店首图不能为空", nil))
    }
    if len(sc.isvUid) == 0 {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "ISV返佣id不能为空", nil))
    }
    if len(sc.requestId) == 0 {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "请求ID不能为空", nil))
    }
    sc.BizContent["store_id"] = sc.storeId
    sc.BizContent["category_id"] = sc.categoryId
    sc.BizContent["main_shop_name"] = sc.mainName
    sc.BizContent["province_code"] = sc.provinceCode
    sc.BizContent["city_code"] = sc.cityCode
    sc.BizContent["district_code"] = sc.districtCode
    sc.BizContent["address"] = sc.address
    sc.BizContent["longitude"] = sc.longitude
    sc.BizContent["latitude"] = sc.latitude
    sc.BizContent["contact_number"] = strings.Join(sc.contactNumber, ",")
    sc.BizContent["main_image"] = sc.mainImage
    sc.BizContent["isv_uid"] = sc.isvUid
    sc.BizContent["request_id"] = sc.requestId

    return sc.GetRequest()
}

func NewShopCreate(appId string) *shopCreate {
    sc := &shopCreate{alipay.NewBase(appId), "", "", "", "", "", "", "", "", "", make([]string, 0), "", "", ""}
    sc.BizContent["biz_version"] = "2.0"
    sc.SetMethod("alipay.offline.market.shop.create")
    sc.SetUrlNotify(true)
    return sc
}
