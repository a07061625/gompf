/**
 * 项目常量
 * User: 姜伟
 * Date: 2019/12/24 0024
 * Time: 15:39
 */
package project

import (
    "os"
    "regexp"
    "sync"

    "github.com/spf13/viper"
)

const (
    // http常量
    HttpHeadKeyCookie        = "Set-Cookie"                                       // 请求头名称-cookie
    HttpHeadKeyContentLength = "Content-Length"                                   // 请求头名称-内容长度
    HttpContentTypeForm      = "application/x-www-form-urlencoded; charset=utf-8" // http内容类型-表单
    HttpContentTypeJson      = "application/json; charset=utf-8"                  // http内容类型-json
    HttpContentTypeXml       = "application/xml; charset=utf-8"                   // http内容类型-xml
    HttpContentTypeHtml      = "text/html; charset=utf-8"                         // http内容类型-html
    HttpContentTypeText      = "text/plain; charset=utf-8"                        // http内容类型-text

    // 正则表达式
    RegexWxOpenid   = `^[0-9a-zA-Z\-_]{28}$`                           // 微信openid
    RegexPhone      = `^1[0-9]{10}$`                                   // 手机号码
    RegexDigit      = `^[0-9]+$`                                       // 数字
    RegexAlpha      = `^[a-zA-Z]+$`                                    // 字母
    RegexDigitAlpha = `^[0-9a-zA-Z]+$`                                 // 数字和字母
    RegexDigitLower = `^[0-9a-z]+$`                                    // 数字和小写字母
    RegexDigitUpper = `^[0-9A-Z]+$`                                    // 数字和大写字母
    RegexIp         = `^(\.(\d|[1-9]\d|1\d{2}|2[0-4]\d|25[0-5])){4}$`  // ip
    RegexEmail      = `^\w+([-+.]\w+)*\@\w+([-.]\w+)*\.\w+([-.]\w+)*$` // 邮箱
    RegexUrlHttp    = `^(http|https)://\S+$`                           // http链接

    // 数据
    DataPrefixXml     = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>" // 前缀-xml
    DataParamKeyReqId = "_mp_req_id"                                 // 请求参数键名-请求ID
    DataParamKeyUrl   = "_mp_url"                                    // 请求参数键名-请求url

    // redis前缀常量内部键名
    RedisPrefixCommonUniqueId = "common_unique_id" // 公用-唯一ID

    RedisPrefixWxAccount                 = "wx_account"                    // 微信-公众号
    RedisPrefixWxOpenAccount             = "wx_open_account"               // 微信开放平台-账号
    RedisPrefixWxOpenAuthorize           = "wx_open_authorize"             // 微信开放平台-授权公众号
    RedisPrefixWxOpenAuthorizeCodeSecret = "wx_open_authorize_code_secret" // 微信开放平台-授权小程序代码保护密钥
    RedisPrefixWxCorp                    = "wx_corp"                       // 企业微信-账号
    RedisPrefixWxProviderAccount         = "wx_provider_account"           // 企业微信服务商-账号
    RedisPrefixWxProviderSuite           = "wx_provider_suite"             // 企业微信服务商-套件
    RedisPrefixWxProviderAuthorize       = "wx_provider_authorize"         // 企业微信服务商-授权企业

    RedisPrefixAliPayAccount = "alipay_account" // 支付宝-商户号

    RedisPrefixDingTalkCorp              = "dingtalk_corp"               // 钉钉-企业
    RedisPrefixDingTalkProviderAccount   = "dingtalk_provider_account"   // 钉钉服务商-账号
    RedisPrefixDingTalkProviderSuite     = "dingtalk_provider_suite"     // 钉钉服务商-套件
    RedisPrefixDingTalkProviderAuthorize = "dingtalk_provider_authorize" // 钉钉服务商-授权企业

    RedisPrefixJpushUidPush     = "jpush_uid_push"     // 极光推送唯一标识符-推送
    RedisPrefixJpushUidSchedule = "jpush_uid_schedule" // 极光推送唯一标识符-定时任务

    RedisPrefixPrintFeiYinAccount = "print_feiyin_account" // 飞印打印-账号

    RedisPrefixMQRedis = "mq_redis" // 消息队列-Redis

    // 微信开放平台常量
    WxOpenAuthorizeStatusCancel             = 0 // 授权公众号状态-取消授权
    WxOpenAuthorizeStatusAllow              = 1 // 授权公众号状态-允许授权
    WxOpenAuthorizeOperateAuthorized        = 1 // 授权公众号操作类型-允许授权
    WxOpenAuthorizeOperateUnauthorized      = 2 // 授权公众号操作类型-取消授权
    WxOpenAuthorizeOperateAuthorizedRefresh = 3 // 授权公众号操作类型-更新授权

    // 企业微信服务商常量
    WxProviderAuthorizeStatusCancel       = 0 // 企业微信状态-取消授权
    WxProviderAuthorizeStatusAllow        = 1 // 企业微信状态-允许授权
    WxProviderAuthorizeOperateAuthCreate  = 1 // 企业微信操作类型-成功授权
    WxProviderAuthorizeOperateAuthCancel  = 2 // 企业微信操作类型-取消授权
    WxProviderAuthorizeOperateAuthRefresh = 3 // 企业微信操作类型-更新授权

    // 微信配置常量
    WxConfigStatusInvalid        = 0           // 状态-无效
    WxConfigStatusValid          = 0           // 状态-有效
    WxConfigAuthorizeStatusEmpty = -1          // 第三方授权状态-不存在
    WxConfigAuthorizeStatusNo    = 0           // 第三方授权状态-未授权
    WxConfigAuthorizeStatusYes   = 1           // 第三方授权状态-已授权
    WxConfigCorpStatusInvalid    = 0           // 企业微信状态-无效
    WxConfigCorpStatusValid      = 1           // 企业微信状态-有效
    WxConfigDefaultClientIp      = "127.0.0.1" // 默认客户端IP

    // 支付宝支付常量
    AliPayStatusInvalid = 0 // 状态-无效
    AliPayStatusValid   = 1 // 状态-有效

    // 钉钉配置常量
    DingTalkConfigCorpStatusInvalid = 0 // 企业钉钉状态-无效
    DingTalkConfigCorpStatusValid   = 1 // 企业钉钉状态-有效
)

