/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/22 0022
 * Time: 16:41
 */
package aliopen

import (
    "regexp"

    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
)

// 大鱼短信配置
type configDySms struct {
    regionId  string // 区域ID
    appKey    string // app key
    appSecret string // app密钥
}

func (c *configDySms) SetRegionId(regionId string) {
    if len(regionId) > 0 {
        c.regionId = regionId
    } else {
        panic(mperr.NewSmsAliYun(errorcode.SmsAliYunParam, "区域ID不合法", nil))
    }
}

func (c *configDySms) GetRegionId() string {
    return c.regionId
}

func (c *configDySms) SetAppKey(appKey string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, appKey)
    if match {
        c.appKey = appKey
    } else {
        panic(mperr.NewSmsAliYun(errorcode.SmsAliYunParam, "app key不合法", nil))
    }
}

func (c *configDySms) GetAppKey() string {
    return c.appKey
}

func (c *configDySms) SetAppSecret(appSecret string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, appSecret)
    if match {
        c.appSecret = appSecret
    } else {
        panic(mperr.NewSmsAliYun(errorcode.SmsAliYunParam, "app密钥不合法", nil))
    }
}

func (c *configDySms) GetAppSecret() string {
    return c.appSecret
}
