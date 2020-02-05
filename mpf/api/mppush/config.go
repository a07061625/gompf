/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/2 0002
 * Time: 22:27
 */
package mppush

import (
    "encoding/base64"
    "regexp"
    "runtime"
    "sync"
    "time"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
)

type configBaiDu struct {
    appKey     string // 应用ID
    appSecret  string // 应用密钥
    deviceType string // 设备类型 0:2.0升级到3.0的应用设置该值 3:安卓 4:苹果
    userAgent  string // HTTP客户端信息
}

func (c *configBaiDu) SetAppKey(appKey string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, appKey)
    if match {
        c.appKey = appKey
    } else {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "应用ID不合法", nil))
    }
}

func (c *configBaiDu) GetAppKey() string {
    return c.appKey
}

func (c *configBaiDu) SetAppSecret(appSecret string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, appSecret)
    if match {
        c.appSecret = appSecret
    } else {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "应用密钥不合法", nil))
    }
}

func (c *configBaiDu) GetAppSecret() string {
    return c.appSecret
}

func (c *configBaiDu) SetDeviceType(deviceType string) {
    switch deviceType {
    case BaiDuDeviceTypeAll:
        c.deviceType = deviceType
    case BaiDuDeviceTypeAndroid:
        c.deviceType = deviceType
    case BaiDuDeviceTypeIOS:
        c.deviceType = deviceType
    default:
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "设备类型不合法", nil))
    }
}

func (c *configBaiDu) GetDeviceType() string {
    return c.deviceType
}

func (c *configBaiDu) GetUserAgent() string {
    return c.userAgent
}

type configXinGe struct {
    appId     string // 应用ID
    appSecret string // 应用密钥
    appAuth   string
    platform  string // 平台类型
}

func (c *configXinGe) SetAppInfo(appId, appSecret string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, appId)
    if !match {
        panic(mperr.NewPushXinGe(errorcode.PushXinGeParam, "应用ID不合法", nil))
    }
    match, _ = regexp.MatchString(project.RegexDigitAlpha, appSecret)
    if !match {
        panic(mperr.NewPushXinGe(errorcode.PushXinGeParam, "应用密钥不合法", nil))
    }

    authStr := appId + ":" + appSecret
    c.appId = appId
    c.appSecret = appSecret
    c.appAuth = "Basic " + base64.StdEncoding.EncodeToString([]byte(authStr))
}

func (c *configXinGe) GetAppId() string {
    return c.appId
}

func (c *configXinGe) GetAppSecret() string {
    return c.appSecret
}

func (c *configXinGe) GetAppAuth() string {
    return c.appAuth
}

func (c *configXinGe) SetPlatform(platform string) {
    if (platform == XinGePlatformTypeAndroid) || (platform == XinGePlatformTypeIOS) {
        c.platform = platform
    } else {
        panic(mperr.NewPushXinGe(errorcode.PushXinGeParam, "平台类型不合法", nil))
    }
}

func (c *configXinGe) GetPlatform() string {
    return c.platform
}

type configJPushDev struct {
    key    string // 标识
    secret string // 密钥
    auth   string // 密文
}

func (c *configJPushDev) SetKeyAndSecret(key, secret string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, key)
    if !match {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "标识不合法", nil))
    }
    match, _ = regexp.MatchString(project.RegexDigitAlpha, secret)
    if !match {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "密钥不合法", nil))
    }

    authStr := key + ":" + secret
    c.key = key
    c.secret = secret
    c.auth = "Basic " + base64.StdEncoding.EncodeToString([]byte(authStr))
}

func (c *configJPushDev) GetKey() string {
    return c.key
}

func (c *configJPushDev) GetSecret() string {
    return c.secret
}

func (c *configJPushDev) GetAuth() string {
    return c.auth
}

type configJPushApp struct {
    key        string // 标识
    secret     string // 密钥
    auth       string // 密文
    valid      bool   // 配置有效状态
    expireTime int    // 配置过期时间戳
}

func (c *configJPushApp) SetKeyAndSecret(key, secret string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, key)
    if !match {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "标识不合法", nil))
    }
    match, _ = regexp.MatchString(project.RegexDigitAlpha, secret)
    if !match {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "密钥不合法", nil))
    }

    authStr := key + ":" + secret
    c.key = key
    c.secret = secret
    c.auth = "Basic " + base64.StdEncoding.EncodeToString([]byte(authStr))
}

