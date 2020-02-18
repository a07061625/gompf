// Package mperr error
// User: 姜伟
// Time: 2020-02-19 04:26:32
package mperr

// IErrorCommon 错误公共接口
type IErrorCommon interface {
    Error() string
    Unwrap() error
}

// ErrorCommon 公共错误
type ErrorCommon struct {
    Type  uint   `json:"type"`
    Title string `json:"title"`
    Code  uint   `json:"code"`
    Msg   string `json:"msg"`
    Err   error
}

// IErrorCommon 错误
func (e *ErrorCommon) Error() string {
    return e.Msg
}

// Unwrap Unwrap
// IErrorCommon 错误
func (e *ErrorCommon) Unwrap() error {
    return e.Err
}

// NewInnerJSON 内部json错误
func NewInnerJSON(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeInnerJSON, "Json错误", code, msg, err}
}

// NewInnerServer 内部服务错误
func NewInnerServer(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeInnerServer, "服务错误", code, msg, err}
}

// NewInnerValidator 内部校验器错误
func NewInnerValidator(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeInnerValidator, "校验器错误", code, msg, err}
}

// NewMapTencent 腾讯地图错误
func NewMapTencent(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeMapTencent, "腾讯地图错误", code, msg, err}
}

// NewMapBaiDu 百度地图错误
func NewMapBaiDu(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeMapBaiDu, "百度地图错误", code, msg, err}
}

// NewMapGaoDe 高德地图错误
func NewMapGaoDe(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeMapGaoDe, "高德地图错误", code, msg, err}
}

// NewSmsAliYun 阿里云短信错误
func NewSmsAliYun(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeSmsAliYun, "阿里云短信错误", code, msg, err}
}

// NewSmsDaYu 阿里大鱼短信错误
func NewSmsDaYu(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeSmsDaYu, "大鱼短信错误", code, msg, err}
}

// NewSmsYun253 253云短信错误
func NewSmsYun253(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeSmsYun253, "253云短信错误", code, msg, err}
}

// NewCacheMem memcache缓存错误
func NewCacheMem(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeCacheMem, "memcache缓存错误", code, msg, err}
}

// NewCacheRedis redis缓存错误
func NewCacheRedis(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeCacheRedis, "redis缓存错误", code, msg, err}
}

// NewDbMysql mysql数据库错误
func NewDbMysql(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeDbMysql, "mysql数据库错误", code, msg, err}
}

// NewDbMonGo mongo数据库错误
func NewDbMonGo(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeDbMonGo, "mongo数据库错误", code, msg, err}
}

// NewConfigViper viper配置错误
func NewConfigViper(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeConfigViper, "viper配置错误", code, msg, err}
}

// NewAliOpen 阿里云开放平台错误
func NewAliOpen(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeAliOpen, "阿里云开放平台错误", code, msg, err}
}

// NewAliPay 支付宝错误
func NewAliPay(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeAliPay, "支付宝错误", code, msg, err}
}

// NewAliPayAuth 支付宝授权错误
func NewAliPayAuth(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeAliPayAuth, "支付宝授权错误", code, msg, err}
}

// NewAliPayFund 支付宝资金错误
func NewAliPayFund(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeAliPayFund, "支付宝资金错误", code, msg, err}
}

// NewAliPayLife 支付宝生活号错误
func NewAliPayLife(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeAliPayLife, "支付宝生活号错误", code, msg, err}
}

// NewAliPayMarket 支付宝店铺错误
func NewAliPayMarket(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeAliPayMarket, "支付宝店铺错误", code, msg, err}
}

// NewAliPayMaterial 支付宝物料错误
func NewAliPayMaterial(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeAliPayMaterial, "支付宝物料错误", code, msg, err}
}

// NewAliPayTrade 支付宝支付错误
func NewAliPayTrade(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeAliPayTrade, "支付宝支付错误", code, msg, err}
}

// NewWx 微信错误
func NewWx(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeWx, "微信错误", code, msg, err}
}

// NewWxAccount 微信公众号错误
func NewWxAccount(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeWxAccount, "微信公众号错误", code, msg, err}
}

// NewWxCorp 微信企业号错误
func NewWxCorp(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeWxCorp, "微信企业号错误", code, msg, err}
}

// NewWxProvider 微信企业号服务商错误
func NewWxProvider(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeWxProvider, "微信企业号服务商错误", code, msg, err}
}

// NewWxMini 微信小程序错误
func NewWxMini(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeWxMini, "微信小程序错误", code, msg, err}
}

// NewWxOpen 微信开放平台错误
func NewWxOpen(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeWxOpen, "微信开放平台错误", code, msg, err}
}

