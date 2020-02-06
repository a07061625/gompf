/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2019/12/25 0025
 * Time: 10:42
 */
package wx

import (
    "regexp"
    "sync"
    "time"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
)

type configSingle struct {
    clientIp     string // 客户端IP
    payMchId     string // 商户号
    payKey       string // 商户支付密钥
    payNotifyUrl string // 支付异步通知URL
    payAuthUrl   string // 支付授权URL
    sslCert      string // CERT PEM证书内容,包含---PUBLIC XXX---内容
    sslKey       string // KEY PEM证书内容,包含---PRIVATE XXX---内容
    valid        bool   // 配置有效状态
    expireTime   int    // 配置过期时间戳
}

func (c *configSingle) SetClientIp(clientIp string) {
    match, _ := regexp.MatchString(project.RegexIp, "."+clientIp)
    if match {
        c.clientIp = clientIp
    } else {
        panic(mperr.NewWx(errorcode.WxParam, "客户端IP不合法", nil))
    }
}

func (c *configSingle) GetClientIp() string {
    return c.clientIp
}

func (c *configSingle) SetPayMchId(payMchId string) {
    match, _ := regexp.MatchString(project.RegexDigit, payMchId)
    if match {
        c.payMchId = payMchId
    } else {
        panic(mperr.NewWx(errorcode.WxParam, "商户号不合法", nil))
    }
}

func (c *configSingle) GetPayMchId() string {
    return c.payMchId
}

func (c *configSingle) SetPayKey(payKey string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{32}$`, payKey)
    if match {
        c.payKey = payKey
    } else {
        panic(mperr.NewWx(errorcode.WxParam, "支付密钥不合法", nil))
    }
}

func (c *configSingle) GetPayKey() string {
    return c.payKey
}

func (c *configSingle) SetPayNotifyUrl(payNotifyUrl string) {
    match, _ := regexp.MatchString(project.RegexUrlHttp, payNotifyUrl)
    if match {
        c.payNotifyUrl = payNotifyUrl
    } else {
        panic(mperr.NewWx(errorcode.WxParam, "支付异步消息通知URL不合法", nil))
    }
}

func (c *configSingle) GetPayNotifyUrl() string {
    return c.payNotifyUrl
}

func (c *configSingle) SetPayAuthUrl(payAuthUrl string) {
    match, _ := regexp.MatchString(project.RegexUrlHttp, payAuthUrl)
    if match {
        c.payAuthUrl = payAuthUrl
    } else {
        panic(mperr.NewWx(errorcode.WxParam, "支付授权URL不合法", nil))
    }
}

func (c *configSingle) GetPayAuthUrl() string {
    return c.payAuthUrl
}

func (c *configSingle) SetSslCert(sslCert string) {
    c.sslCert = sslCert
}

func (c *configSingle) GetSslCert() string {
    return c.sslCert
}

func (c *configSingle) SetSslKey(sslKey string) {
    c.sslKey = sslKey
}

func (c *configSingle) GetSslKey() string {
    return c.sslKey
}

func (c *configSingle) SetValid(valid bool) {
    c.valid = valid
}

func (c *configSingle) IsValid() bool {
    return c.valid
}

func (c *configSingle) SetExpireTime(expireTime int) {
    c.expireTime = expireTime
}

func (c *configSingle) GetExpireTime() int {
    return c.expireTime
}

func newConfigSingle() configSingle {
    return configSingle{"", "", "", "", "", "", "", false, 0}
}

// 公众号配置
type configAccount struct {
    configSingle
    appId          string            // 微信号
    secret         string            // 微信密钥
    sslCompanyBank string            // 企业付款银行卡公钥内容,包含---PUBLIC XXX---内容
    templates      map[string]string // 模板列表
    merchantAppId  string            // 服务商微信号
}

func (c *configAccount) SetAppId(appId string) {
    match, _ := regexp.MatchString(`^[0-9a-z]{18}$`, appId)
    if match {
        c.appId = appId
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "微信号不合法", nil))
    }
}

func (c *configAccount) GetAppId() string {
    return c.appId
}

func (c *configAccount) SetSecret(secret string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{32}$`, secret)
    if match {
        c.secret = secret
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "微信密钥不合法", nil))
    }
}

func (c *configAccount) GetSecret() string {
    return c.secret
}

func (c *configAccount) SetSslCompanyBank(sslCompanyBank string) {
    c.sslCompanyBank = sslCompanyBank
}

