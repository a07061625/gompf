// Package mperr constant
// User: 姜伟
// Time: 2020-02-19 04:51:13
package mperr

const (
    // 类型,1-9999为框架内部预留

    // TypeInnerJSON JSON
    TypeInnerJSON uint = 1
    // TypeInnerServer 服务
    TypeInnerServer uint = 2
    // TypeInnerValidator 校验器
    TypeInnerValidator uint = 3

    // 地图

    // TypeMapTencent 腾讯
    TypeMapTencent uint = 10000
    // TypeMapGaoDe 高德
    TypeMapGaoDe uint = 10001
    // TypeMapBaiDu 百度
    TypeMapBaiDu uint = 10002

    // 短信

    // TypeSmsAliYun 阿里云
    TypeSmsAliYun uint = 10100
    // TypeSmsDaYu 阿里大鱼
    TypeSmsDaYu uint = 10101
    // TypeSmsYun253 253云
    TypeSmsYun253 uint = 10102

    // 缓存

    // TypeCacheMem memcache
    TypeCacheMem uint = 10300
    // TypeCacheRedis redis
    TypeCacheRedis uint = 10301

    // 数据库

    // TypeDbMysql mysql
    TypeDbMysql uint = 10400
    // TypeDbMonGo mongo
    TypeDbMonGo uint = 10401

    // 配置

    // TypeConfigViper viper
    TypeConfigViper uint = 10500

    // 阿里云开放平台

    // TypeAliOpen 阿里云开放平台
    TypeAliOpen uint = 10600

    // 支付宝

    // TypeAliPay 支付宝
    TypeAliPay uint = 10800
    // TypeAliPayAuth 授权
    TypeAliPayAuth uint = 10801
    // TypeAliPayFund 资金
    TypeAliPayFund uint = 10802
    // TypeAliPayLife 生活号
    TypeAliPayLife uint = 10803
    // TypeAliPayMarket 店铺
    TypeAliPayMarket uint = 10804
    // TypeAliPayMaterial 物料
    TypeAliPayMaterial uint = 10805
    // TypeAliPayTrade 支付
    TypeAliPayTrade uint = 10806

    // 微信

    // TypeWx 微信
    TypeWx uint = 11000
    // TypeWxAccount 公众号
    TypeWxAccount uint = 11001
    // TypeWxCorp 企业号
    TypeWxCorp uint = 11002
    // TypeWxProvider 企业号服务商
    TypeWxProvider uint = 11003
    // TypeWxMini 小程序
    TypeWxMini uint = 11004
    // TypeWxOpen 开放平台
    TypeWxOpen uint = 11005
    // TypeWxOpenAccount 开放平台公众号
    TypeWxOpenAccount uint = 11006
    // TypeWxOpenMini 开放平台小程序
    TypeWxOpenMini uint = 11007

    // 打印

    // TypePrintFeYin 飞印
    TypePrintFeYin uint = 11200

    // 物流

    // TypeLogisticsAMAli 阿里云市场阿里
    TypeLogisticsAMAli uint = 11300
    // TypeLogisticsKd100 快递100
    TypeLogisticsKd100 uint = 11301
    // TypeLogisticsKdBird 快递鸟
    TypeLogisticsKdBird uint = 11302
    // TypeLogisticsTaoBao 淘宝
    TypeLogisticsTaoBao uint = 11303

    // 即时通讯

    // TypeIMTencent 腾讯
    TypeIMTencent uint = 11400

    // 货币

    // TypeCurrencyAMJiSu 阿里云市场极速
    TypeCurrencyAMJiSu uint = 11500
    // TypeCurrencyAMYiYuan 阿里云市场易圆
    TypeCurrencyAMYiYuan uint = 11501

    // 腾讯云

    // TypeQCloud 腾讯云
    TypeQCloud uint = 11600
    // TypeQCloudCos 对象存储
    TypeQCloudCos uint = 11601

    // 七牛

    // TypeQiNiu 七牛云
    TypeQiNiu uint = 11800
    // TypeQiNiuKodo 对象存储
    TypeQiNiuKodo uint = 11801

    // 阿里云OSS

    // TypeAliOss 阿里云OSS
    TypeAliOss uint = 11900

    // 钉钉

    // TypeDingTalk 钉钉
    TypeDingTalk uint = 12000
    // TypeDingTalkCorp 企业号
    TypeDingTalkCorp uint = 12001
    // TypeDingTalkProvider 服务商
    TypeDingTalkProvider uint = 12002

    // 物联网

    // TypeIotAliYun 阿里云
    TypeIotAliYun uint = 12100
    // TypeIotBaiDu 百度
    TypeIotBaiDu uint = 12101
    // TypeIotTencent 腾讯
    TypeIotTencent uint = 12102

    // 消息推送

    // TypePushAliYun 阿里云
    TypePushAliYun uint = 12200
    // TypePushBaiDu 百度
    TypePushBaiDu uint = 12201
    // TypePushXinGe 信鸽
    TypePushXinGe uint = 12202
    // TypePushJPush 极光
    TypePushJPush uint = 12203

    // 消息队列

    // TypeMQ 消息队列
    TypeMQ uint = 12300
    // TypeMQRedis redis
    TypeMQRedis uint = 12301
    // TypeMQRabbit rabbit
    TypeMQRabbit uint = 12302

    // 协议

    // TypeProtocol 协议
    TypeProtocol uint = 12400
)
