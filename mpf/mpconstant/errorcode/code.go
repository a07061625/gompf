/**
 * 错误码
 * User: 姜伟
 * Date: 2020/1/17 0017
 * Time: 19:24
 */
package errorcode

const (
    // 公共基础,取值范围10000-99999
    CommonBaseSuccess uint = 0     // 成功
    CommonBaseMin     uint = 10000 // 最小错误码
    CommonBaseInner   uint = 10000 // 内部错误
    CommonBaseServer  uint = 10001 // 服务错误

    // 公共请求
    CommonRequestFail          uint = 30001 // 请求失败
    CommonRequestParam         uint = 30002 // 请求参数错误
    CommonRequestTimeout       uint = 30003 // 请求超时
    CommonRequestResourceEmpty uint = 30004 // 请求资源不存在
    CommonRequestFormat        uint = 30005 // 请求格式错误
    CommonRequestMethod        uint = 30006 // 请求方式错误

    // 公共响应
    CommonResponseEmpty  uint = 30050 // 响应为空
    CommonResponseFormat uint = 30051 // 响应格式错误

    // 公共json
    CommonJsonMarshal   uint = 30100 // 压缩错误
    CommonJsonUnmarshal uint = 30101 // 解析错误

    // 公共校验器
    CommonValidatorSign uint = 30150 // 校验器-签名错误
    CommonValidatorRule uint = 30151 // 校验器-规则错误

    // 缓存,取值范围100000-100099
    CacheRedisConnect    uint = 100000 // redis-连接错误
    CacheRedisOperate    uint = 100001 // redis-操作错误
    CacheMemCacheConnect uint = 100010 // memcache-连接错误
    CacheMemCacheOperate uint = 100011 // memcache-操作错误

    // 配置,取值范围100100-100199
    ConfigViperGet uint = 100100 // viper-获取错误

    // 数据库,取值范围100200-100299
    DbMysqlConnect       uint = 100200 // mysql-连接错误
    DbMysqlOperateInsert uint = 100201 // mysql-插入操作错误
    DbMysqlOperateDelete uint = 100202 // mysql-删除操作错误
    DbMysqlOperateUpdate uint = 100203 // mysql-更新操作错误
    DbMysqlOperateSelect uint = 100204 // mysql-查询操作错误
    DbMonGoConnect       uint = 100210 // mongo-连接错误
    DbMonGoOperateInsert uint = 100211 // mongo-插入操作错误
    DbMonGoOperateDelete uint = 100212 // mongo-删除操作错误
    DbMonGoOperateUpdate uint = 100213 // mongo-更新操作错误
    DbMonGoOperateSelect uint = 100214 // mongo-查询操作错误

    // 日志,取值范围100300-100399
    LogZapConfig uint = 100300 // zap-配置错误

    // 短信,取值范围100400-100799
    SmsDaYuParam         uint = 100400 // 大鱼-参数错误
    SmsDaYuRequestPost   uint = 100401 // 大鱼-post请求出错
    SmsDaYuRequestGet    uint = 100402 // 大鱼-get请求出错
    SmsDaYuRequest       uint = 100403 // 大鱼-请求出错
    SmsAliYunParam       uint = 100410 // 阿里云-参数错误
    SmsAliYunRequestPost uint = 100411 // 阿里云-post请求出错
    SmsAliYunRequestGet  uint = 100412 // 阿里云-get请求出错
    SmsAliYunRequest     uint = 100413 // 阿里云-请求出错
    SmsYun253Param       uint = 100420 // 253云-参数错误
    SmsYun253RequestPost uint = 100421 // 253云-post请求出错
    SmsYun253RequestGet  uint = 100422 // 253云-get请求出错
    SmsYun253Request     uint = 100423 // 253云-请求出错

    // 地图,取值范围100800-100999
    MapTencentParam       uint = 100800 // 腾讯-参数错误
    MapTencentRequestPost uint = 100801 // 腾讯-post请求出错
    MapTencentRequestGet  uint = 100802 // 腾讯-get请求出错
    MapTencentRequest     uint = 100803 // 腾讯-请求出错
    MapBaiDuParam         uint = 100810 // 百度-参数错误
    MapBaiDuRequestPost   uint = 100811 // 百度-post请求出错
    MapBaiDuRequestGet    uint = 100812 // 百度-get请求出错
    MapBaiDuRequest       uint = 100813 // 百度-请求出错
    MapGaoDeParam         uint = 100820 // 高德-参数错误
    MapGaoDeRequestPost   uint = 100821 // 高德-post请求出错
    MapGaoDeRequestGet    uint = 100822 // 高德-get请求出错
    MapGaoDeRequest       uint = 100823 // 高德-请求出错

    // 打印,取值范围101000-101099
    PrintFeYinParam       uint = 101000 // 飞印-参数错误
    PrintFeYinRequestPost uint = 101001 // 飞印-post请求出错
    PrintFeYinRequestGet  uint = 101002 // 飞印-get请求出错
    PrintFeYinRequest     uint = 101003 // 飞印-请求出错

    // 微信,取值范围101100-101299
    WxParam               uint = 101100 // 微信-参数错误
    WxRequestPost         uint = 101101 // 微信-post请求出错
    WxRequestGet          uint = 101102 // 微信-get请求出错
    WxRequest             uint = 101103 // 微信-请求出错
    WxAccountParam        uint = 101130 // 公众号-参数错误
    WxAccountRequestPost  uint = 101131 // 公众号-post请求出错
    WxAccountRequestGet   uint = 101132 // 公众号-get请求出错
    WxAccountRequest      uint = 101133 // 公众号-请求出错
    WxCorpParam           uint = 101140 // 企业号-参数错误
    WxCorpRequestPost     uint = 101141 // 企业号-post请求出错
    WxCorpRequestGet      uint = 101142 // 企业号-get请求出错
    WxCorpRequest         uint = 101143 // 企业号-请求出错
    WxMiniParam           uint = 101150 // 小程序-参数错误
    WxMiniRequestPost     uint = 101151 // 小程序-post请求出错
    WxMiniRequestGet      uint = 101152 // 小程序-get请求出错
    WxMiniRequest         uint = 101153 // 小程序-请求出错
    WxOpenParam           uint = 101160 // 第三方开放平台-参数错误
    WxOpenRequestPost     uint = 101161 // 第三方开放平台-post请求出错
    WxOpenRequestGet      uint = 101162 // 第三方开放平台-get请求出错
    WxOpenRequest         uint = 101163 // 第三方开放平台-请求出错
    WxProviderParam       uint = 101170 // 企业服务商-参数错误
    WxProviderRequestPost uint = 101171 // 企业服务商-post请求出错
    WxProviderRequestGet  uint = 101172 // 企业服务商-get请求出错
    WxProviderRequest     uint = 101173 // 企业服务商-请求出错

    // 物流,取值范围101300-101499
    LogisticsAMAliParam        uint = 101300 // 阿里云市场阿里-参数错误
    LogisticsAMAliRequestPost  uint = 101301 // 阿里云市场阿里-post请求出错
    LogisticsAMAliRequestGet   uint = 101302 // 阿里云市场阿里-get请求出错
    LogisticsAMAliRequest      uint = 101303 // 阿里云市场阿里-请求出错
    LogisticsKd100Param        uint = 101310 // 快递100-参数错误
    LogisticsKd100RequestPost  uint = 101311 // 快递100-post请求出错
    LogisticsKd100RequestGet   uint = 101312 // 快递100-get请求出错
    LogisticsKd100Request      uint = 101313 // 快递100-请求出错
    LogisticsKdBirdParam       uint = 101320 // 快递鸟-参数错误
    LogisticsKdBirdRequestPost uint = 101321 // 快递鸟-post请求出错
    LogisticsKdBirdRequestGet  uint = 101322 // 快递鸟-get请求出错
    LogisticsKdBirdRequest     uint = 101323 // 快递鸟-请求出错
    LogisticsTaoBaoParam       uint = 101330 // 淘宝-参数错误
    LogisticsTaoBaoRequestPost uint = 101331 // 淘宝-post请求出错
    LogisticsTaoBaoRequestGet  uint = 101332 // 淘宝-get请求出错
    LogisticsTaoBaoRequest     uint = 101333 // 淘宝-请求出错

    // 即时通讯,取值范围101500-101699
    IMTencentParam       uint = 101500 // 腾讯-参数错误
    IMTencentSign        uint = 101501 // 腾讯-签名出错
    IMTencentRequestPost uint = 101502 // 腾讯-post请求出错
    IMTencentRequestGet  uint = 101503 // 腾讯-get请求出错
    IMTencentRequest     uint = 101504 // 腾讯-请求出错

    // 货币,取值范围101700-101899
    CurrencyAMJiSuParam         uint = 101700 // 阿里云市场极速-参数错误
    CurrencyAMJiSuRequestPost   uint = 101701 // 阿里云市场极速-post请求出错
    CurrencyAMJiSuRequestGet    uint = 101702 // 阿里云市场极速-get请求出错
    CurrencyAMJiSuRequest       uint = 101703 // 阿里云市场极速-请求出错
    CurrencyAMYiYuanParam       uint = 101710 // 阿里云市场易圆-参数错误
    CurrencyAMYiYuanRequestPost uint = 101711 // 阿里云市场易圆-post请求出错
    CurrencyAMYiYuanRequestGet  uint = 101712 // 阿里云市场易圆-get请求出错
    CurrencyAMYiYuanRequest     uint = 101713 // 阿里云市场易圆-请求出错

    // 支付宝,取值范围101900-102099
    AliPayParam         uint = 101900 // 支付宝-参数错误
    AliPayRequestPost   uint = 101901 // 支付宝-post请求出错
    AliPayRequestGet    uint = 101902 // 支付宝-get请求出错
    AliPayRequest       uint = 101903 // 支付宝-请求出错
    AliPaySign          uint = 101904 // 支付宝-签名出错
    AliPayAuthParam     uint = 101950 // 支付宝-授权参数错误
    AliPayFundParam     uint = 101951 // 支付宝-资金参数错误
    AliPayLifeParam     uint = 101952 // 支付宝-生活号参数错误
    AliPayMarketParam   uint = 101953 // 支付宝-店铺参数错误
    AliPayMaterialParam uint = 101954 // 支付宝-物料参数错误
    AliPayTradeParam    uint = 101955 // 支付宝-支付参数错误

    // 腾讯云,取值范围102100-102299
    QCloudParam       uint = 102100 // 腾讯云-参数错误
    QCloudRequestPost uint = 102101 // 腾讯云-post请求出错
    QCloudRequestGet  uint = 102102 // 腾讯云-get请求出错
    QCloudRequest     uint = 102103 // 腾讯云-请求出错
    QCloudCosParam    uint = 102150 // 腾讯云-对象存储参数错误

    // 阿里云开放平台,取值范围102300-102499
    AliOpenParam       uint = 102300 // 阿里云开放平台-参数错误
    AliOpenRequestPost uint = 102301 // 阿里云开放平台-post请求出错
    AliOpenRequestGet  uint = 102302 // 阿里云开放平台-get请求出错
    AliOpenRequest     uint = 102303 // 阿里云开放平台-请求出错

    // 阿里云OSS,取值范围102500-102699
    AliOssParam       uint = 102500 // 阿里云OSS-参数错误
    AliOssRequestPost uint = 102501 // 阿里云OSS-post请求出错
    AliOssRequestGet  uint = 102502 // 阿里云OSS-get请求出错
    AliOssRequest     uint = 102503 // 阿里云OSS-请求出错

    // 七牛云,取值范围102700-102899
    QiNiuKodoParam       uint = 102700 // 七牛云图片-参数错误
    QiNiuKodoRequestPost uint = 102701 // 七牛云图片-post请求出错
    QiNiuKodoRequestGet  uint = 102702 // 七牛云图片-get请求出错
    QiNiuKodoRequest     uint = 102703 // 七牛云图片-请求出错

    // 钉钉,取值范围102900-103099
    DingTalkParam               uint = 102900 // 钉钉-参数错误
    DingTalkRequestPost         uint = 102901 // 钉钉-post请求出错
    DingTalkRequestGet          uint = 102902 // 钉钉-get请求出错
    DingTalkRequest             uint = 102903 // 钉钉-请求出错
    DingTalkCorpParam           uint = 102930 // 企业号-参数错误
    DingTalkCorpRequestPost     uint = 102931 // 企业号-post请求出错
    DingTalkCorpRequestGet      uint = 102932 // 企业号-get请求出错
    DingTalkCorpRequest         uint = 102933 // 企业号-请求出错
    DingTalkProviderParam       uint = 102940 // 服务商-参数错误
    DingTalkProviderRequestPost uint = 102941 // 服务商-post请求出错
    DingTalkProviderRequestGet  uint = 102942 // 服务商-get请求出错
    DingTalkProviderRequest     uint = 102943 // 服务商-请求出错

    // 物联网,取值范围103100-103299
    IotAliYunParam        uint = 103100 // 阿里云-参数错误
    IotAliYunRequestPost  uint = 103101 // 阿里云-post请求出错
    IotAliYunRequestGet   uint = 103102 // 阿里云-get请求出错
    IotAliYunRequest      uint = 103100 // 阿里云-请求出错
    IotBaiDuParam         uint = 103110 // 百度-参数错误
    IotBaiDuRequestPost   uint = 103111 // 百度-post请求出错
    IotBaiDuRequestGet    uint = 103112 // 百度-get请求出错
    IotBaiDuRequest       uint = 103113 // 百度-请求出错
    IotTencentParam       uint = 103120 // 腾讯-参数错误
    IotTencentRequestPost uint = 103121 // 腾讯-post请求出错
    IotTencentRequestGet  uint = 103122 // 腾讯-get请求出错
    IotTencentRequest     uint = 103123 // 腾讯-请求出错

    // 消息推送,取值范围103300-103499
    PushAliYunParam       uint = 103300 // 阿里云-参数错误
    PushAliYunRequestPost uint = 103301 // 阿里云-post请求出错
    PushAliYunRequestGet  uint = 103302 // 阿里云-get请求出错
    PushAliYunRequest     uint = 103303 // 阿里云-请求出错
    PushBaiDuParam        uint = 103310 // 百度-参数错误
    PushBaiDuRequestPost  uint = 103311 // 百度-post请求出错
    PushBaiDuRequestGet   uint = 103312 // 百度-get请求出错
    PushBaiDuRequest      uint = 103313 // 百度-请求出错
    PushXinGeParam        uint = 103320 // 信鸽-参数错误
    PushXinGeRequestPost  uint = 103321 // 信鸽-post请求出错
    PushXinGeRequestGet   uint = 103322 // 信鸽-get请求出错
    PushXinGeRequest      uint = 103323 // 信鸽-请求出错
    PushJPushParam        uint = 103330 // 极光-参数错误
    PushJPushRequestPost  uint = 103331 // 极光-post请求出错
    PushJPushRequestGet   uint = 103332 // 极光-get请求出错
    PushJPushRequest      uint = 103333 // 极光-请求出错

    // 消息队列,取值范围103500-103699
    MQParam          uint = 103500 // 通用-参数错误
    MQRedisParam     uint = 103510 // redis-参数错误
    MQRedisConnect   uint = 103511 // redis-连接错误
    MQRedisProducer  uint = 103512 // redis-生产者错误
    MQRedisConsumer  uint = 103513 // redis-消费者错误
    MQRabbitParam    uint = 103520 // rabbit-参数错误
    MQRabbitConnect  uint = 103521 // rabbit-连接错误
    MQRabbitProducer uint = 103522 // rabbit-生产者错误
    MQRabbitConsumer uint = 103523 // rabbit-消费者错误
)
