/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/20 0020
 * Time: 13:24
 */
package qcloud

import (
    "regexp"
    "sync"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
)

type configCos struct {
    appId        string // 应用ID
    secretId     string // 帐户ID
    secretKey    string // 帐户密钥
    bucketName   string // 桶名称
    bucketDomain string // 桶用户域名
    regionTag    string // 地域标识
    bucketHost   string // 桶域名
    controlHost  string // 控制台域名
    hostFlag     bool   // 请求标识 true:未生成域名 false:已生成域名
}

func (c *configCos) SetAppId(appId string) {
    match, _ := regexp.MatchString(project.RegexDigit, appId)
    if match {
        c.appId = appId
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "应用ID不合法", nil))
    }
}

func (c *configCos) GetAppId() string {
    return c.appId
}

func (c *configCos) SetSecretId(secretId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, secretId)
    if match {
        c.secretId = secretId
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "帐户ID不合法", nil))
    }
}

func (c *configCos) GetSecretId() string {
    return c.secretId
}

func (c *configCos) SetSecretKey(secretKey string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, secretKey)
    if match {
        c.secretKey = secretKey
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "帐户密钥不合法", nil))
    }
}

func (c *configCos) GetSecretKey() string {
    return c.secretKey
}

func (c *configCos) SetBucketName(bucketName string) {
    if len(bucketName) > 0 {
        c.bucketName = bucketName
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "桶名称不合法", nil))
    }
}

func (c *configCos) GetBucketName() string {
    return c.bucketName
}

func (c *configCos) SetBucketDomain(bucketDomain string) {
    match, _ := regexp.MatchString(project.RegexUrlHttp, bucketDomain)
    if match {
        c.bucketDomain = bucketDomain
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "桶用户域名不合法", nil))
    }
}

func (c *configCos) GetBucketDomain() string {
    return c.bucketDomain
}

func (c *configCos) SetRegionTag(regionTag string) {
    if len(regionTag) > 0 {
        c.regionTag = regionTag
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "地域标识不合法", nil))
    }
}

func (c *configCos) GetRegionTag() string {
    return c.regionTag
}

func (c *configCos) CreateHost() {
    if c.hostFlag {
        c.bucketHost = c.bucketName + "-" + c.appId + ".cos." + c.regionTag + ".myqcloud.com"
        c.controlHost = c.appId + ".cos-control." + c.regionTag + ".myqcloud.com"
        c.hostFlag = false
    }
}

func (c *configCos) GetBucketHost() string {
    return c.bucketHost
}

func (c *configCos) GetControlHost() string {
    return c.controlHost
}

func newConfigCos() *configCos {
    return &configCos{"", "", "", "", "", "", "", "", true}
}

type configQCloud struct {
    cos *configCos
}

func (c *configQCloud) GetCos() *configCos {
    onceConfigCos.Do(func() {
        conf := mpf.NewConfig().GetConfig("qcloud")
        c.cos.SetAppId(conf.GetString("cos." + mpf.EnvProjectKey() + ".app.id"))
        c.cos.SetSecretId(conf.GetString("cos." + mpf.EnvProjectKey() + ".secret.id"))
        c.cos.SetSecretKey(conf.GetString("cos." + mpf.EnvProjectKey() + ".secret.key"))
        c.cos.SetBucketName(conf.GetString("cos." + mpf.EnvProjectKey() + ".bucket.name"))
        c.cos.SetBucketDomain(conf.GetString("cos." + mpf.EnvProjectKey() + ".bucket.domain"))
        c.cos.SetRegionTag(conf.GetString("cos." + mpf.EnvProjectKey() + ".region.tag"))
        c.cos.CreateHost()
    })

    return c.cos
}

var (
    onceConfigCos sync.Once
    insConfig     *configQCloud
)

func init() {
    insConfig = &configQCloud{}
    insConfig.cos = newConfigCos()
}

func NewConfig() *configQCloud {
    return insConfig
}
