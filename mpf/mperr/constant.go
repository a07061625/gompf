/**
 * 自定义错误常量
 * User: 姜伟
 * Date: 2019/12/25 0025
 * Time: 13:19
 */
package mperr

const (
    /**
     * 类型
     */

    // 1-9999为框架内部预留
    TypeInnerJson      uint = 1
    TypeInnerServer    uint = 2
    TypeInnerValidator uint = 3

    // 地图
    TypeMapTencent uint = 10000
    TypeMapGaoDe   uint = 10001
    TypeMapBaiDu   uint = 10002

    // 短信
    TypeSmsAliYun uint = 10100
    TypeSmsDaYu   uint = 10101
    TypeSmsYun253 uint = 10102

    // 缓存
    TypeCacheMem   uint = 10300
    TypeCacheRedis uint = 10301

    // 数据库
    TypeDbMysql uint = 10400
    TypeDbMonGo uint = 10401

    // 配置
    TypeConfigViper uint = 10500

    // 阿里云开放平台
    TypeAliOpen uint = 10600

    // 支付宝
    TypeAliPay         uint = 10800
    TypeAliPayAuth     uint = 10801
    TypeAliPayFund     uint = 10802
    TypeAliPayLife     uint = 10803
    TypeAliPayMarket   uint = 10804
    TypeAliPayMaterial uint = 10805
    TypeAliPayTrade    uint = 10806

    // 微信
    TypeWx            uint = 11000
    TypeWxAccount     uint = 11001
    TypeWxCorp        uint = 11002
    TypeWxProvider    uint = 11003
    TypeWxMini        uint = 11004
    TypeWxOpen        uint = 11005
    TypeWxOpenAccount uint = 11006
    TypeWxOpenMini    uint = 11007

    // 打印
    TypePrintFeYin uint = 11200

    // 物流
    TypeLogisticsAMAli  uint = 11300
    TypeLogisticsKd100  uint = 11301
    TypeLogisticsKdBird uint = 11302
    TypeLogisticsTaoBao uint = 11303

    // 即时通讯
    TypeIMTencent uint = 11400

    // 货币
    TypeCurrencyAMJiSu   uint = 11500
    TypeCurrencyAMYiYuan uint = 11501

    // 腾讯云
    TypeQCloud    uint = 11600
    TypeQCloudCos uint = 11601

    // 七牛
    TypeQiNiu     uint = 11800
    TypeQiNiuKodo uint = 11801

    // 阿里云OSS
    TypeAliOss uint = 11900

    // 钉钉
    TypeDingTalk         uint = 12000
    TypeDingTalkCorp     uint = 12001
    TypeDingTalkProvider uint = 12002

    // 物联网
    TypeIotAliYun  uint = 12100
    TypeIotBaiDu   uint = 12101
    TypeIotTencent uint = 12102

    // 消息推送
    TypePushAliYun uint = 12200
    TypePushBaiDu  uint = 12201
    TypePushXinGe  uint = 12202
    TypePushJPush  uint = 12203

    // 消息队列
    TypeMQ       uint = 12300
    TypeMQRedis  uint = 12301
    TypeMQRabbit uint = 12302
)