type project struct {
    // 默认常量
    defaultPage  uint // 页数
    defaultLimit uint // 分页限制

    // 时间常量
    timeClearLocalWxAccount     int64 // 清理本地微信公众号,单位为秒
    timeClearLocalWxCorp        int64 // 清理本地微信企业号,单位为秒
    timeClearLocalAliPayAccount int64 // 清理本地支付宝账号,单位为秒
    timeClearLocalDingTalkCorp  int64 // 清理本地钉钉企业号,单位为秒
    timeClearLocalJPushApp      int64 // 清理本地极光app,单位为秒
    timeClearLocalJPushGroup    int64 // 清理本地极光分组,单位为秒

    // redis前缀常量,6位长度由数字和小写字母组成的字符串,纯数字的是框架内部用,带字母的为项目用
    redisPrefix map[string]string
}

func DefaultPage() uint {
    return ins.defaultPage
}

func DefaultLimit() uint {
    return ins.defaultLimit
}

func TimeClearLocalWxAccount() int64 {
    return ins.timeClearLocalWxAccount
}

func TimeClearLocalWxCorp() int64 {
    return ins.timeClearLocalWxCorp
}

func TimeClearLocalAliPayAccount() int64 {
    return ins.timeClearLocalAliPayAccount
}

func TimeClearLocalDingTalkCorp() int64 {
    return ins.timeClearLocalDingTalkCorp
}

func TimeClearLocalJPushApp() int64 {
    return ins.timeClearLocalJPushApp
}

func TimeClearLocalJPushGroup() int64 {
    return ins.timeClearLocalJPushGroup
}

func RedisPrefix(key string) string {
    val, ok := ins.redisPrefix[key]
    if ok {
        return val
    } else {
        return ""
    }
}