func (c *configJPushApp) GetKey() string {
    return c.key
}

func (c *configJPushApp) GetSecret() string {
    return c.secret
}

func (c *configJPushApp) GetAuth() string {
    return c.auth
}

func (c *configJPushApp) SetValid(valid bool) {
    c.valid = valid
}

func (c *configJPushApp) IsValid() bool {
    return c.valid
}

func (c *configJPushApp) SetExpireTime(expireTime int) {
    c.expireTime = expireTime
}

func (c *configJPushApp) GetExpireTime() int {
    return c.expireTime
}

func NewConfigJPushApp() *configJPushApp {
    return &configJPushApp{"", "", "", false, 0}
}

type configJPushGroup struct {
    key        string // 标识
    secret     string // 密钥
    auth       string // 密文
    valid      bool   // 配置有效状态
    expireTime int    // 配置过期时间戳
}

func (c *configJPushGroup) SetKeyAndSecret(key, secret string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, key)
    if !match {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "标识不合法", nil))
    }
    match, _ = regexp.MatchString(project.RegexDigitAlpha, secret)
    if !match {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "密钥不合法", nil))
    }

    authStr := "group-" + key + ":" + secret
    c.key = key
    c.secret = secret
    c.auth = "Basic " + base64.StdEncoding.EncodeToString([]byte(authStr))
}

func (c *configJPushGroup) GetKey() string {
    return c.key
}

func (c *configJPushGroup) GetSecret() string {
    return c.secret
}

func (c *configJPushGroup) GetAuth() string {
    return c.auth
}

func (c *configJPushGroup) SetValid(valid bool) {
    c.valid = valid
}

func (c *configJPushGroup) IsValid() bool {
    return c.valid
}

func (c *configJPushGroup) SetExpireTime(expireTime int) {
    c.expireTime = expireTime
}

func (c *configJPushGroup) GetExpireTime() int {
    return c.expireTime
}

func NewConfigJPushGroup() *configJPushGroup {
    return &configJPushGroup{"", "", "", false, 0}
}

type IPushConfig interface {
    RefreshConfigJPushApp(key string) *configJPushApp     // 刷新极光App配置
    RefreshConfigJPushGroup(key string) *configJPushGroup // 刷新极光分组配置
}

type configPush struct {
    outer               IPushConfig
    baiDu               *configBaiDu
    xinGeAndroid        *configXinGe
    xinGeIos            *configXinGe
    jPushDev            *configJPushDev
    jPushAppClearTime   int                          // 极光App本地清理时间戳
    jPushAppList        map[string]*configJPushApp   // 极光App本地配置集合
    jPushGroupClearTime int                          // 极光分组本地清理时间戳
    jPushGroupList      map[string]*configJPushGroup // 极光分组本地配置集合
}

func (c *configPush) GetBaiDu() *configBaiDu {
    onceConfigBaiDu.Do(func() {
        conf := mpf.NewConfig().GetConfig("mppush")
        c.baiDu.SetAppKey(conf.GetString("baidu." + mpf.EnvProjectKey() + ".app.key"))
        c.baiDu.SetAppSecret(conf.GetString("baidu." + mpf.EnvProjectKey() + ".app.secret"))
        c.baiDu.SetDeviceType(conf.GetString("baidu." + mpf.EnvProjectKey() + ".device.type"))
    })

    return c.baiDu
}

func (c *configPush) GetXinGeAndroid() *configXinGe {
    onceConfigXinGeAndroid.Do(func() {
        conf := mpf.NewConfig().GetConfig("mppush")
        appId := conf.GetString("xinge." + mpf.EnvProjectKey() + ".android.app.id")
        appSecret := conf.GetString("xinge." + mpf.EnvProjectKey() + ".android.app.secret")
        c.xinGeAndroid.SetAppInfo(appId, appSecret)
        c.xinGeAndroid.SetPlatform(XinGePlatformTypeAndroid)
    })

    return c.xinGeAndroid
}

func (c *configPush) GetXinGeIos() *configXinGe {
    onceConfigXinGeIos.Do(func() {
        conf := mpf.NewConfig().GetConfig("mppush")
        appId := conf.GetString("xinge." + mpf.EnvProjectKey() + ".ios.app.id")
        appSecret := conf.GetString("xinge." + mpf.EnvProjectKey() + ".ios.app.secret")
        c.xinGeIos.SetAppInfo(appId, appSecret)
        c.xinGeIos.SetPlatform(XinGePlatformTypeIOS)
    })

    return c.xinGeIos
}