func (c *configAccount) GetSslCompanyBank() string {
    return c.sslCompanyBank
}

func (c *configAccount) SetTemplates(templates map[string]string) {
    trueTemplates := make(map[string]string)
    for k, v := range templates {
        if (len(k) > 0) && (len(v) > 0) {
            trueTemplates[k] = v
        }
    }
    c.templates = trueTemplates
}

func (c *configAccount) GetTemplates() map[string]string {
    return c.templates
}

func (c *configAccount) GetTemplateId(tplTag string) string {
    tplId, ok := c.templates[tplTag]
    if ok {
        return tplId
    } else {
        return ""
    }
}

func (c *configAccount) SetMerchantAppId(merchantAppId string) {
    if len(merchantAppId) == 0 {
        c.merchantAppId = ""
    } else {
        match, _ := regexp.MatchString(`^[0-9a-z]{18}$`, merchantAppId)
        if match {
            c.merchantAppId = merchantAppId
        } else {
            panic(mperr.NewWxAccount(errorcode.WxAccountParam, "服务商微信号不合法", nil))
        }
    }
}

func (c *configAccount) GetMerchantAppId() string {
    return c.merchantAppId
}

func NewConfigAccount() *configAccount {
    return &configAccount{newConfigSingle(), "", "", "", make(map[string]string), ""}
}

// 企业号配置
type configCorp struct {
    configSingle
    corpId       string                       // 企业ID
    agents       map[string]map[string]string // 应用列表
    urlAuthLogin string                       // 登录授权地址
}

func (c *configCorp) SetCorpId(corpId string) {
    match, _ := regexp.MatchString(project.RegexDigitLower, corpId)
    if match {
        c.corpId = corpId
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "企业ID不合法", nil))
    }
}

func (c *configCorp) GetCorpId() string {
    return c.corpId
}

func (c *configCorp) SetAgents(agents map[string]map[string]string) {
    c.agents = agents
}

func (c *configCorp) GetAgents() map[string]map[string]string {
    return c.agents
}

func (c *configCorp) GetAgentInfo(agentTag string) map[string]string {
    info, ok := c.agents[agentTag]
    if ok {
        return info
    } else {
        return make(map[string]string)
    }
}

func (c *configCorp) SetUrlAuthLogin(urlAuthLogin string) {
    match, _ := regexp.MatchString(project.RegexUrlHttp, urlAuthLogin)
    if match {
        c.urlAuthLogin = urlAuthLogin
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "登录授权地址不合法", nil))
    }
}

func (c *configCorp) GetUrlAuthLogin() string {
    return c.urlAuthLogin
}

func NewConfigCorp() *configCorp {
    return &configCorp{newConfigSingle(), "", make(map[string]map[string]string), ""}
}

// 第三方开放平台配置
type configOpen struct {
    appId               string   // 微信号
    secret              string   // 微信密钥
    token               string   // 消息校验token
    aesKeyBefore        string   // 旧消息加解密key
    aesKeyNow           string   // 新消息加解密key
    urlAuth             string   // 授权页面域名
    urlAuthCallback     string   // 授权页面回跳地址
    urlMiniRebindAdmin  string   // 换绑小程序管理员回跳地址
    urlMiniFastRegister string   // 快速注册小程序回跳地址
    domainMiniServers   []string // 小程序服务域名列表
    domainMiniWebViews  []string // 小程序业务域名列表
}

func (c *configOpen) SetAppId(appId string) {
    match, _ := regexp.MatchString(`^[0-9a-z]{18}$`, appId)
    if match {
        c.appId = appId
    } else {
        panic(mperr.NewWxOpen(errorcode.WxOpenParam, "微信号不合法", nil))
    }
}

func (c *configOpen) GetAppId() string {
    return c.appId
}

func (c *configOpen) SetSecret(secret string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{32}$`, secret)
    if match {
        c.secret = secret
    } else {
        panic(mperr.NewWxOpen(errorcode.WxOpenParam, "微信密钥不合法", nil))
    }
}

func (c *configOpen) GetSecret() string {
    return c.secret
}

func (c *configOpen) SetToken(token string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{1,32}$`, token)
    if match {
        c.token = token
    } else {
        panic(mperr.NewWxOpen(errorcode.WxOpenParam, "消息校验token不合法", nil))
    }
}

func (c *configOpen) GetToken() string {
    return c.token
}

