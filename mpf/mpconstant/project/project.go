// Package project 项目常量
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

    // HTTPHeadKeyCookie 请求头名称-cookie
    HTTPHeadKeyCookie = "Set-Cookie"
    // HTTPHeadKeyContentType 请求头名称-内容类型
    HTTPHeadKeyContentType = "Content-Type"
    // HTTPHeadKeyContentLength 请求头名称-内容长度
    HTTPHeadKeyContentLength = "Content-Length"
    // HTTPContentTypeForm http内容类型-表单
    HTTPContentTypeForm = "application/x-www-form-urlencoded; charset=utf-8"
    // HTTPContentTypeJSON http内容类型-json
    HTTPContentTypeJSON = "application/json; charset=utf-8"
    // HTTPContentTypeXML http内容类型-xml
    HTTPContentTypeXML = "application/xml; charset=utf-8"
    // HTTPContentTypeHTML http内容类型-html
    HTTPContentTypeHTML = "text/html; charset=utf-8"
    // HTTPContentTypeText http内容类型-text
    HTTPContentTypeText = "text/plain; charset=utf-8"

    // 正则表达式

    // RegexWxOpenid 微信openid
    RegexWxOpenid = `^[0-9a-zA-Z\-_]{28}$`
    // RegexPhone 手机号码
    RegexPhone = `^1[0-9]{10}$`
    // RegexDigit 数字
    RegexDigit = `^[0-9]+$`
    // RegexAlpha 字母
    RegexAlpha = `^[a-zA-Z]+$`
    // RegexDigitAlpha 数字和字母
    RegexDigitAlpha = `^[0-9a-zA-Z]+$`
    // RegexDigitLower 数字和小写字母
    RegexDigitLower = `^[0-9a-z]+$`
    // RegexDigitUpper 数字和大写字母
    RegexDigitUpper = `^[0-9A-Z]+$`
    // RegexIP ip
    RegexIP = `^(\.(\d|[1-9]\d|1\d{2}|2[0-4]\d|25[0-5])){4}$`
    // RegexEmail 邮箱
    RegexEmail = `^\w+([-+.]\w+)*\@\w+([-.]\w+)*\.\w+([-.]\w+)*$`
    // RegexURLHTTP http链接
    RegexURLHTTP = `^(http|https)://\S+$`

    // 数据

    // DataPrefixXML 前缀-xml
    DataPrefixXML = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>"
    // DataParamKeyReqID 参数键名-请求ID
    DataParamKeyReqID = "_mp_req_id"
    // DataParamKeyReqURL 参数键名-请求url
    DataParamKeyReqURL = "_mp_req_url"
    // DataParamKeyRespData 参数键名-响应数据
    DataParamKeyRespData = "_mp_resp_data"
    // DataParamKeyRespType 参数键名-响应类型
    DataParamKeyRespType = "_mp_resp_type"

    // redis前缀常量内部键名

    // RedisPrefixCommonUniqueID 公用-唯一ID
    RedisPrefixCommonUniqueID = "common_unique_id"
    // RedisPrefixWxAccount 微信-公众号
    RedisPrefixWxAccount = "wx_account"
    // RedisPrefixWxOpenAccount 微信开放平台-账号
    RedisPrefixWxOpenAccount = "wx_open_account"
    // RedisPrefixWxOpenAuthorize 微信开放平台-授权公众号
    RedisPrefixWxOpenAuthorize = "wx_open_authorize"
    // RedisPrefixWxOpenAuthorizeCodeSecret 微信开放平台-授权小程序代码保护密钥
    RedisPrefixWxOpenAuthorizeCodeSecret = "wx_open_authorize_code_secret"
    // RedisPrefixWxCorp 企业微信-账号
    RedisPrefixWxCorp = "wx_corp"
    // RedisPrefixWxProviderAccount 企业微信服务商-账号
    RedisPrefixWxProviderAccount = "wx_provider_account"
    // RedisPrefixWxProviderSuite 企业微信服务商-套件
    RedisPrefixWxProviderSuite = "wx_provider_suite"
    // RedisPrefixWxProviderAuthorize 企业微信服务商-授权企业
    RedisPrefixWxProviderAuthorize = "wx_provider_authorize"
    // RedisPrefixAliPayAccount 支付宝-商户号
    RedisPrefixAliPayAccount = "alipay_account"
    // RedisPrefixDingTalkCorp 钉钉-企业
    RedisPrefixDingTalkCorp = "dingtalk_corp"
    // RedisPrefixDingTalkProviderAccount 钉钉服务商-账号
    RedisPrefixDingTalkProviderAccount = "dingtalk_provider_account"
    // RedisPrefixDingTalkProviderSuite 钉钉服务商-套件
    RedisPrefixDingTalkProviderSuite = "dingtalk_provider_suite"
    // RedisPrefixDingTalkProviderAuthorize 钉钉服务商-授权企业
    RedisPrefixDingTalkProviderAuthorize = "dingtalk_provider_authorize"
    // RedisPrefixJpushUIDPush 极光推送唯一标识符-推送
    RedisPrefixJpushUIDPush = "jpush_uid_push"
    // RedisPrefixJpushUIDSchedule 极光推送唯一标识符-定时任务
    RedisPrefixJpushUIDSchedule = "jpush_uid_schedule"
    // RedisPrefixPrintFeiYinAccount 飞印打印-账号
    RedisPrefixPrintFeiYinAccount = "print_feiyin_account"
    // RedisPrefixMQRedis 消息队列-Redis
    RedisPrefixMQRedis = "mq_redis"

    // 微信开放平台常量

    // WxOpenAuthorizeStatusCancel 授权公众号状态-取消授权
    WxOpenAuthorizeStatusCancel = 0
    // WxOpenAuthorizeStatusAllow 授权公众号状态-允许授权
    WxOpenAuthorizeStatusAllow = 1
    // WxOpenAuthorizeOperateAuthorized 授权公众号操作类型-允许授权
    WxOpenAuthorizeOperateAuthorized = 1
    // WxOpenAuthorizeOperateUnauthorized 授权公众号操作类型-取消授权
    WxOpenAuthorizeOperateUnauthorized = 2
    // WxOpenAuthorizeOperateAuthorizedRefresh 授权公众号操作类型-更新授权
    WxOpenAuthorizeOperateAuthorizedRefresh = 3

    // 企业微信服务商常量

    // WxProviderAuthorizeStatusCancel 企业微信状态-取消授权
    WxProviderAuthorizeStatusCancel = 0
    // WxProviderAuthorizeStatusAllow 企业微信状态-允许授权
    WxProviderAuthorizeStatusAllow = 1
    // WxProviderAuthorizeOperateAuthCreate 企业微信操作类型-成功授权
    WxProviderAuthorizeOperateAuthCreate = 1
    // WxProviderAuthorizeOperateAuthCancel 企业微信操作类型-取消授权
    WxProviderAuthorizeOperateAuthCancel = 2
    // WxProviderAuthorizeOperateAuthRefresh 企业微信操作类型-更新授权
    WxProviderAuthorizeOperateAuthRefresh = 3

    // 微信配置常量

    // WxConfigStatusInvalid 状态-无效
    WxConfigStatusInvalid = 0
    // WxConfigStatusValid 状态-有效
    WxConfigStatusValid = 0
    // WxConfigAuthorizeStatusEmpty 第三方授权状态-不存在
    WxConfigAuthorizeStatusEmpty = -1
    // WxConfigAuthorizeStatusNo 第三方授权状态-未授权
    WxConfigAuthorizeStatusNo = 0
    // WxConfigAuthorizeStatusYes 第三方授权状态-已授权
    WxConfigAuthorizeStatusYes = 1
    // WxConfigCorpStatusInvalid 企业微信状态-无效
    WxConfigCorpStatusInvalid = 0
    // WxConfigCorpStatusValid 企业微信状态-有效
    WxConfigCorpStatusValid = 1
    // WxConfigDefaultClientIP 默认客户端IP
    WxConfigDefaultClientIP = "127.0.0.1"

    // 支付宝支付常量

    // AliPayStatusInvalid 状态-无效
    AliPayStatusInvalid = 0
    // AliPayStatusValid  状态-有效
    AliPayStatusValid = 1

    // 钉钉配置常量

    // DingTalkConfigCorpStatusInvalid 企业钉钉状态-无效
    DingTalkConfigCorpStatusInvalid = 0
    // DingTalkConfigCorpStatusValid 企业钉钉状态-有效
    DingTalkConfigCorpStatusValid = 1
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

