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

// 消息推送配置
type configPush struct {
    regionId  string // 区域ID
    appKey    string // app key
    appSecret string // app密钥
}

func (c *configPush) SetRegionId(regionId string) {
    if len(regionId) > 0 {
        c.regionId = regionId
    } else {
        panic(mperr.NewPushAliYun(errorcode.PushAliYunParam, "区域ID不合法", nil))
    }
}

func (c *configPush) GetRegionId() string {
    return c.regionId
}

func (c *configPush) SetAppKey(appKey string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, appKey)
    if match {
        c.appKey = appKey
    } else {
        panic(mperr.NewPushAliYun(errorcode.PushAliYunParam, "app key不合法", nil))
    }
}

func (c *configPush) GetAppKey() string {
    return c.appKey
}

func (c *configPush) SetAppSecret(appSecret string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, appSecret)
    if match {
        c.appSecret = appSecret
    } else {
        panic(mperr.NewPushAliYun(errorcode.PushAliYunParam, "app密钥不合法", nil))
    }
}

func (c *configPush) GetAppSecret() string {
    return c.appSecret
}