func (c *configOpen) SetAesKeyBefore(aesKeyBefore string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{43}$`, aesKeyBefore)
    if match {
        c.aesKeyBefore = aesKeyBefore
    } else {
        panic(mperr.NewWxOpen(errorcode.WxOpenParam, "旧消息加解密key不合法", nil))
    }
}

func (c *configOpen) GetAesKeyBefore() string {
    return c.aesKeyBefore
}

func (c *configOpen) SetAesKeyNow(aesKeyNow string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{43}$`, aesKeyNow)
    if match {
        c.aesKeyNow = aesKeyNow
    } else {
        panic(mperr.NewWxOpen(errorcode.WxOpenParam, "新消息加解密key不合法", nil))
    }
}

func (c *configOpen) GetAesKeyNow() string {
    return c.aesKeyNow
}

func (c *configOpen) SetUrlAuth(urlAuth string) {
    match, _ := regexp.MatchString(project.RegexUrlHttp, urlAuth)
    if match {
        c.urlAuth = urlAuth
    } else {
        panic(mperr.NewWxOpen(errorcode.WxOpenParam, "授权页面域名不合法", nil))
    }
}

func (c *configOpen) GetUrlAuth() string {
    return c.urlAuth
}

func (c *configOpen) SetUrlAuthCallback(urlAuthCallback string) {
    match, _ := regexp.MatchString(project.RegexUrlHttp, urlAuthCallback)
    if match {
        c.urlAuthCallback = urlAuthCallback
    } else {
        panic(mperr.NewWxOpen(errorcode.WxOpenParam, "授权页面回跳地址不合法", nil))
    }
}

func (c *configOpen) GetUrlAuthCallback() string {
    return c.urlAuthCallback
}

func (c *configOpen) SetUrlMiniRebindAdmin(urlMiniRebindAdmin string) {
    match, _ := regexp.MatchString(project.RegexUrlHttp, urlMiniRebindAdmin)
    if match {
        c.urlMiniRebindAdmin = urlMiniRebindAdmin
    } else {
        panic(mperr.NewWxOpen(errorcode.WxOpenParam, "换绑小程序管理员回跳地址不合法", nil))
    }
}

func (c *configOpen) GetUrlMiniRebindAdmin() string {
    return c.urlMiniRebindAdmin
}

func (c *configOpen) SetUrlMiniFastRegister(urlMiniFastRegister string) {
    match, _ := regexp.MatchString(project.RegexUrlHttp, urlMiniFastRegister)
    if match {
        c.urlMiniFastRegister = urlMiniFastRegister
    } else {
        panic(mperr.NewWxOpen(errorcode.WxOpenParam, "快速注册小程序回跳地址不合法", nil))
    }
}

func (c *configOpen) GetUrlMiniFastRegister() string {
    return c.urlMiniFastRegister
}

func (c *configOpen) SetDomainMiniServers(domainMiniServers []string) {
    c.domainMiniServers = domainMiniServers
}

func (c *configOpen) GetDomainMiniServers() []string {
    return c.domainMiniServers
}

func (c *configOpen) SetDomainMiniWebViews(domainMiniWebViews []string) {
    c.domainMiniWebViews = domainMiniWebViews
}

func (c *configOpen) GetDomainMiniWebViews() []string {
    return c.domainMiniWebViews
}

// 企业服务号配置
type configProvider struct {
    corpId       string // 企业ID
    corpSecret   string // 企业密钥
    token        string // 消息校验token
    aesKey       string // 消息加解密key
    suiteId      string // 套件ID
    suiteSecret  string // 套件密钥
    urlAuthSuite string // 套件授权地址
    urlAuthLogin string // 登录授权地址
}

func (c *configProvider) SetCorpId(corpId string) {
    match, _ := regexp.MatchString(project.RegexDigitLower, corpId)
    if match {
        c.corpId = corpId
    } else {
        panic(mperr.NewWxProvider(errorcode.WxProviderParam, "企业ID不合法", nil))
    }
}

func (c *configProvider) GetCorpId() string {
    return c.corpId
}

func (c *configProvider) SetCorpSecret(corpSecret string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, corpSecret)
    if match {
        c.corpSecret = corpSecret
    } else {
        panic(mperr.NewWxProvider(errorcode.WxProviderParam, "企业密钥不合法", nil))
    }
}

func (c *configProvider) GetCorpSecret() string {
    return c.corpSecret
}

