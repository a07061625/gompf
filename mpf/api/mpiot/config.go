/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/1 0001
 * Time: 15:15
 */
package mpiot

import (
    "regexp"

    "sync"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
)

type configBaiDu struct {
    accessKey    string // 访问ID
    accessSecret string // 访问密钥
}

func (c *configBaiDu) SetAccessKey(accessKey string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, accessKey)
    if match {
        c.accessKey = accessKey
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "访问ID不合法", nil))
    }
}

func (c *configBaiDu) GetAccessKey() string {
    return c.accessKey
}

func (c *configBaiDu) SetAccessSecret(accessSecret string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, accessSecret)
    if match {
        c.accessSecret = accessSecret
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "访问密钥不合法", nil))
    }
}

func (c *configBaiDu) GetAccessSecret() string {
    return c.accessSecret
}

type configTencent struct {
    regionId  string // 区域ID
    secretId  string // 应用ID
    secretKey string // 应用密钥
}

func (c *configTencent) SetRegionId(regionId string) {
    if len(regionId) > 0 {
        c.regionId = regionId
    } else {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "区域ID不合法", nil))
    }
}

func (c *configTencent) GetRegionId() string {
    return c.regionId
}

func (c *configTencent) SetSecretId(secretId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, secretId)
    if match {
        c.secretId = secretId
    } else {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "应用ID不合法", nil))
    }
}

func (c *configTencent) GetSecretId() string {
    return c.secretId
}

func (c *configTencent) SetSecretKey(secretKey string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, secretKey)
    if match {
        c.secretKey = secretKey
    } else {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "应用密钥不合法", nil))
    }
}

func (c *configTencent) GetSecretKey() string {
    return c.secretKey
}

type configIot struct {
    baidu   *configBaiDu
    tencent *configTencent
}

func (c *configIot) GetBaiDu() *configBaiDu {
    onceConfigBaiDu.Do(func() {
        conf := mpf.NewConfig().GetConfig("mpiot")
        c.baidu.SetAccessKey(conf.GetString("baidu." + mpf.EnvProjectKey() + ".access.key"))
        c.baidu.SetAccessSecret(conf.GetString("baidu." + mpf.EnvProjectKey() + ".access.secret"))
    })

    return c.baidu
}

func (c *configIot) GetTencent() *configTencent {
    onceConfigTencent.Do(func() {
        conf := mpf.NewConfig().GetConfig("mpiot")
        c.tencent.SetRegionId(conf.GetString("tencent." + mpf.EnvProjectKey() + ".region.id"))
        c.tencent.SetSecretId(conf.GetString("tencent." + mpf.EnvProjectKey() + ".secret.id"))
        c.tencent.SetSecretKey(conf.GetString("tencent." + mpf.EnvProjectKey() + ".secret.key"))
    })

    return c.tencent
}

var (
    onceConfigBaiDu   sync.Once
    onceConfigTencent sync.Once
    insConfig         *configIot
)

func init() {
    insConfig = &configIot{}
    insConfig.baidu = &configBaiDu{}
    insConfig.tencent = &configTencent{}
}

func NewConfig() *configIot {
    return insConfig
}
