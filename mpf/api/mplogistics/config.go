/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/15 0015
 * Time: 14:06
 */
package mplogistics

import (
    "regexp"
    "sync"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
)

type configAMAli struct {
    appKey          string // APP KEY
    appSecret       string // APP 密钥
    appCode         string // APP 编码
    serviceProtocol string // 服务协议
    serviceDomain   string // 服务域名
    serviceAddress  string // 服务地址
}

func (c *configAMAli) SetAppKey(appKey string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, appKey)
    if match {
        c.appKey = appKey
    } else {
        panic(mperr.NewLogisticsAMAli(errorcode.LogisticsAMAliParam, "app key不合法", nil))
    }
}

func (c *configAMAli) GetAppKey() string {
    return c.appKey
}

func (c *configAMAli) SetAppSecret(appSecret string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, appSecret)
    if match {
        c.appSecret = appSecret
    } else {
        panic(mperr.NewLogisticsAMAli(errorcode.LogisticsAMAliParam, "app 密钥不合法", nil))
    }
}

func (c *configAMAli) GetAppSecret() string {
    return c.appSecret
}

func (c *configAMAli) SetAppCode(appCode string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, appCode)
    if match {
        c.appCode = appCode
    } else {
        panic(mperr.NewLogisticsAMAli(errorcode.LogisticsAMAliParam, "app 编码不合法", nil))
    }
}

func (c *configAMAli) GetAppCode() string {
    return c.appCode
}

func (c *configAMAli) SetServiceAddress(protocol, domain string) {
    if (protocol != "http") && (protocol != "https") {
        panic(mperr.NewLogisticsAMAli(errorcode.LogisticsAMAliParam, "服务协议不合法", nil))
    }
    if len(domain) == 0 {
        panic(mperr.NewLogisticsAMAli(errorcode.LogisticsAMAliParam, "服务域名不合法", nil))
    }
    c.serviceProtocol = protocol
    c.serviceDomain = domain
    c.serviceAddress = protocol + "://" + domain
}

func (c *configAMAli) GetServiceProtocol() string {
    return c.serviceProtocol
}

func (c *configAMAli) GetServiceDomain() string {
    return c.serviceDomain
}

func (c *configAMAli) GetServiceAddress() string {
    return c.serviceAddress
}

type configKd100 struct {
    appId  string // 应用ID
    appKey string // 应用密钥
}

func (c *configKd100) SetAppId(appId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, appId)
    if match {
        c.appId = appId
    } else {
        panic(mperr.NewLogisticsKd100(errorcode.LogisticsKd100Param, "应用ID不合法", nil))
    }
}

func (c *configKd100) GetAppId() string {
    return c.appId
}

func (c *configKd100) SetAppKey(appKey string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, appKey)
    if match {
        c.appKey = appKey
    } else {
        panic(mperr.NewLogisticsKd100(errorcode.LogisticsKd100Param, "应用密钥不合法", nil))
    }
}

func (c *configKd100) GetAppKey() string {
    return c.appKey
}

type configKdBird struct {
    businessId string // 商户ID
    appKey     string // 应用密钥
}

func (c *configKdBird) SetBusinessId(businessId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, businessId)
    if match {
        c.businessId = businessId
    } else {
        panic(mperr.NewLogisticsKdBird(errorcode.LogisticsKdBirdParam, "商户ID不合法", nil))
    }
}

func (c *configKdBird) GetBusinessId() string {
    return c.businessId
}

func (c *configKdBird) SetAppKey(appKey string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, appKey)
    if match {
        c.appKey = appKey
    } else {
        panic(mperr.NewLogisticsKdBird(errorcode.LogisticsKdBirdParam, "应用密钥不合法", nil))
    }
}

func (c *configKdBird) GetAppKey() string {
    return c.appKey
}

var (
    onceConfigAMAli  sync.Once
    onceConfigKd100  sync.Once
    onceConfigKdBird sync.Once
    insConfigAMAli   *configAMAli
    insConfigKd100   *configKd100
    insConfigKdBird  *configKdBird
)

func init() {
    insConfigAMAli = &configAMAli{}
    insConfigKd100 = &configKd100{}
    insConfigKdBird = &configKdBird{}
}

func NewConfigAMAli() *configAMAli {
    onceConfigAMAli.Do(func() {
        conf := mpf.NewConfig().GetConfig("mplogistics")
        protocol := conf.GetString("amali." + mpf.EnvProjectKey() + ".service.protocol")
        domain := conf.GetString("amali." + mpf.EnvProjectKey() + ".service.domain")
        insConfigAMAli.SetAppKey(conf.GetString("amali." + mpf.EnvProjectKey() + ".app.key"))
        insConfigAMAli.SetAppSecret(conf.GetString("amali." + mpf.EnvProjectKey() + ".app.secret"))
        insConfigAMAli.SetAppCode(conf.GetString("amali." + mpf.EnvProjectKey() + ".app.code"))
        insConfigAMAli.SetServiceAddress(protocol, domain)
    })

    return insConfigAMAli
}

func NewConfigKd100() *configKd100 {
    onceConfigKd100.Do(func() {
        conf := mpf.NewConfig().GetConfig("mplogistics")
        insConfigKd100.SetAppId(conf.GetString("kd100." + mpf.EnvProjectKey() + ".app.id"))
        insConfigKd100.SetAppKey(conf.GetString("kd100." + mpf.EnvProjectKey() + ".app.key"))
    })

    return insConfigKd100
}

func NewConfigKdBird() *configKdBird {
    onceConfigKdBird.Do(func() {
        conf := mpf.NewConfig().GetConfig("mplogistics")
        insConfigKdBird.SetBusinessId(conf.GetString("kdbird." + mpf.EnvProjectKey() + ".business.id"))
        insConfigKdBird.SetAppKey(conf.GetString("kdbird." + mpf.EnvProjectKey() + ".app.key"))
    })

    return insConfigKdBird
}
