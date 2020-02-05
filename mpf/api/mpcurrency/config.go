/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/17 0017
 * Time: 16:19
 */
package mpcurrency

import (
    "regexp"
    "sync"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
)

type configAMJiSu struct {
    appKey          string // APP KEY
    appSecret       string // APP 密钥
    appCode         string // APP 编码
    serviceProtocol string // 服务协议
    serviceDomain   string // 服务域名
    serviceAddress  string // 服务地址
}

func (c *configAMJiSu) SetAppKey(appKey string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, appKey)
    if match {
        c.appKey = appKey
    } else {
        panic(mperr.NewLogisticsAMAli(errorcode.LogisticsAMAliParam, "app key不合法", nil))
    }
}

func (c *configAMJiSu) GetAppKey() string {
    return c.appKey
}

func (c *configAMJiSu) SetAppSecret(appSecret string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, appSecret)
    if match {
        c.appSecret = appSecret
    } else {
        panic(mperr.NewLogisticsAMAli(errorcode.LogisticsAMAliParam, "app 密钥不合法", nil))
    }
}

func (c *configAMJiSu) GetAppSecret() string {
    return c.appSecret
}

func (c *configAMJiSu) SetAppCode(appCode string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, appCode)
    if match {
        c.appCode = appCode
    } else {
        panic(mperr.NewLogisticsAMAli(errorcode.LogisticsAMAliParam, "app 编码不合法", nil))
    }
}

func (c *configAMJiSu) GetAppCode() string {
    return c.appCode
}

func (c *configAMJiSu) SetServiceAddress(protocol, domain string) {
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

func (c *configAMJiSu) GetServiceProtocol() string {
    return c.serviceProtocol
}

func (c *configAMJiSu) GetServiceDomain() string {
    return c.serviceDomain
}

func (c *configAMJiSu) GetServiceAddress() string {
    return c.serviceAddress
}

type configAMYiYuan struct {
    appKey          string // APP KEY
    appSecret       string // APP 密钥
    appCode         string // APP 编码
    serviceProtocol string // 服务协议
    serviceDomain   string // 服务域名
    serviceAddress  string // 服务地址
}

func (c *configAMYiYuan) SetAppKey(appKey string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, appKey)
    if match {
        c.appKey = appKey
    } else {
        panic(mperr.NewLogisticsAMAli(errorcode.LogisticsAMAliParam, "app key不合法", nil))
    }
}

func (c *configAMYiYuan) GetAppKey() string {
    return c.appKey
}

func (c *configAMYiYuan) SetAppSecret(appSecret string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, appSecret)
    if match {
        c.appSecret = appSecret
    } else {
        panic(mperr.NewLogisticsAMAli(errorcode.LogisticsAMAliParam, "app 密钥不合法", nil))
    }
}

func (c *configAMYiYuan) GetAppSecret() string {
    return c.appSecret
}

func (c *configAMYiYuan) SetAppCode(appCode string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, appCode)
    if match {
        c.appCode = appCode
    } else {
        panic(mperr.NewLogisticsAMAli(errorcode.LogisticsAMAliParam, "app 编码不合法", nil))
    }
}

func (c *configAMYiYuan) GetAppCode() string {
    return c.appCode
}

func (c *configAMYiYuan) SetServiceAddress(protocol, domain string) {
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

func (c *configAMYiYuan) GetServiceProtocol() string {
    return c.serviceProtocol
}

func (c *configAMYiYuan) GetServiceDomain() string {
    return c.serviceDomain
}

func (c *configAMYiYuan) GetServiceAddress() string {
    return c.serviceAddress
}

var (
    onceConfigAMJiSu   sync.Once
    onceConfigAMYiYuan sync.Once
    insConfigAMJiSu    *configAMJiSu
    insConfigAMYiYuan  *configAMYiYuan
)

func init() {
    insConfigAMJiSu = &configAMJiSu{}
    insConfigAMYiYuan = &configAMYiYuan{}
}

func NewConfigAMJiSu() *configAMJiSu {
    onceConfigAMJiSu.Do(func() {
        conf := mpf.NewConfig().GetConfig("mpcurrency")
        protocol := conf.GetString("amjisu." + mpf.EnvProjectKey() + ".service.protocol")
        domain := conf.GetString("amjisu." + mpf.EnvProjectKey() + ".service.domain")
        insConfigAMJiSu.SetAppKey(conf.GetString("amjisu." + mpf.EnvProjectKey() + ".app.key"))
        insConfigAMJiSu.SetAppSecret(conf.GetString("amjisu." + mpf.EnvProjectKey() + ".app.secret"))
        insConfigAMJiSu.SetAppCode(conf.GetString("amjisu." + mpf.EnvProjectKey() + ".app.code"))
        insConfigAMJiSu.SetServiceAddress(protocol, domain)
    })

    return insConfigAMJiSu
}

func NewConfigAMYiYuan() *configAMYiYuan {
    onceConfigAMYiYuan.Do(func() {
        conf := mpf.NewConfig().GetConfig("mpcurrency")
        protocol := conf.GetString("amyiyuan." + mpf.EnvProjectKey() + ".service.protocol")
        domain := conf.GetString("amyiyuan." + mpf.EnvProjectKey() + ".service.domain")
        insConfigAMYiYuan.SetAppKey(conf.GetString("amyiyuan." + mpf.EnvProjectKey() + ".app.key"))
        insConfigAMYiYuan.SetAppSecret(conf.GetString("amyiyuan." + mpf.EnvProjectKey() + ".app.secret"))
        insConfigAMYiYuan.SetAppCode(conf.GetString("amyiyuan." + mpf.EnvProjectKey() + ".app.code"))
        insConfigAMYiYuan.SetServiceAddress(protocol, domain)
    })

    return insConfigAMYiYuan
}