var (
    once sync.Once
    ins  *project
)

func init() {
    ins = &project{}
    ins.redisPrefix = make(map[string]string)

    // 初始化redis缓存前缀
    redisKeyPrefix := "mp" + os.Getenv("MP_PROJECT_TAG")
    ins.redisPrefix[RedisPrefixCommonUniqueId] = redisKeyPrefix + "000000_"
    ins.redisPrefix[RedisPrefixWxAccount] = redisKeyPrefix + "100000_"
    ins.redisPrefix[RedisPrefixWxOpenAccount] = redisKeyPrefix + "100100_"
    ins.redisPrefix[RedisPrefixWxOpenAuthorize] = redisKeyPrefix + "100101_"
    ins.redisPrefix[RedisPrefixWxOpenAuthorizeCodeSecret] = redisKeyPrefix + "100102_"
    ins.redisPrefix[RedisPrefixWxCorp] = redisKeyPrefix + "100200_"
    ins.redisPrefix[RedisPrefixWxProviderAccount] = redisKeyPrefix + "100300_"
    ins.redisPrefix[RedisPrefixWxProviderSuite] = redisKeyPrefix + "100301_"
    ins.redisPrefix[RedisPrefixWxProviderAuthorize] = redisKeyPrefix + "100302_"
    ins.redisPrefix[RedisPrefixAliPayAccount] = redisKeyPrefix + "110000_"
    ins.redisPrefix[RedisPrefixDingTalkCorp] = redisKeyPrefix + "120000_"
    ins.redisPrefix[RedisPrefixDingTalkProviderAccount] = redisKeyPrefix + "120100_"
    ins.redisPrefix[RedisPrefixDingTalkProviderSuite] = redisKeyPrefix + "120101_"
    ins.redisPrefix[RedisPrefixDingTalkProviderAuthorize] = redisKeyPrefix + "120102_"
    ins.redisPrefix[RedisPrefixJpushUidPush] = redisKeyPrefix + "130000_"
    ins.redisPrefix[RedisPrefixJpushUidSchedule] = redisKeyPrefix + "130001_"
    ins.redisPrefix[RedisPrefixPrintFeiYinAccount] = redisKeyPrefix + "130100_"
    ins.redisPrefix[RedisPrefixMQRedis] = redisKeyPrefix + "130200_"
}

func LoadProject(conf *viper.Viper) {
    once.Do(func() {
        projectKey := os.Getenv("MP_PROJECT_KEY")
        ins.defaultPage = conf.GetUint(projectKey + ".defaultpage")
        ins.defaultLimit = conf.GetUint(projectKey + ".defaultlimit")
        ins.timeClearLocalWxAccount = int64(conf.GetUint64(projectKey + ".time.clearlocal.wxaccount"))
        ins.timeClearLocalWxCorp = int64(conf.GetUint64(projectKey + ".time.clearlocal.wxcorp"))
        ins.timeClearLocalAliPayAccount = int64(conf.GetUint64(projectKey + ".time.clearlocal.alipayaccount"))
        ins.timeClearLocalDingTalkCorp = int64(conf.GetUint64(projectKey + ".time.clearlocal.dingtalkcorp"))
        ins.timeClearLocalJPushApp = int64(conf.GetUint64(projectKey + ".time.clearlocal.jpushapp"))
        ins.timeClearLocalJPushGroup = int64(conf.GetUint64(projectKey + ".time.clearlocal.jpushgroup"))

        redisKeyPrefix := "mp" + os.Getenv("MP_PROJECT_TAG")
        redisPrefixList := conf.GetStringMapString(projectKey + ".prefix.redis")
        for k, v := range redisPrefixList {
            match, _ := regexp.MatchString(`^[0-9a-z]{6,9}$`, v)
            if !match {
                continue
            }
            if len(v) == 6 {
                ins.redisPrefix[k] = redisKeyPrefix + v + "_"
            } else if len(v) == 9 {
                ins.redisPrefix[k] = "mp" + v + "_"
            }
        }
    })
}