func (c *configProvider) SetToken(token string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{1,32}$`, token)
    if match {
        c.token = token
    } else {
        panic(mperr.NewWxProvider(errorcode.WxProviderParam, "消息校验token不合法", nil))
    }
}

func (c *configProvider) GetToken() string {
    return c.token
}

func (c *configProvider) SetAesKey(aesKey string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{43}$`, aesKey)
    if match {
        c.aesKey = aesKey
    } else {
        panic(mperr.NewWxProvider(errorcode.WxProviderParam, "消息加解密key不合法", nil))
    }
}

func (c *configProvider) GetAesKey() string {
    return c.aesKey
}

func (c *configProvider) SetSuiteId(suiteId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, suiteId)
    if match {
        c.suiteId = suiteId
    } else {
        panic(mperr.NewWxProvider(errorcode.WxProviderParam, "套件ID不合法", nil))
    }
}

func (c *configProvider) GetSuiteId() string {
    return c.suiteId
}

func (c *configProvider) SetSuiteSecret(suiteSecret string) {
    if len(suiteSecret) > 0 {
        c.suiteSecret = suiteSecret
    } else {
        panic(mperr.NewWxProvider(errorcode.WxProviderParam, "套件密钥不合法", nil))
    }
}

func (c *configProvider) GetSuiteSecret() string {
    return c.suiteSecret
}

func (c *configProvider) SetUrlAuthSuite(urlAuthSuite string) {
    match, _ := regexp.MatchString(project.RegexUrlHttp, urlAuthSuite)
    if match {
        c.urlAuthSuite = urlAuthSuite
    } else {
        panic(mperr.NewWxProvider(errorcode.WxProviderParam, "套件授权地址不合法", nil))
    }
}

func (c *configProvider) GetUrlAuthSuite() string {
    return c.urlAuthSuite
}

func (c *configProvider) SetUrlAuthLogin(urlAuthLogin string) {
    match, _ := regexp.MatchString(project.RegexUrlHttp, urlAuthLogin)
    if match {
        c.urlAuthLogin = urlAuthLogin
    } else {
        panic(mperr.NewWxProvider(errorcode.WxProviderParam, "登录授权地址不合法", nil))
    }
}

func (c *configProvider) GetUrlAuthLogin() string {
    return c.urlAuthLogin
}

type IWxConfig interface {
    RefreshConfigAccount(appId string) *configAccount // 刷新公众号配置
    RefreshConfigCorp(appId string) *configCorp       // 刷新企业号配置
}

type configWx struct {
    outer            IWxConfig                 // 项目实现的接口,用于根据项目获取到相应的配置实例
    open             *configOpen               // 第三方平台配置
    provider         *configProvider           // 企业服务商配置
    accountClearTime int                       // 公众号本地清理时间戳
    accountList      map[string]*configAccount // 公众号本地配置集合
    corpClearTime    int                       // 企业号本地清理时间戳
    corpList         map[string]*configCorp    // 企业号本地配置集合
}

func (c *configWx) GetOpen() *configOpen {
    onceConfigOpen.Do(func() {
        conf := mpf.NewConfig().GetConfig("wx")
        c.open.SetAesKeyBefore(conf.GetString("open." + mpf.EnvProjectKey() + ".aeskeybefore"))
        c.open.SetAesKeyNow(conf.GetString("open." + mpf.EnvProjectKey() + ".aeskeynow"))
        c.open.SetAppId(conf.GetString("open." + mpf.EnvProjectKey() + ".appid"))
        c.open.SetDomainMiniServers(conf.GetStringSlice("open." + mpf.EnvProjectKey() + ".miniserver"))
        c.open.SetDomainMiniWebViews(conf.GetStringSlice("open." + mpf.EnvProjectKey() + ".miniwebview"))
        c.open.SetSecret(conf.GetString("open." + mpf.EnvProjectKey() + ".secret"))
        c.open.SetToken(conf.GetString("open." + mpf.EnvProjectKey() + ".token"))
        c.open.SetUrlAuth(conf.GetString("open." + mpf.EnvProjectKey() + ".url.auth"))
        c.open.SetUrlAuthCallback(conf.GetString("open." + mpf.EnvProjectKey() + ".url.authcallback"))
        c.open.SetUrlMiniFastRegister(conf.GetString("open." + mpf.EnvProjectKey() + ".url.minifastregister"))
        c.open.SetUrlMiniRebindAdmin(conf.GetString("open." + mpf.EnvProjectKey() + ".url.minirebindadmin"))
    })

    return c.open
}

