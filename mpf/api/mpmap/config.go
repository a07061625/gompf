/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2019/12/25 0025
 * Time: 11:42
 */
package mpmap

import (
    "regexp"
    "sync"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
)

type configBaiDu struct {
    ak       string // 开发密钥
    serverIp string // 服务器IP
}

func (c *configBaiDu) SetAk(ak string) {
    match, _ := regexp.MatchString(`^[0-9a-z]{32}$`, ak)
    if match {
        c.ak = ak
    } else {
        panic(mperr.NewMapBaiDu(errorcode.MapBaiDuParam, "密钥不合法", nil))
    }
}

func (c *configBaiDu) GetAk() string {
    return c.ak
}

func (c *configBaiDu) SetServerIp(serverIp string) {
    match, _ := regexp.MatchString(project.RegexIP, "."+serverIp)
    if match {
        c.serverIp = serverIp
    } else {
        panic(mperr.NewMapBaiDu(errorcode.MapBaiDuParam, "服务器IP不合法", nil))
    }
}

func (c *configBaiDu) GetServerIp() string {
    return c.serverIp
}

type configGaoDe struct {
    key    string // 应用KEY
    secret string // 应用密钥
}

func (c *configGaoDe) SetKey(key string) {
    match, _ := regexp.MatchString(project.RegexDigitLower, key)
    if match {
        c.key = key
    } else {
        panic(mperr.NewMapGaoDe(errorcode.MapGaoDeParam, "应用KEY不合法", nil))
    }
}

func (c *configGaoDe) GetKey() string {
    return c.key
}

func (c *configGaoDe) SetSecret(secret string) {
    match, _ := regexp.MatchString(project.RegexDigitLower, secret)
    if match {
        c.secret = secret
    } else {
        panic(mperr.NewMapGaoDe(errorcode.MapGaoDeParam, "应用密钥不合法", nil))
    }
}

func (c *configGaoDe) GetSecret() string {
    return c.secret
}

type configTencent struct {
    key      string // 开发密钥
    serverIp string // 服务器IP
    domain   string // 域名
}

func (c *configTencent) SetKey(key string) {
    match, _ := regexp.MatchString(`^(\-[0-9A-Z]{5}){6}$`, "-"+key)
    if match {
        c.key = key
    } else {
        panic(mperr.NewMapTencent(errorcode.MapTencentParam, "开发密钥不合法", nil))
    }
}

func (c *configTencent) GetKey() string {
    return c.key
}

func (c *configTencent) SetServerIp(serverIp string) {
    match, _ := regexp.MatchString(project.RegexIP, "."+serverIp)
    if match {
        c.serverIp = serverIp
    } else {
        panic(mperr.NewMapTencent(errorcode.MapTencentParam, "服务器IP不合法", nil))
    }
}

func (c *configTencent) GetServerIp() string {
    return c.serverIp
}

func (c *configTencent) SetDomain(domain string) {
    match, _ := regexp.MatchString(project.RegexURLHTTP, domain)
    if match {
        c.domain = domain
    } else {
        panic(mperr.NewMapTencent(errorcode.MapTencentParam, "域名不合法", nil))
    }
}

func (c *configTencent) GetDomain() string {
    return c.domain
}

var (
    onceConfigBaidu   sync.Once
    onceConfigGaoDe   sync.Once
    onceConfigTencent sync.Once
    insConfigBaidu    *configBaiDu
    insConfigGaoDe    *configGaoDe
    insConfigTencent  *configTencent
)

func init() {
    insConfigBaidu = &configBaiDu{"", ""}
    insConfigGaoDe = &configGaoDe{"", ""}
    insConfigTencent = &configTencent{"", "", ""}
}

func NewConfigBaiDu() *configBaiDu {
    onceConfigBaidu.Do(func() {
        conf := mpf.NewConfig().GetConfig("mpmap")
        insConfigBaidu.SetAk(conf.GetString("baidu." + mpf.EnvProjectKey() + ".ak"))
        insConfigBaidu.SetServerIp(conf.GetString("baidu." + mpf.EnvProjectKey() + ".serverip"))
    })

    return insConfigBaidu
}

func NewConfigGaoDe() *configGaoDe {
    onceConfigGaoDe.Do(func() {
        conf := mpf.NewConfig().GetConfig("mpmap")
        insConfigGaoDe.SetKey(conf.GetString("gaode." + mpf.EnvProjectKey() + ".key"))
        insConfigGaoDe.SetSecret(conf.GetString("gaode." + mpf.EnvProjectKey() + ".secret"))
    })

    return insConfigGaoDe
}

func NewConfigTencent() *configTencent {
    onceConfigTencent.Do(func() {
        conf := mpf.NewConfig().GetConfig("mpmap")
        insConfigTencent.SetKey(conf.GetString("tencent." + mpf.EnvProjectKey() + ".key"))
        insConfigTencent.SetServerIp(conf.GetString("tencent." + mpf.EnvProjectKey() + ".serverip"))
        insConfigTencent.SetDomain(conf.GetString("tencent." + mpf.EnvProjectKey() + ".domain"))
    })

    return insConfigTencent
}
