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

// 添加卡券配置
type singleCardAddConfig struct {
    BaseWxSingle
    appId     string                       // 应用ID
    cardList  map[string]map[string]string // 卡券列表
    timestamp int                          // 时间戳
    nonceStr  string                       // 随机字符串
    needJs    bool                         // JS签名标识,true:需要JS签名 false:不需要JS签名
}

// 设置卡券列表,每个卡券信息包含card_id,code,openid,outer_str,fixed_begintimestamp
func (cec *singleCardAddConfig) SetCardList(cardList []map[string]string) {
    for _, v := range cardList {
        eCard := v
        delete(eCard, "card_id")
        cec.cardList[v["card_id"]] = eCard
    }
}

func (cec *singleCardAddConfig) SetNeedJs(needJs bool) {
    cec.needJs = needJs
}

func (cec *singleCardAddConfig) checkData() {
    if len(cec.cardList) == 0 {
        panic(mperr.NewWx(errorcode.WxParam, "卡券列表不能为空", nil))
    }
}

func (cec *singleCardAddConfig) GetResult(getType string) map[string]interface{} {
    cec.checkData()
    ticket := NewUtilWx().GetSingleCache(cec.appId, getType)
    timeStr := strconv.Itoa(cec.timestamp)
    result := make(map[string]interface{})

    cards := make([]map[string]string, 0)
    for k, v := range cec.cardList {
        cardInfo := v
        signData := make([]string, 0)
        signData = append(signData, ticket, timeStr, k, v["code"], v["openid"], cec.nonceStr)
        sort.Strings(signData)
        cardInfo["signature"] = mpf.HashSha1(strings.Join(signData, ""), "")
        cardInfo["timestamp"] = timeStr
        cardInfo["nonce_str"] = cec.nonceStr
        eCard := make(map[string]string)
        eCard["cardId"] = k
        eCard["cardExt"] = mpf.JsonMarshal(cardInfo)
        cards = append(cards, eCard)
    }
    result["card_list"] = cards

    if cec.needJs {
        jc := NewSingleJsConfig(cec.appId)
        jc.SetTimestamp(cec.timestamp)
        jc.SetNonceStr(cec.nonceStr)
        result["js_config"] = jc.GetResult(getType)
    }

    return result
}

func NewSingleCardAddConfig(appId string) *singleCardAddConfig {
    cec := &singleCardAddConfig{NewBaseWxSingle(), "", make(map[string]map[string]string), 0, "", false}
    cec.appId = appId
    cec.timestamp = time.Now().Second()
    cec.nonceStr = mpf.ToolCreateNonceStr(32, "numlower")
    cec.needJs = false
    return cec
}
