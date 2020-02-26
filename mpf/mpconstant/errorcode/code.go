// Package errorcode code
// User: 姜伟
// Time: 2020-02-26 14:47:37
package errorcode

const (
    // 公共状态码,取值范围10000-99999
    // 公共基础

    // CommonBaseSuccess 成功
    CommonBaseSuccess uint = 0
    // CommonBaseMin 最小错误码
    CommonBaseMin uint = 10000
    // CommonBaseInner 内部错误
    CommonBaseInner uint = 10000
    // CommonBaseServer 服务错误
    CommonBaseServer uint = 10001
    // CommonBaseWebSocket WebSocket服务错误
    CommonBaseWebSocket uint = 10002
    // CommonBaseFuse 熔断器错误
    CommonBaseFuse uint = 10003

    // 公共请求

    // CommonRequestFail 请求失败
    CommonRequestFail uint = 30001
    // CommonRequestParam 请求参数错误
    CommonRequestParam uint = 30002
    // CommonRequestTimeout 请求超时
    CommonRequestTimeout uint = 30003
    // CommonRequestResourceEmpty 请求资源不存在
    CommonRequestResourceEmpty uint = 30004
    // CommonRequestFormat 请求格式错误
    CommonRequestFormat uint = 30005
    // CommonRequestMethod 请求方式错误
    CommonRequestMethod uint = 30006

    // 公共响应

    // CommonResponseEmpty 响应为空
    CommonResponseEmpty uint = 30050
    // CommonResponseFormat 响应格式错误
    CommonResponseFormat uint = 30051

    // 公共json

    // CommonJSONMarshal 压缩错误
    CommonJSONMarshal uint = 30100
    // CommonJSONUnmarshal 解析错误
    CommonJSONUnmarshal uint = 30101

    // 公共校验器

    // CommonValidatorSign 校验器-签名错误
    CommonValidatorSign uint = 30150
    // CommonValidatorRule 校验器-规则错误
    CommonValidatorRule uint = 30151

    // 缓存,取值范围100000-100099

    // CacheRedisConnect redis-连接错误
    CacheRedisConnect uint = 100000
    // CacheRedisOperate redis-操作错误
    CacheRedisOperate uint = 100001
    // CacheMemCacheConnect memcache-连接错误
    CacheMemCacheConnect uint = 100010
    // CacheMemCacheOperate memcache-操作错误
    CacheMemCacheOperate uint = 100011

    // 配置,取值范围100100-100199

    // ConfigViperGet viper-获取错误
    ConfigViperGet uint = 100100

    // 数据库,取值范围100200-100299

    // DbMysqlConnect mysql-连接错误
    DbMysqlConnect uint = 100200
    // DbMysqlOperate mysql-操作错误
    DbMysqlOperate uint = 100201
    // DbMysqlOperateInsert mysql-插入操作错误
    DbMysqlOperateInsert uint = 100202
    // DbMysqlOperateDelete mysql-删除操作错误
    DbMysqlOperateDelete uint = 100203
    // DbMysqlOperateUpdate mysql-更新操作错误
    DbMysqlOperateUpdate uint = 100204
    // DbMysqlOperateSelect mysql-查询操作错误
    DbMysqlOperateSelect uint = 100205
    // DbMonGoConnect mongo-连接错误
    DbMonGoConnect uint = 100210
    // DbMonGoOperate mongo-操作错误
    DbMonGoOperate uint = 100211
    // DbMonGoOperateInsert mongo-插入操作错误
    DbMonGoOperateInsert uint = 100212
    // DbMonGoOperateDelete mongo-删除操作错误
    DbMonGoOperateDelete uint = 100213
    // DbMonGoOperateUpdate mongo-更新操作错误
    DbMonGoOperateUpdate uint = 100214
    // DbMonGoOperateSelect mongo-查询操作错误
    DbMonGoOperateSelect uint = 100215

    // 日志,取值范围100300-100399

    // LogZapConfig zap-配置错误
    LogZapConfig uint = 100300

    // 短信,取值范围100400-100799

    // SmsDaYuParam 大鱼-参数错误
    SmsDaYuParam uint = 100400
    // SmsDaYuRequestPost 大鱼-post请求出错
    SmsDaYuRequestPost uint = 100401
    // SmsDaYuRequestGet 大鱼-get请求出错
    SmsDaYuRequestGet uint = 100402
    // SmsDaYuRequest 大鱼-请求出错
    SmsDaYuRequest uint = 100403
    // SmsAliYunParam 阿里云-参数错误
    SmsAliYunParam uint = 100410
    // SmsAliYunRequestPost 阿里云-post请求出错
    SmsAliYunRequestPost uint = 100411
    // SmsAliYunRequestGet 阿里云-get请求出错
    SmsAliYunRequestGet uint = 100412
    // SmsAliYunRequest 阿里云-请求出错
    SmsAliYunRequest uint = 100413
    // SmsYun253Param 253云-参数错误
    SmsYun253Param uint = 100420
    // SmsYun253RequestPost 253云-post请求出错
    SmsYun253RequestPost uint = 100421
    // SmsYun253RequestGet 253云-get请求出错
    SmsYun253RequestGet uint = 100422
    // SmsYun253Request 253云-请求出错
    SmsYun253Request uint = 100423

    // 地图,取值范围100800-100999

    // MapTencentParam 腾讯-参数错误
    MapTencentParam uint = 100800
    // MapTencentRequestPost 腾讯-post请求出错
    MapTencentRequestPost uint = 100801
    // MapTencentRequestGet 腾讯-get请求出错
    MapTencentRequestGet uint = 100802
    // MapTencentRequest 腾讯-请求出错
    MapTencentRequest uint = 100803
    // MapBaiDuParam 百度-参数错误
    MapBaiDuParam uint = 100810
    // MapBaiDuRequestPost 百度-post请求出错
    MapBaiDuRequestPost uint = 100811
    // MapBaiDuRequestGet 百度-get请求出错
    MapBaiDuRequestGet uint = 100812
    // MapBaiDuRequest 百度-请求出错
    MapBaiDuRequest uint = 100813
    // MapGaoDeParam 高德-参数错误
    MapGaoDeParam uint = 100820
    // MapGaoDeRequestPost 高德-post请求出错
    MapGaoDeRequestPost uint = 100821
    // MapGaoDeRequestGet 高德-get请求出错
    MapGaoDeRequestGet uint = 100822
    // MapGaoDeRequest 高德-请求出错
    MapGaoDeRequest uint = 100823

    // 打印,取值范围101000-101099

    // PrintFeYinParam 飞印-参数错误
    PrintFeYinParam uint = 101000
    // PrintFeYinRequestPost 飞印-post请求出错
    PrintFeYinRequestPost uint = 101001
    // PrintFeYinRequestGet 飞印-get请求出错
    PrintFeYinRequestGet uint = 101002
    // PrintFeYinRequest 飞印-请求出错
    PrintFeYinRequest uint = 101003

    // 微信,取值范围101100-101299

    // WxParam 微信-参数错误
    WxParam uint = 101100
    // WxRequestPost 微信-post请求出错
    WxRequestPost uint = 101101
    // WxRequestGet 微信-get请求出错
    WxRequestGet uint = 101102
    // WxRequest 微信-请求出错
    WxRequest uint = 101103
    // WxAccountParam 公众号-参数错误
    WxAccountParam uint = 101130
    // WxAccountRequestPost 公众号-post请求出错
    WxAccountRequestPost uint = 101131
    // WxAccountRequestGet 公众号-get请求出错
    WxAccountRequestGet uint = 101132
    // WxAccountRequest 公众号-请求出错
    WxAccountRequest uint = 101133
    // WxCorpParam 企业号-参数错误
    WxCorpParam uint = 101140
    // WxCorpRequestPost 企业号-post请求出错
    WxCorpRequestPost uint = 101141
    // WxCorpRequestGet 企业号-get请求出错
    WxCorpRequestGet uint = 101142
    // WxCorpRequest 企业号-请求出错
    WxCorpRequest uint = 101143
    // WxMiniParam 小程序-参数错误
    WxMiniParam uint = 101150
    // WxMiniRequestPost 小程序-post请求出错
    WxMiniRequestPost uint = 101151
    // WxMiniRequestGet 小程序-get请求出错
    WxMiniRequestGet uint = 101152
    // WxMiniRequest 小程序-请求出错
    WxMiniRequest uint = 101153
    // WxOpenParam 第三方开放平台-参数错误
    WxOpenParam uint = 101160
    // WxOpenRequestPost 第三方开放平台-post请求出错
    WxOpenRequestPost uint = 101161
    // WxOpenRequestGet 第三方开放平台-get请求出错
    WxOpenRequestGet uint = 101162
    // WxOpenRequest 第三方开放平台-请求出错
    WxOpenRequest uint = 101163
    // WxProviderParam 企业服务商-参数错误
    WxProviderParam uint = 101170
    // WxProviderRequestPost 企业服务商-post请求出错
    WxProviderRequestPost uint = 101171
    // WxProviderRequestGet 企业服务商-get请求出错
    WxProviderRequestGet uint = 101172
    // WxProviderRequest 企业服务商-请求出错
    WxProviderRequest uint = 101173

    // 物流,取值范围101300-101499

    // LogisticsAMAliParam 阿里云市场阿里-参数错误
    LogisticsAMAliParam uint = 101300
    // LogisticsAMAliRequestPost 阿里云市场阿里-post请求出错
    LogisticsAMAliRequestPost uint = 101301
    // LogisticsAMAliRequestGet 阿里云市场阿里-get请求出错
    LogisticsAMAliRequestGet uint = 101302
    // LogisticsAMAliRequest 阿里云市场阿里-请求出错
    LogisticsAMAliRequest uint = 101303
    // LogisticsKd100Param 快递100-参数错误
    LogisticsKd100Param uint = 101310
    // LogisticsKd100RequestPost 快递100-post请求出错
    LogisticsKd100RequestPost uint = 101311
    // LogisticsKd100RequestGet 快递100-get请求出错
    LogisticsKd100RequestGet uint = 101312
    // LogisticsKd100Request 快递100-请求出错
    LogisticsKd100Request uint = 101313
    // LogisticsKdBirdParam 快递鸟-参数错误
    LogisticsKdBirdParam uint = 101320
    // LogisticsKdBirdRequestPost 快递鸟-post请求出错
    LogisticsKdBirdRequestPost uint = 101321
    // LogisticsKdBirdRequestGet 快递鸟-get请求出错
    LogisticsKdBirdRequestGet uint = 101322
    // LogisticsKdBirdRequest 快递鸟-请求出错
    LogisticsKdBirdRequest uint = 101323
    // LogisticsTaoBaoParam 淘宝-参数错误
    LogisticsTaoBaoParam uint = 101330
    // LogisticsTaoBaoRequestPost 淘宝-post请求出错
    LogisticsTaoBaoRequestPost uint = 101331
    // LogisticsTaoBaoRequestGet 淘宝-get请求出错
    LogisticsTaoBaoRequestGet uint = 101332
    // LogisticsTaoBaoRequest 淘宝-请求出错
    LogisticsTaoBaoRequest uint = 101333

    // 即时通讯,取值范围101500-101699

    // IMTencentParam 腾讯-参数错误
    IMTencentParam uint = 101500
    // IMTencentSign 腾讯-签名出错
    IMTencentSign uint = 101501
    // IMTencentRequestPost 腾讯-post请求出错
    IMTencentRequestPost uint = 101502
    // IMTencentRequestGet 腾讯-get请求出错
    IMTencentRequestGet uint = 101503
    // IMTencentRequest 腾讯-请求出错
    IMTencentRequest uint = 101504

    // 货币,取值范围101700-101899

    // CurrencyAMJiSuParam 阿里云市场极速-参数错误
    CurrencyAMJiSuParam uint = 101700
    // CurrencyAMJiSuRequestPost 阿里云市场极速-post请求出错
    CurrencyAMJiSuRequestPost uint = 101701
    // CurrencyAMJiSuRequestGet 阿里云市场极速-get请求出错
    CurrencyAMJiSuRequestGet uint = 101702
    // CurrencyAMJiSuRequest 阿里云市场极速-请求出错
    CurrencyAMJiSuRequest uint = 101703
    // CurrencyAMYiYuanParam 阿里云市场易圆-参数错误
    CurrencyAMYiYuanParam uint = 101710
    // CurrencyAMYiYuanRequestPost 阿里云市场易圆-post请求出错
    CurrencyAMYiYuanRequestPost uint = 101711
    // CurrencyAMYiYuanRequestGet 阿里云市场易圆-get请求出错
    CurrencyAMYiYuanRequestGet uint = 101712
    // CurrencyAMYiYuanRequest 阿里云市场易圆-请求出错
    CurrencyAMYiYuanRequest uint = 101713

    // 支付宝,取值范围101900-102099

    // AliPayParam 支付宝-参数错误
    AliPayParam uint = 101900
    // AliPayRequestPost 支付宝-post请求出错
    AliPayRequestPost uint = 101901
    // AliPayRequestGet 支付宝-get请求出错
    AliPayRequestGet uint = 101902
    // AliPayRequest 支付宝-请求出错
    AliPayRequest uint = 101903
    // AliPaySign 支付宝-签名出错
    AliPaySign uint = 101904
    // AliPayAuthParam 支付宝-授权参数错误
    AliPayAuthParam uint = 101950
    // AliPayFundParam 支付宝-资金参数错误
    AliPayFundParam uint = 101951
    // AliPayLifeParam 支付宝-生活号参数错误
    AliPayLifeParam uint = 101952
    // AliPayMarketParam 支付宝-店铺参数错误
    AliPayMarketParam uint = 101953
    // AliPayMaterialParam 支付宝-物料参数错误
    AliPayMaterialParam uint = 101954
    // AliPayTradeParam 支付宝-支付参数错误
    AliPayTradeParam uint = 101955

    // 腾讯云,取值范围102100-102299

    // QCloudParam 腾讯云-参数错误
    QCloudParam uint = 102100
    // QCloudRequestPost 腾讯云-post请求出错
    QCloudRequestPost uint = 102101
    // QCloudRequestGet 腾讯云-get请求出错
    QCloudRequestGet uint = 102102
    // QCloudRequest 腾讯云-请求出错
    QCloudRequest uint = 102103
    // QCloudCosParam 腾讯云-对象存储参数错误
    QCloudCosParam uint = 102150

    // 阿里云开放平台,取值范围102300-102499

    // AliOpenParam 阿里云开放平台-参数错误
    AliOpenParam uint = 102300
    // AliOpenRequestPost 阿里云开放平台-post请求出错
    AliOpenRequestPost uint = 102301
    // AliOpenRequestGet 阿里云开放平台-get请求出错
    AliOpenRequestGet uint = 102302
    // AliOpenRequest 阿里云开放平台-请求出错
    AliOpenRequest uint = 102303

    // 阿里云OSS,取值范围102500-102699

    // AliOssParam 阿里云OSS-参数错误
    AliOssParam uint = 102500
    // AliOssRequestPost 阿里云OSS-post请求出错
    AliOssRequestPost uint = 102501
    // AliOssRequestGet 阿里云OSS-get请求出错
    AliOssRequestGet uint = 102502
    // AliOssRequest 阿里云OSS-请求出错
    AliOssRequest uint = 102503

    // 七牛云,取值范围102700-102899

    // QiNiuKodoParam 七牛云图片-参数错误
    QiNiuKodoParam uint = 102700
    // QiNiuKodoRequestPost 七牛云图片-post请求出错
    QiNiuKodoRequestPost uint = 102701
    // QiNiuKodoRequestGet 七牛云图片-get请求出错
    QiNiuKodoRequestGet uint = 102702
    // QiNiuKodoRequest 七牛云图片-请求出错
    QiNiuKodoRequest uint = 102703

    // 钉钉,取值范围102900-103099

    // DingTalkParam 钉钉-参数错误
    DingTalkParam uint = 102900
    // DingTalkRequestPost 钉钉-post请求出错
    DingTalkRequestPost uint = 102901
    // DingTalkRequestGet 钉钉-get请求出错
    DingTalkRequestGet uint = 102902
    // DingTalkRequest 钉钉-请求出错
    DingTalkRequest uint = 102903
    // DingTalkCorpParam 企业号-参数错误
    DingTalkCorpParam uint = 102930
    // DingTalkCorpRequestPost 企业号-post请求出错
    DingTalkCorpRequestPost uint = 102931
    // DingTalkCorpRequestGet 企业号-get请求出错
    DingTalkCorpRequestGet uint = 102932
    // DingTalkCorpRequest 企业号-请求出错
    DingTalkCorpRequest uint = 102933
    // DingTalkProviderParam 服务商-参数错误
    DingTalkProviderParam uint = 102940
    // DingTalkProviderRequestPost 服务商-post请求出错
    DingTalkProviderRequestPost uint = 102941
    // DingTalkProviderRequestGet 服务商-get请求出错
    DingTalkProviderRequestGet uint = 102942
    // DingTalkProviderRequest 服务商-请求出错
    DingTalkProviderRequest uint = 102943

    // 物联网,取值范围103100-103299

    // IotAliYunParam 阿里云-参数错误
    IotAliYunParam uint = 103100
    // IotAliYunRequestPost 阿里云-post请求出错
    IotAliYunRequestPost uint = 103101
    // IotAliYunRequestGet 阿里云-get请求出错
    IotAliYunRequestGet uint = 103102
    // IotAliYunRequest 阿里云-请求出错
    IotAliYunRequest uint = 103100
    // IotBaiDuParam 百度-参数错误
    IotBaiDuParam uint = 103110
    // IotBaiDuRequestPost 百度-post请求出错
    IotBaiDuRequestPost uint = 103111
    // IotBaiDuRequestGet 百度-get请求出错
    IotBaiDuRequestGet uint = 103112
    // IotBaiDuRequest 百度-请求出错
    IotBaiDuRequest uint = 103113
    // IotTencentParam 腾讯-参数错误
    IotTencentParam uint = 103120
    // IotTencentRequestPost 腾讯-post请求出错
    IotTencentRequestPost uint = 103121
    // IotTencentRequestGet 腾讯-get请求出错
    IotTencentRequestGet uint = 103122
    // IotTencentRequest 腾讯-请求出错
    IotTencentRequest uint = 103123

    // 消息推送,取值范围103300-103499

    // PushAliYunParam 阿里云-参数错误
    PushAliYunParam uint = 103300
    // PushAliYunRequestPost 阿里云-post请求出错
    PushAliYunRequestPost uint = 103301
    // PushAliYunRequestGet 阿里云-get请求出错
    PushAliYunRequestGet uint = 103302
    // PushAliYunRequest 阿里云-请求出错
    PushAliYunRequest uint = 103303
    // PushBaiDuParam 百度-参数错误
    PushBaiDuParam uint = 103310
    // PushBaiDuRequestPost 百度-post请求出错
    PushBaiDuRequestPost uint = 103311
    // PushBaiDuRequestGet 百度-get请求出错
    PushBaiDuRequestGet uint = 103312
    // PushBaiDuRequest 百度-请求出错
    PushBaiDuRequest uint = 103313
    // PushXinGeParam 信鸽-参数错误
    PushXinGeParam uint = 103320
    // PushXinGeRequestPost 信鸽-post请求出错
    PushXinGeRequestPost uint = 103321
    // PushXinGeRequestGet 信鸽-get请求出错
    PushXinGeRequestGet uint = 103322
    // PushXinGeRequest 信鸽-请求出错
    PushXinGeRequest uint = 103323
    // PushJPushParam 极光-参数错误
    PushJPushParam uint = 103330
    // PushJPushRequestPost 极光-post请求出错
    PushJPushRequestPost uint = 103331
    // PushJPushRequestGet 极光-get请求出错
    PushJPushRequestGet uint = 103332
    // PushJPushRequest 极光-请求出错
    PushJPushRequest uint = 103333

    // 消息队列,取值范围103500-103699

    // MQParam 通用-参数错误
    MQParam uint = 103500
    // MQRedisParam redis-参数错误
    MQRedisParam uint = 103510
    // MQRedisConnect redis-连接错误
    MQRedisConnect uint = 103511
    // MQRedisProducer redis-生产者错误
    MQRedisProducer uint = 103512
    // MQRedisConsumer redis-消费者错误
    MQRedisConsumer uint = 103513
    // MQRabbitParam rabbit-参数错误
    MQRabbitParam uint = 103520
    // MQRabbitConnect rabbit-连接错误
    MQRabbitConnect uint = 103521
    // MQRabbitProducer rabbit-生产者错误
    MQRabbitProducer uint = 103522
    // MQRabbitConsumer rabbit-消费者错误
    MQRabbitConsumer uint = 103523

    // 协议,取值范围103700-103899

    // ProtocolPacket 协议解析-打包错误
    ProtocolPacket uint = 103700
    // ProtocolUnPacket 协议解析-解包错误
    ProtocolUnPacket uint = 103701
)
