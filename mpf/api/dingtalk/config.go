package dingtalk

import (
    "regexp"
    "sync"
    "time"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
)

// 企业号配置
type configCorp struct {
    corpId           string                       // 企业ID
    ssoSecret        string                       // 免登密钥
    agents           map[string]map[string]string // 应用列表
    loginAppId       string                       // 登陆应用ID
    loginAppSecret   string                       // 登陆应用密钥
    loginUrlCallback string                       // 登陆应用回调地址
    valid            bool                         // 配置有效状态
    expireTime       int64                        // 配置过期时间戳
}

func (c *configCorp) SetCorpId(corpId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, corpId)
    if match {
        c.corpId = corpId
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "企业ID不合法", nil))
    }
}

func (c *configCorp) GetCorpId() string {
    return c.corpId
}

func (c *configCorp) SetSsoSecret(ssoSecret string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, ssoSecret)
    if match {
        c.ssoSecret = ssoSecret
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "免登密钥不合法", nil))
    }
}

func (c *configCorp) GetSsoSecret() string {
    return c.ssoSecret
}

func (c *configCorp) SetAgents(agents map[string]map[string]string) {
    if len(agents) > 0 {
        c.agents = agents
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "应用列表不合法", nil))
    }
}

// 获取应用信息
// 返回数据结构:
//   id: 应用ID
//   key: 应用标识
//   secret: 应用密钥
//   token: 消息校验token
//   aes_key: 加密密钥
//   callback_tags: 监听事件类型列表
//   callback_url: 回调地址
func (c *configCorp) GetAgentInfo(agentTag string) map[string]string {
    info, ok := c.agents[agentTag]
    if ok {
        return info
    } else {
        return make(map[string]string)
    }
}

func (c *configCorp) SetLoginAppId(loginAppId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, loginAppId)
    if match {
        c.loginAppId = loginAppId
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "登陆应用ID不合法", nil))
    }
}

func (c *configCorp) GetLoginAppId() string {
    return c.loginAppId
}

func (c *configCorp) SetLoginAppSecret(loginAppSecret string) {
    if len(loginAppSecret) > 0 {
        c.loginAppSecret = loginAppSecret
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "登陆应用密钥不合法", nil))
    }
}

func (c *configCorp) GetLoginAppSecret() string {
    return c.loginAppSecret
}

func (c *configCorp) SetLoginUrlCallback(loginUrlCallback string) {
    match, _ := regexp.MatchString(project.RegexURLHTTP, loginUrlCallback)
    if match {
        c.loginUrlCallback = loginUrlCallback
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "登陆应用回调地址不合法", nil))
    }
}

func (c *configCorp) GetLoginUrlCallback() string {
    return c.loginUrlCallback
}

func (c *configCorp) SetValid(valid bool) {
    c.valid = valid
}

func (c *configCorp) IsValid() bool {
    return c.valid
}

func (c *configCorp) SetExpireTime(expireTime int64) {
    c.expireTime = expireTime
}

func (c *configCorp) GetExpireTime() int64 {
    return c.expireTime
}

func NewConfigCorp() *configCorp {
    return &configCorp{"", "", make(map[string]map[string]string), "", "", "", false, 0}
}

// 服务商配置
type configProvider struct {
    corpId           string // 企业ID
    ssoSecret        string // 免登密钥
    token            string // 消息校验token
    aesKey           string // 加密密钥
    suiteId          int    // 套件ID
    suiteKey         string // 套件标识
    suiteSecret      string // 套件密钥
    loginAppId       string // 登陆应用ID
    loginAppSecret   string // 登陆应用密钥
    loginUrlCallback string // 登陆应用回调地址
}

func (c *configProvider) SetCorpId(corpId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, corpId)
    if match {
        c.corpId = corpId
    } else {
        panic(mperr.NewDingTalkProvider(errorcode.DingTalkProviderParam, "企业ID不合法", nil))
    }
}

func (c *configProvider) GetCorpId() string {
    return c.corpId
}

func (c *configProvider) SetSsoSecret(ssoSecret string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, ssoSecret)
    if match {
        c.ssoSecret = ssoSecret
    } else {
        panic(mperr.NewDingTalkProvider(errorcode.DingTalkProviderParam, "免登密钥不合法", nil))
    }
}

func (c *configProvider) GetSsoSecret() string {
    return c.ssoSecret
}