func (c *configPush) GetJPushDev() *configJPushDev {
    onceConfigJPushDev.Do(func() {
        conf := mpf.NewConfig().GetConfig("mppush")
        key := conf.GetString("jpush." + mpf.EnvProjectKey() + ".dev.key")
        secret := conf.GetString("jpush." + mpf.EnvProjectKey() + ".dev.secret")
        c.jPushDev.SetKeyAndSecret(key, secret)
    })

    return c.jPushDev
}

func (c *configPush) getLocalJPushApp(key string) *configJPushApp {
    nowTime := time.Now().Second()
    expireTime := nowTime + project.TimeClearLocalJPushApp()
    if c.jPushAppClearTime < nowTime {
        delList := make([]string, 0)
        for k, v := range c.jPushAppList {
            if v.GetExpireTime() < nowTime {
                delList = append(delList, k)
            }
        }
        for _, delId := range delList {
            delete(c.jPushAppList, delId)
        }
        c.jPushAppClearTime = expireTime
    }

    conf, ok := c.jPushAppList[key]
    if !ok {
        conf = c.outer.RefreshConfigJPushApp(key)
        conf.SetExpireTime(expireTime)
        c.jPushAppList[key] = conf
    }

    return conf
}

func (c *configPush) GetJPushApp(key string) *configJPushApp {
    conf := c.getLocalJPushApp(key)
    if conf.IsValid() {
        return conf
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "极光app配置不存在", nil))
    }
}

func (c *configPush) RemoveJPushApp(key string) {
    delete(c.jPushAppList, key)
}

func (c *configPush) GetJPushAppList() map[string]*configJPushApp {
    return c.jPushAppList
}

func (c *configPush) getLocalJPushGroup(key string) *configJPushGroup {
    nowTime := time.Now().Second()
    expireTime := nowTime + project.TimeClearLocalJPushGroup()
    if c.jPushGroupClearTime < nowTime {
        delList := make([]string, 0)
        for k, v := range c.jPushGroupList {
            if v.GetExpireTime() < nowTime {
                delList = append(delList, k)
            }
        }
        for _, delId := range delList {
            delete(c.jPushGroupList, delId)
        }
        c.jPushGroupClearTime = expireTime
    }

    conf, ok := c.jPushGroupList[key]
    if !ok {
        conf = c.outer.RefreshConfigJPushGroup(key)
        conf.SetExpireTime(expireTime)
        c.jPushGroupList[key] = conf
    }

    return conf
}

func (c *configPush) GetJPushGroup(key string) *configJPushGroup {
    conf := c.getLocalJPushGroup(key)
    if conf.IsValid() {
        return conf
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "极光分组配置不存在", nil))
    }
}

func (c *configPush) RemoveJPushGroup(key string) {
    delete(c.jPushGroupList, key)
}

func (c *configPush) GetJPushGroupList() map[string]*configJPushGroup {
    return c.jPushGroupList
}

var (
    onceConfig             sync.Once
    onceConfigBaiDu        sync.Once
    onceConfigXinGeAndroid sync.Once
    onceConfigXinGeIos     sync.Once
    onceConfigJPushDev     sync.Once
    insConfig              *configPush
)

func init() {
    insConfig = &configPush{}
    insConfig.baiDu = &configBaiDu{}
    insConfig.baiDu.userAgent = "BCCS_SDK/3.0 (Darwin; Darwin Kernel Version 14.0.0; x86_64) GO/" + runtime.Version() + " (Baidu Push Server SDK V3.0.0) cli/Unknown ZEND/2.6.0"
    insConfig.xinGeAndroid = &configXinGe{}
    insConfig.xinGeIos = &configXinGe{}
    insConfig.jPushDev = &configJPushDev{}
    insConfig.jPushAppClearTime = 0
    insConfig.jPushAppList = make(map[string]*configJPushApp)
    insConfig.jPushGroupClearTime = 0
    insConfig.jPushGroupList = make(map[string]*configJPushGroup)
}

func LoadConfig(outer IPushConfig) {
    onceConfig.Do(func() {
        insConfig.outer = outer
    })
}

func NewConfig() *configPush {
    return insConfig
}
