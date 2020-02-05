/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/19 0019
 * Time: 9:32
 */
package alipay

const (
    UrlGateWay = "https://openapi.alipay.com/gateway.do"

    FundCurrencyAUD = "AUD"
    FundCurrencyNZD = "NZD"
    FundCurrencyTWD = "TWD"
    FundCurrencyUSD = "USD"
    FundCurrencyEUR = "EUR"
    FundCurrencyGBP = "GBP"
    FundCurrencyCNY = "CNY"
)

var (
    FundCurrencyList map[string]string
)

func init() {
    FundCurrencyList = make(map[string]string)
    FundCurrencyList[FundCurrencyAUD] = "澳元"
    FundCurrencyList[FundCurrencyNZD] = "新西兰元"
    FundCurrencyList[FundCurrencyTWD] = "台币"
    FundCurrencyList[FundCurrencyUSD] = "美元"
    FundCurrencyList[FundCurrencyEUR] = "欧元"
    FundCurrencyList[FundCurrencyGBP] = "英镑"
}