// NewWxOpenAccount 微信开放平台公众号错误
func NewWxOpenAccount(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeWxOpenAccount, "微信开放平台公众号错误", code, msg, err}
}

// NewWxOpenMini 微信开放平台小程序错误
func NewWxOpenMini(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeWxOpenMini, "微信开放平台小程序错误", code, msg, err}
}

// NewPrintFeYin 飞印打印错误
func NewPrintFeYin(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypePrintFeYin, "飞印打印错误", code, msg, err}
}

// NewLogisticsAMAli 阿里云市场阿里物流错误
func NewLogisticsAMAli(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeLogisticsAMAli, "阿里云市场阿里物流错误", code, msg, err}
}

// NewLogisticsKd100 快递100物流错误
func NewLogisticsKd100(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeLogisticsKd100, "快递100物流错误", code, msg, err}
}

// NewLogisticsKdBird 快递鸟物流错误
func NewLogisticsKdBird(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeLogisticsKdBird, "快递鸟物流错误", code, msg, err}
}

// NewLogisticsTaoBao 淘宝物流错误
func NewLogisticsTaoBao(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeLogisticsTaoBao, "淘宝物流错误", code, msg, err}
}

// NewIMTencent 腾讯即时通讯错误
func NewIMTencent(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeIMTencent, "腾讯即时通讯错误", code, msg, err}
}

// NewCurrencyAMJiSu 阿里云市场极速货币错误
func NewCurrencyAMJiSu(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeCurrencyAMJiSu, "阿里云市场极速货币错误", code, msg, err}
}

// NewCurrencyAMYiYuan 阿里云市场易圆货币错误
func NewCurrencyAMYiYuan(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeCurrencyAMYiYuan, "阿里云市场易圆货币错误", code, msg, err}
}

// NewQCloud 腾讯云错误
func NewQCloud(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeQCloud, "腾讯云错误", code, msg, err}
}

// NewQCloudCos 腾讯云对象存储错误
func NewQCloudCos(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeQCloudCos, "腾讯云对象存储错误", code, msg, err}
}

// NewQiNiu 七牛云错误
func NewQiNiu(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeQiNiu, "七牛云错误", code, msg, err}
}

// NewQiNiuKodo 七牛云对象存储错误
func NewQiNiuKodo(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeQiNiuKodo, "七牛云对象存储错误", code, msg, err}
}

// NewAliOss 阿里云OSS错误
func NewAliOss(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeAliOss, "阿里云OSS错误", code, msg, err}
}

// NewDingTalk 钉钉错误
func NewDingTalk(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeDingTalk, "钉钉错误", code, msg, err}
}

// NewDingTalkCorp 钉钉企业号错误
func NewDingTalkCorp(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeDingTalkCorp, "钉钉企业号错误", code, msg, err}
}

// NewDingTalkProvider 钉钉服务商错误
func NewDingTalkProvider(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeDingTalkProvider, "钉钉服务商错误", code, msg, err}
}

// NewIotAliYun 阿里云物联网错误
func NewIotAliYun(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeIotAliYun, "阿里云物联网错误", code, msg, err}
}

// NewIotBaiDu 百度物联网错误
func NewIotBaiDu(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeIotBaiDu, "百度物联网错误", code, msg, err}
}

// NewIotTencent 腾讯物联网错误
func NewIotTencent(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeIotTencent, "腾讯物联网错误", code, msg, err}
}

// NewPushAliYun 阿里云消息推送错误
func NewPushAliYun(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypePushAliYun, "阿里云消息推送错误", code, msg, err}
}

// NewPushBaiDu 百度消息推送错误
func NewPushBaiDu(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypePushBaiDu, "百度消息推送错误", code, msg, err}
}

// NewPushXinGe 信鸽消息推送错误
func NewPushXinGe(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypePushXinGe, "信鸽消息推送错误", code, msg, err}
}

// NewPushJPush 极光消息推送错误
func NewPushJPush(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypePushJPush, "极光消息推送错误", code, msg, err}
}

// NewMQ 消息队列错误
func NewMQ(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeMQ, "消息队列错误", code, msg, err}
}

// NewMQRedis Redis消息队列错误
func NewMQRedis(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeMQRedis, "Redis消息队列错误", code, msg, err}
}

// NewMQRabbit Rabbit消息队列错误
func NewMQRabbit(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeMQRabbit, "Rabbit消息队列错误", code, msg, err}
}

// NewProtocol 协议错误
func NewProtocol(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeProtocol, "协议错误", code, msg, err}
}