// DefaultPage 默认页数
func DefaultPage() uint {
    return ins.defaultPage
}

// DefaultLimit 默认每页条数
func DefaultLimit() uint {
    return ins.defaultLimit
}

// TimeClearLocalWxAccount 清理本地微信帐号配置时间
func TimeClearLocalWxAccount() int64 {
    return ins.timeClearLocalWxAccount
}

// TimeClearLocalWxCorp 清理本地微信企业号配置时间
func TimeClearLocalWxCorp() int64 {
    return ins.timeClearLocalWxCorp
}

// TimeClearLocalAliPayAccount 清理本地支付宝帐号配置时间
func TimeClearLocalAliPayAccount() int64 {
    return ins.timeClearLocalAliPayAccount
}

// TimeClearLocalDingTalkCorp 清理本地钉钉企业号配置时间
func TimeClearLocalDingTalkCorp() int64 {
    return ins.timeClearLocalDingTalkCorp
}

// TimeClearLocalJPushApp 清理本地极光推送app配置时间
func TimeClearLocalJPushApp() int64 {
    return ins.timeClearLocalJPushApp
}

// TimeClearLocalJPushGroup 清理本地极光推送分组配置时间
func TimeClearLocalJPushGroup() int64 {
    return ins.timeClearLocalJPushGroup
}

// RedisPrefix redis前缀
func RedisPrefix(key string) string {
    val, ok := ins.redisPrefix[key]
    if ok {
        return val
    }
    return ""
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
    ins.redisPrefix[RedisPrefixCommonUniqueID] = redisKeyPrefix + "000000_"
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
    ins.redisPrefix[RedisPrefixJpushUIDPush] = redisKeyPrefix + "130000_"
    ins.redisPrefix[RedisPrefixJpushUIDSchedule] = redisKeyPrefix + "130001_"
    ins.redisPrefix[RedisPrefixPrintFeiYinAccount] = redisKeyPrefix + "130100_"
    ins.redisPrefix[RedisPrefixMQRedis] = redisKeyPrefix + "130200_"
}

// LoadProject 加载配置
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