func (c *configProvider) SetToken(token string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{1,32}$`, token)
    if match {
        c.token = token
    } else {
        panic(mperr.NewDingTalkProvider(errorcode.DingTalkProviderParam, "消息校验token不合法", nil))
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
        panic(mperr.NewDingTalkProvider(errorcode.DingTalkProviderParam, "加密密钥不合法", nil))
    }
}

func (c *configProvider) GetAesKey() string {
    return c.aesKey
}

func (c *configProvider) SetSuiteId(suiteId int) {
    if suiteId > 0 {
        c.suiteId = suiteId
    } else {
        panic(mperr.NewDingTalkProvider(errorcode.DingTalkProviderParam, "套件ID不合法", nil))
    }
}

func (c *configProvider) GetSuiteId() int {
    return c.suiteId
}

func (c *configProvider) SetSuiteKey(suiteKey string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, suiteKey)
    if match {
        c.suiteKey = suiteKey
    } else {
        panic(mperr.NewDingTalkProvider(errorcode.DingTalkProviderParam, "套件标识不合法", nil))
    }
}

func (c *configProvider) GetSuiteKey() string {
    return c.suiteKey
}

func (c *configProvider) SetSuiteSecret(suiteSecret string) {
    if len(suiteSecret) > 0 {
        c.suiteSecret = suiteSecret
    } else {
        panic(mperr.NewDingTalkProvider(errorcode.DingTalkProviderParam, "套件密钥不合法", nil))
    }
}

func (c *configProvider) GetSuiteSecret() string {
    return c.suiteSecret
}

func (c *configProvider) SetLoginAppId(loginAppId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, loginAppId)
    if match {
        c.loginAppId = loginAppId
    } else {
        panic(mperr.NewDingTalkProvider(errorcode.DingTalkProviderParam, "登陆应用ID不合法", nil))
    }
}

func (c *configProvider) GetLoginAppId() string {
    return c.loginAppId
}

func (c *configProvider) SetLoginAppSecret(loginAppSecret string) {
    if len(loginAppSecret) > 0 {
        c.loginAppSecret = loginAppSecret
    } else {
        panic(mperr.NewDingTalkProvider(errorcode.DingTalkProviderParam, "登陆应用密钥不合法", nil))
    }
}

func (c *configProvider) GetLoginAppSecret() string {
    return c.loginAppSecret
}

func (c *configProvider) SetLoginUrlCallback(loginUrlCallback string) {
    match, _ := regexp.MatchString(project.RegexURLHTTP, loginUrlCallback)
    if match {
        c.loginUrlCallback = loginUrlCallback
    } else {
        panic(mperr.NewDingTalkProvider(errorcode.DingTalkProviderParam, "登陆应用回调地址不合法", nil))
    }
}

func (c *configProvider) GetLoginUrlCallback() string {
    return c.loginUrlCallback
}

type IDingTalkConfig interface {
    RefreshConfigCorp(corpId string) *configCorp // 刷新企业号配置
}

type configDingTalk struct {
    outer         IDingTalkConfig
    provider      *configProvider
    corpClearTime int64                  // 企业号本地清理时间戳
    corpList      map[string]*configCorp // 企业号本地配置集合
}

func (c *configDingTalk) getLocalCorp(corpId string) *configCorp {
    nowTime := time.Now().Unix()
    expireTime := nowTime + project.TimeClearLocalDingTalkCorp()
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

func (c *configDingTalk) GetCorp(corpId string) *configCorp {
    conf := c.getLocalCorp(corpId)
    if conf.IsValid() {
        return conf
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "企业号配置不存在", nil))
    }
}

func (c *configDingTalk) RemoveCorp(corpId string) {
    delete(c.corpList, corpId)
}

func (c *configDingTalk) GetCorpList() map[string]*configCorp {
    return c.corpList
}

func (c *configDingTalk) GetProvider() *configProvider {
    onceConfigProvider.Do(func() {
        conf := mpf.NewConfig().GetConfig("dingtalk")
        c.provider.SetCorpId(conf.GetString("provider." + mpf.EnvProjectKey() + ".corp.id"))
        c.provider.SetSsoSecret(conf.GetString("provider." + mpf.EnvProjectKey() + ".sso.secret"))
        c.provider.SetToken(conf.GetString("provider." + mpf.EnvProjectKey() + ".token"))
        c.provider.SetAesKey(conf.GetString("provider." + mpf.EnvProjectKey() + ".aeskey"))
        c.provider.SetSuiteId(conf.GetInt("provider." + mpf.EnvProjectKey() + ".suite.id"))
        c.provider.SetSuiteKey(conf.GetString("provider." + mpf.EnvProjectKey() + ".suite.key"))
        c.provider.SetSuiteSecret(conf.GetString("provider." + mpf.EnvProjectKey() + ".suite.secret"))
        c.provider.SetLoginAppId(conf.GetString("provider." + mpf.EnvProjectKey() + ".login.appid"))
        c.provider.SetLoginAppSecret(conf.GetString("provider." + mpf.EnvProjectKey() + ".login.appsecret"))
        c.provider.SetLoginUrlCallback(conf.GetString("provider." + mpf.EnvProjectKey() + ".login.urlcallback"))
    })

    return c.provider
}

var (
    onceConfig         sync.Once
    onceConfigProvider sync.Once
    insConfig          *configDingTalk
)

func init() {
    insConfig = &configDingTalk{}
    insConfig.provider = &configProvider{}
    insConfig.corpClearTime = 0
    insConfig.corpList = make(map[string]*configCorp)
}

func LoadConfig(outer IDingTalkConfig) {
    onceConfig.Do(func() {
        insConfig.outer = outer
    })
}

func NewConfig() *configDingTalk {
    return insConfig
}