func (c *configWx) GetProvider() *configProvider {
    onceConfigProvider.Do(func() {
        conf := mpf.NewConfig().GetConfig("wx")
        c.provider.SetAesKey(conf.GetString("provider." + mpf.EnvProjectKey() + ".aeskey"))
        c.provider.SetCorpId(conf.GetString("provider." + mpf.EnvProjectKey() + ".corpid"))
        c.provider.SetCorpSecret(conf.GetString("provider." + mpf.EnvProjectKey() + ".corpsecret"))
        c.provider.SetSuiteId(conf.GetString("provider." + mpf.EnvProjectKey() + ".suiteid"))
        c.provider.SetSuiteSecret(conf.GetString("provider." + mpf.EnvProjectKey() + ".suitesecret"))
        c.provider.SetToken(conf.GetString("provider." + mpf.EnvProjectKey() + ".token"))
        c.provider.SetUrlAuthLogin(conf.GetString("provider." + mpf.EnvProjectKey() + ".url.authlogin"))
        c.provider.SetUrlAuthSuite(conf.GetString("provider." + mpf.EnvProjectKey() + ".url.authsuite"))
    })

    return c.provider
}

func (c *configWx) getLocalAccount(appId string) *configAccount {
    nowTime := time.Now().Second()
    expireTime := nowTime + project.TimeClearLocalWxAccount()
    if c.accountClearTime < nowTime {
        delList := make([]string, 0)
        for k, v := range c.accountList {
            if v.GetExpireTime() < nowTime {
                delList = append(delList, k)
            }
        }
        for _, delId := range delList {
            delete(c.accountList, delId)
        }
        c.accountClearTime = expireTime
    }

    conf, ok := c.accountList[appId]
    if !ok {
        conf = c.outer.RefreshConfigAccount(appId)
        conf.SetExpireTime(expireTime)
        c.accountList[appId] = conf
    }

    return conf
}

func (c *configWx) GetAccount(appId string) *configAccount {
    conf := c.getLocalAccount(appId)
    if conf.IsValid() {
        return conf
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "公众号配置不存在", nil))
    }
}

func (c *configWx) RemoveAccount(appId string) {
    delete(c.accountList, appId)
}

func (c *configWx) GetAccountList() map[string]*configAccount {
    return c.accountList
}

func (c *configWx) getLocalCorp(corpId string) *configCorp {
    nowTime := time.Now().Second()
    expireTime := nowTime + project.TimeClearLocalWxCorp()
    if c.corpClearTime < nowTime {
        delList := make([]string, 0)
        for k, v := range c.corpList {
            if v.GetExpireTime() < nowTime {
                delList = append(delList, k)
            }
        }
        for _, delId := range delList {
            delete(c.corpList, delId)
        }
        c.corpClearTime = expireTime
    }

    conf, ok := c.corpList[corpId]
    if !ok {
        conf = c.outer.RefreshConfigCorp(corpId)
        conf.SetExpireTime(expireTime)
        c.corpList[corpId] = conf
    }

    return conf
}

func (c *configWx) GetCorp(corpId string) *configCorp {
    conf := c.getLocalCorp(corpId)
    if conf.IsValid() {
        return conf
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "企业号配置不存在", nil))
    }
}

func (c *configWx) RemoveCorp(corpId string) {
    delete(c.corpList, corpId)
}

func (c *configWx) GetCorpList() map[string]*configCorp {
    return c.corpList
}

var (
    onceConfig         sync.Once
    onceConfigOpen     sync.Once
    onceConfigProvider sync.Once
    insConfig          *configWx
)

func init() {
    insConfig = &configWx{}
    insConfig.open = &configOpen{}
    insConfig.open.domainMiniServers = make([]string, 0)
    insConfig.open.domainMiniWebViews = make([]string, 0)
    insConfig.provider = &configProvider{}
    insConfig.accountClearTime = 0
    insConfig.accountList = make(map[string]*configAccount)
    insConfig.corpClearTime = 0
    insConfig.corpList = make(map[string]*configCorp)
}

func LoadConfig(outer IWxConfig) {
    onceConfig.Do(func() {
        insConfig.outer = outer
    })
}

func NewConfig() *configWx {
    return insConfig
}
