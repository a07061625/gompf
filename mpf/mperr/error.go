/**
 * 自定义错误结构体
 * User: 姜伟
 * Date: 2019/12/25 0025
 * Time: 13:18
 */
package mperr

type IErrorCommon interface {
    Error() string
    Unwrap() error
}

type ErrorCommon struct {
    Type  uint   `json:"type"`
    Title string `json:"title"`
    Code  uint   `json:"code"`
    Msg   string `json:"msg"`
    Err   error
}

func (e *ErrorCommon) Error() string {
    return e.Msg
}

func (e *ErrorCommon) Unwrap() error {
    return e.Err
}

func NewInnerJson(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeInnerJson, "Json错误", code, msg, err}
}

func NewInnerServer(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeInnerServer, "服务错误", code, msg, err}
}

func NewInnerValidator(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeInnerValidator, "校验器错误", code, msg, err}
}

func NewMapTencent(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeMapTencent, "腾讯地图错误", code, msg, err}
}

func NewMapBaiDu(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeMapBaiDu, "百度地图错误", code, msg, err}
}

func NewMapGaoDe(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeMapGaoDe, "高德地图错误", code, msg, err}
}

func NewSmsAliYun(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeSmsAliYun, "阿里云短信错误", code, msg, err}
}

func NewSmsDaYu(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeSmsDaYu, "大鱼短信错误", code, msg, err}
}

func NewSmsYun253(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeSmsYun253, "253云短信错误", code, msg, err}
}

func NewCacheMem(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeCacheMem, "memcache缓存错误", code, msg, err}
}

func NewCacheRedis(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeCacheRedis, "redis缓存错误", code, msg, err}
}

func NewDbMysql(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeDbMysql, "mysql数据库错误", code, msg, err}
}

func NewDbMonGo(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeDbMonGo, "mongo数据库错误", code, msg, err}
}

func NewConfigViper(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeConfigViper, "viper配置错误", code, msg, err}
}

func NewAliOpen(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeAliOpen, "阿里云开放平台错误", code, msg, err}
}

func NewAliPay(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeAliPay, "支付宝错误", code, msg, err}
}

func NewAliPayAuth(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeAliPayAuth, "支付宝授权错误", code, msg, err}
}

func NewAliPayFund(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeAliPayFund, "支付宝资金错误", code, msg, err}
}

func NewAliPayLife(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeAliPayLife, "支付宝生活号错误", code, msg, err}
}

func NewAliPayMarket(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeAliPayMarket, "支付宝店铺错误", code, msg, err}
}

func NewAliPayMaterial(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeAliPayMaterial, "支付宝物料错误", code, msg, err}
}

func NewAliPayTrade(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeAliPayTrade, "支付宝支付错误", code, msg, err}
}

func NewWx(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeWx, "微信错误", code, msg, err}
}

func NewWxAccount(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeWxAccount, "微信公众号错误", code, msg, err}
}

func NewWxCorp(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeWxCorp, "微信企业号错误", code, msg, err}
}

func NewWxProvider(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeWxProvider, "微信企业号服务商错误", code, msg, err}
}

func NewWxMini(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeWxMini, "微信小程序错误", code, msg, err}
}

func NewWxOpen(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeWxOpen, "微信开放平台错误", code, msg, err}
}

func NewWxOpenAccount(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeWxOpenAccount, "微信开放平台公众号错误", code, msg, err}
}

func NewWxOpenMini(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeWxOpenMini, "微信开放平台小程序错误", code, msg, err}
}

func NewPrintFeYin(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypePrintFeYin, "飞印打印错误", code, msg, err}
}

func NewLogisticsAMAli(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeLogisticsAMAli, "阿里云市场阿里物流错误", code, msg, err}
}

func NewLogisticsKd100(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeLogisticsKd100, "快递100物流错误", code, msg, err}
}

func NewLogisticsKdBird(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeLogisticsKdBird, "快递鸟物流错误", code, msg, err}
}

func NewLogisticsTaoBao(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeLogisticsTaoBao, "淘宝物流错误", code, msg, err}
}

func NewIMTencent(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeIMTencent, "腾讯即时通讯错误", code, msg, err}
}

func NewCurrencyAMJiSu(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeCurrencyAMJiSu, "阿里云市场极速货币错误", code, msg, err}
}

func NewCurrencyAMYiYuan(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeCurrencyAMYiYuan, "阿里云市场易圆货币错误", code, msg, err}
}

func NewQCloud(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeQCloud, "腾讯云错误", code, msg, err}
}

func NewQCloudCos(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeQCloudCos, "腾讯云对象存储错误", code, msg, err}
}

func NewQiNiu(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeQiNiu, "七牛云错误", code, msg, err}
}

func NewQiNiuKodo(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeQiNiuKodo, "七牛云对象存储错误", code, msg, err}
}

func NewAliOss(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeAliOss, "阿里云OSS错误", code, msg, err}
}

func NewDingTalk(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeDingTalk, "钉钉错误", code, msg, err}
}

func NewDingTalkCorp(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeDingTalkCorp, "钉钉企业号错误", code, msg, err}
}

func NewDingTalkProvider(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeDingTalkProvider, "钉钉服务商错误", code, msg, err}
}

func NewIotAliYun(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeIotAliYun, "阿里云物联网错误", code, msg, err}
}

func NewIotBaiDu(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeIotBaiDu, "百度物联网错误", code, msg, err}
}

func NewIotTencent(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeIotTencent, "腾讯物联网错误", code, msg, err}
}

func NewPushAliYun(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypePushAliYun, "阿里云消息推送错误", code, msg, err}
}

func NewPushBaiDu(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypePushBaiDu, "百度消息推送错误", code, msg, err}
}

func NewPushXinGe(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypePushXinGe, "信鸽消息推送错误", code, msg, err}
}

func NewPushJPush(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypePushJPush, "极光消息推送错误", code, msg, err}
}

func NewMQ(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeMQ, "消息队列错误", code, msg, err}
}

func NewMQRedis(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeMQRedis, "Redis消息队列错误", code, msg, err}
}

func NewMQRabbit(code uint, msg string, err error) *ErrorCommon {
    return &ErrorCommon{TypeMQRabbit, "Rabbit消息队列错误", code, msg, err}
}
