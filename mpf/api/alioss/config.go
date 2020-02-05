/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/22 0022
 * Time: 17:51
 */
package alioss

import (
    "regexp"
    "sync"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
)

type configAliOss struct {
    endpoint         string // 终端节点
    endpointDomain   string // 终端节点域名
    endpointProtocol string // 终端节点协议
    accessKeyId      string // 帐号ID
    accessKeySecret  string // 帐号密钥
    bucketName       string // 桶名称
    bucketDomain     string // 桶域名
}

func (c *configAliOss) SetEndpoint(protocol, domain string) {
    if (protocol != "http") && (protocol != "https") {
        panic(mperr.NewAliOss(errorcode.AliOssParam, "终端节点协议不合法", nil))
    } else if len(domain) == 0 {
        panic(mperr.NewAliOss(errorcode.AliOssParam, "终端节点域名不合法", nil))
    }

    c.endpointProtocol = protocol
    c.endpointDomain = domain
    c.endpoint = protocol + "://" + domain
}

func (c *configAliOss) GetEndpoint() string {
    return c.endpoint
}

func (c *configAliOss) GetEndpointDomain() string {
    return c.endpointDomain
}

func (c *configAliOss) GetEndpointProtocol() string {
    return c.endpointProtocol
}

func (c *configAliOss) SetAccessKeyId(accessKeyId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, accessKeyId)
    if match {
        c.accessKeyId = accessKeyId
    } else {
        panic(mperr.NewAliOss(errorcode.AliOssParam, "帐号ID不合法", nil))
    }
}

func (c *configAliOss) GetAccessKeyId() string {
    return c.accessKeyId
}

func (c *configAliOss) SetAccessKeySecret(accessKeySecret string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, accessKeySecret)
    if match {
        c.accessKeySecret = accessKeySecret
    } else {
        panic(mperr.NewAliOss(errorcode.AliOssParam, "帐号密钥不合法", nil))
    }
}

func (c *configAliOss) GetAccessKeySecret() string {
    return c.accessKeySecret
}

func (c *configAliOss) SetBucketName(bucketName string) {
    if len(bucketName) > 0 {
        c.bucketName = bucketName
    } else {
        panic(mperr.NewAliOss(errorcode.AliOssParam, "桶名称不合法", nil))
    }
}

func (c *configAliOss) GetBucketName() string {
    return c.bucketName
}

func (c *configAliOss) SetBucketDomain(bucketDomain string) {
    match, _ := regexp.MatchString(project.RegexUrlHttp, bucketDomain)
    if match {
        c.bucketDomain = bucketDomain
    } else {
        panic(mperr.NewAliOss(errorcode.AliOssParam, "桶域名不合法", nil))
    }
}

func (c *configAliOss) GetBucketDomain() string {
    return c.bucketDomain
}

var (
    onceConfig sync.Once
    insConfig  *configAliOss
)

func init() {
    insConfig = &configAliOss{}
}

func NewConfig() *configAliOss {
    onceConfig.Do(func() {
        conf := mpf.NewConfig().GetConfig("alioss")
        protocol := conf.GetString(mpf.EnvProjectKey() + ".endpoint.protocol")
        domain := conf.GetString(mpf.EnvProjectKey() + ".endpoint.domain")
        insConfig.SetEndpoint(protocol, domain)
        insConfig.SetAccessKeyId(conf.GetString(mpf.EnvProjectKey() + ".access.key.id"))
        insConfig.SetAccessKeySecret(conf.GetString(mpf.EnvProjectKey() + ".access.key.secret"))
        insConfig.SetBucketName(conf.GetString(mpf.EnvProjectKey() + ".bucket.name"))
        insConfig.SetBucketDomain(conf.GetString(mpf.EnvProjectKey() + ".bucket.domain"))
    })

    return insConfig
}
