/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2019/12/27 0027
 * Time: 19:50
 */
package alipay

import (
    "regexp"
    "sync"
    "time"

    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
)

type configAccount struct {
    appId      string
    sellerId   string // 卖家ID
    urlNotify  string // 异步消息通知URL
    urlReturn  string // 同步消息通知URL
    priKeyRsa  string // rsa私钥,包含-----BEGIN RSA PRIVATE KEY-----
    pubKeyRsa  string // rsa公钥,包含-----BEGIN RSA PUBLIC KEY-----
    pubKeyAli  string // 支付宝公钥,包含-----BEGIN PUBLIC KEY-----
    valid      bool   // 配置有效状态
    expireTime int64  // 配置过期时间戳
}

func (c *configAccount) SetAppId(appId string) {
    match, _ := regexp.MatchString(`^[0-9]{16}$`, appId)
    if match {
        c.appId = appId
    } else {
        panic(mperr.NewAliPay(errorcode.AliPayParam, "app id不合法", nil))
    }
}

func (c *configAccount) GetAppId() string {
    return c.appId
}

func (c *configAccount) SetSellerId(sellerId string) {
    match, _ := regexp.MatchString(`^[0-9]{16}$`, sellerId)
    if match && (sellerId[:4] == "2008") {
        c.sellerId = sellerId
    } else {
        panic(mperr.NewAliPay(errorcode.AliPayParam, "卖家ID不合法", nil))
    }
}

func (c *configAccount) GetSellerId() string {
    return c.sellerId
}

func (c *configAccount) SetUrlNotify(urlNotify string) {
    match, _ := regexp.MatchString(project.RegexUrlHttp, urlNotify)
    if match {
        c.urlNotify = urlNotify
    } else {
        panic(mperr.NewAliPay(errorcode.AliPayParam, "异步消息通知URL不合法", nil))
    }
}

func (c *configAccount) GetUrlNotify() string {
    return c.urlNotify
}

func (c *configAccount) SetUrlReturn(urlReturn string) {
    match, _ := regexp.MatchString(project.RegexUrlHttp, urlReturn)
    if match {
        c.urlReturn = urlReturn
    } else {
        panic(mperr.NewAliPay(errorcode.AliPayParam, "同步消息通知URL不合法", nil))
    }
}

func (c *configAccount) GetUrlReturn() string {
    return c.urlReturn
}

func (c *configAccount) SetPriKeyRsa(priKeyRsa string) {
    if len(priKeyRsa) >= 1024 {
        c.priKeyRsa = priKeyRsa
    } else {
        panic(mperr.NewAliPay(errorcode.AliPayParam, "rsa私钥不合法", nil))
    }
}

func (c *configAccount) GetPriKeyRsa() string {
    return c.priKeyRsa
}

func (c *configAccount) SetPubKeyRsa(pubKeyRsa string) {
    if len(pubKeyRsa) >= 256 {
        c.pubKeyRsa = pubKeyRsa
    } else {
        panic(mperr.NewAliPay(errorcode.AliPayParam, "rsa公钥不合法", nil))
    }
}

func (c *configAccount) GetPubKeyRsa() string {
    return c.pubKeyRsa
}

func (c *configAccount) SetPubKeyAli(pubKeyAli string) {
    if len(pubKeyAli) >= 256 {
        c.pubKeyAli = pubKeyAli
    } else {
        panic(mperr.NewAliPay(errorcode.AliPayParam, "支付宝公钥不合法", nil))
    }
}

func (c *configAccount) GetPubKeyAli() string {
    return c.pubKeyAli
}

func (c *configAccount) SetValid(valid bool) {
    c.valid = valid
}

func (c *configAccount) IsValid() bool {
    return c.valid
}

func (c *configAccount) SetExpireTime(expireTime int64) {
    c.expireTime = expireTime
}

func (c *configAccount) GetExpireTime() int64 {
    return c.expireTime
}

func NewConfigAccount() *configAccount {
    return &configAccount{"", "", "", "", "", "", "", false, 0}
}

type IAliPayConfig interface {
    RefreshConfigAccount(appId string) *configAccount // 刷新账号配置
}

type configAliPay struct {
    outer            IAliPayConfig             // 项目实现的接口,用于根据项目获取到相应的配置实例
    accountClearTime int64                     // 账号本地清理时间戳
    accountList      map[string]*configAccount // 账号本地配置集合
}

func (c *configAliPay) getLocalAccount(appId string) *configAccount {
    nowTime := time.Now().Unix()
    expireTime := nowTime + project.TimeClearLocalAliPayAccount()
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

func (c *configAliPay) GetAccount(appId string) *configAccount {
    conf := c.getLocalAccount(appId)
    if conf.IsValid() {
        return conf
    } else {
        panic(mperr.NewAliPay(errorcode.AliPayParam, "账号配置不存在", nil))
    }
}

func (c *configAliPay) RemoveAccount(appId string) {
    delete(c.accountList, appId)
}

func (c *configAliPay) GetAccountList() map[string]*configAccount {
    return c.accountList
}

var (
    onceConfig sync.Once
    insConfig  *configAliPay
)

func init() {
    insConfig = &configAliPay{}
    insConfig.accountClearTime = 0
    insConfig.accountList = make(map[string]*configAccount)
}

func LoadConfig(outer IAliPayConfig) {
    onceConfig.Do(func() {
        insConfig.outer = outer
    })
}

func NewConfig() *configAliPay {
    return insConfig
}
