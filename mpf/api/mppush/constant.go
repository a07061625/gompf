/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/2 0002
 * Time: 22:27
 */
package mppush

const (
    BaiDuDeviceTypeAll     = "0"
    BaiDuDeviceTypeAndroid = "3"
    BaiDuDeviceTypeIOS     = "4"
    BaiDuServiceDomain     = "https://api.tuisong.baidu.com"
    BaiDuServiceUriPrefix  = "/rest/3.0"

    XinGePlatformTypeAll     = "all"
    XinGePlatformTypeAndroid = "android"
    XinGePlatformTypeIOS     = "ios"
    XinGeServiceDomain       = "https://openapi.xg.qq.com/v3/"

    JPushPlatformTypeAll      = "all"
    JPushPlatformTypeAndroid  = "android"
    JPushPlatformTypeIOS      = "ios"
    JPushPlatformTypeWinPhone = "winphone"
    JPushServiceDomainAdmin   = "https://admin.jpush.cn"
    JPushServiceDomainDevice  = "https://device.jpush.cn"
    JPushServiceDomainApi     = "https://api.jpush.cn"
    JPushServiceDomainReport  = "https://report.jpush.cn"
)

var (
    BaiDuDeviceTypes   map[string]string
    XinGePlatformTypes map[string]string
    JPushPlatformTypes map[string]string
)

func init() {
    BaiDuDeviceTypes = make(map[string]string)
    BaiDuDeviceTypes[BaiDuDeviceTypeAndroid] = "安卓"
    BaiDuDeviceTypes[BaiDuDeviceTypeIOS] = "苹果"

    XinGePlatformTypes = make(map[string]string)
    XinGePlatformTypes[XinGePlatformTypeAndroid] = "安卓"
    XinGePlatformTypes[XinGePlatformTypeIOS] = "苹果"

    JPushPlatformTypes = make(map[string]string)
    JPushPlatformTypes[JPushPlatformTypeAndroid] = "安卓"
    JPushPlatformTypes[JPushPlatformTypeIOS] = "苹果"
    JPushPlatformTypes[JPushPlatformTypeWinPhone] = "微软手机"
}
