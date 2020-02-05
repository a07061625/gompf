/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2019/12/31 0031
 * Time: 0:38
 */
package wx

const (
    // 公众号或小程序缓存类型
    SingleCacheTypeAccessToken     = "single_accesstoken"
    SingleCacheTypeJsTicket        = "single_jsticket"
    SingleCacheTypeCardTicket      = "single_cardticket"
    SingleCacheTypeOpenAccessToken = "open_accesstoken"
    SingleCacheTypeOpenJsTicket    = "open_jsticket"
    SingleCacheTypeOpenCardTicket  = "open_cardticket"

    // 企业号缓存类型
    CorpCacheTypeAccessToken         = "corp_accesstoken"
    CorpCacheTypeJsTicket            = "corp_jsticket"
    CorpCacheTypeProviderAccessToken = "provider_accesstoken"
    CorpCacheTypeProviderJsTicket    = "provider_jsticket"

    AccountMerchantTypeSelf = "self" // 自身
    AccountMerchantTypeSub  = "sub"  // 子商户
)
