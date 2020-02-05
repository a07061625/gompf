/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/5 0005
 * Time: 0:18
 */
package wx

import (
    "sort"
    "strconv"
    "strings"
    "time"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
)

// 拉取卡券配置
type singleCardConfig struct {
    BaseWxSingle
    appId     string
    timestamp int
    nonceStr  string
    signType  string // 签名方式
    shopId    string // 门店ID
    cardId    string // 卡券ID
    cardType  string // 卡券类型
    needJs    bool   // JS签名标识,true:需要JS签名 false:不需要JS签名
}

func (cc *singleCardConfig) SetShopId(shopId string) {
    if len(shopId) <= 24 {
        cc.ReqData["shopId"] = shopId
    } else {
        panic(mperr.NewWx(errorcode.WxParam, "门店ID不合法", nil))
    }
}

func (cc *singleCardConfig) SetCardId(cardId string) {
    if len(cardId) <= 32 {
        cc.ReqData["cardId"] = cardId
    } else {
        panic(mperr.NewWx(errorcode.WxParam, "卡券ID不合法", nil))
    }
}

func (cc *singleCardConfig) SetCardType(cardType string) {
    if len(cardType) <= 24 {
        cc.ReqData["cardType"] = cardType
    } else {
        panic(mperr.NewWx(errorcode.WxParam, "卡券类型不合法", nil))
    }
}

func (cc *singleCardConfig) SetNeedJs(needJs bool) {
    cc.needJs = needJs
}

func (cc *singleCardConfig) checkData() {
    cc.ReqData["timestamp"] = strconv.Itoa(cc.timestamp)
}

func (cc *singleCardConfig) GetResult(getType string) map[string]interface{} {
    cc.checkData()
    ticket := NewUtilWx().GetSingleCache(cc.appId, getType)
    signData := make([]string, 0)
    signData = append(signData, ticket, cc.appId, cc.ReqData["shopId"], cc.ReqData["timestamp"], cc.ReqData["nonceStr"], cc.ReqData["cardId"], cc.ReqData["cardType"])
    sort.Strings(signData)
    cc.ReqData["cardSign"] = mpf.HashSha1(strings.Join(signData, ""), "")

    result := make(map[string]interface{})
    for k, v := range cc.ReqData {
        result[k] = v
    }

    if cc.needJs {
        jc := NewSingleJsConfig(cc.appId)
        jc.SetTimestamp(cc.timestamp)
        jc.SetNonceStr(cc.ReqData["nonceStr"])
        result["jsConfig"] = jc.GetResult(getType)
    }

    return result
}

func NewSingleCardConfig(appId string) *singleCardConfig {
    cc := &singleCardConfig{NewBaseWxSingle(), "", 0, "", "", "", "", "", false}
    cc.appId = appId
    cc.timestamp = time.Now().Second()
    cc.needJs = false
    cc.ReqData["nonceStr"] = mpf.ToolCreateNonceStr(32, "numlower")
    cc.ReqData["signType"] = "SHA1"
    cc.ReqData["shopId"] = ""
    cc.ReqData["cardId"] = ""
    cc.ReqData["cardType"] = ""
    return cc
}
