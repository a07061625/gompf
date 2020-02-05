/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/15 0015
 * Time: 0:41
 */
package single

const (
    CompanyBankCodeCMB  = "1001"
    CompanyBankCodeICBC = "1002"
    CompanyBankCodeCCB  = "1003"
    CompanyBankCodeSPDB = "1004"
    CompanyBankCodeABC  = "1005"
    CompanyBankCodeCMBC = "1006"
    CompanyBankCodeCIB  = "1009"
    CompanyBankCodePAB  = "1010"
    CompanyBankCodeBCM  = "1020"
    CompanyBankCodeZXB  = "1021"
    CompanyBankCodeCEB  = "1022"
    CompanyBankCodeHXB  = "1025"
    CompanyBankCodeBOC  = "1026"
    CompanyBankCodeGDB  = "1027"
    CompanyBankCodeBOB  = "1032"
    CompanyBankCodeNBB  = "1056"
    CompanyBankCodePSBC = "1066"

    TradeTypeJsApi     = "JSAPI"
    TradeTypeNative    = "NATIVE"
    TradeTypeApp       = "APP"
    TradeTypeMobileWeb = "MWEB"
)

var (
    CompanyBankCodes map[string]string
    TradeTypes       map[string]string
)

func init() {
    CompanyBankCodes = make(map[string]string)
    CompanyBankCodes[CompanyBankCodeCMB] = "招商银行"
    CompanyBankCodes[CompanyBankCodeICBC] = "工商银行"
    CompanyBankCodes[CompanyBankCodeCCB] = "建设银行"
    CompanyBankCodes[CompanyBankCodeSPDB] = "浦发银行"
    CompanyBankCodes[CompanyBankCodeABC] = "农业银行"
    CompanyBankCodes[CompanyBankCodeCMBC] = "民生银行"
    CompanyBankCodes[CompanyBankCodeCIB] = "兴业银行"
    CompanyBankCodes[CompanyBankCodePAB] = "平安银行"
    CompanyBankCodes[CompanyBankCodeBCM] = "交通银行"
    CompanyBankCodes[CompanyBankCodeZXB] = "中信银行"
    CompanyBankCodes[CompanyBankCodeCEB] = "光大银行"
    CompanyBankCodes[CompanyBankCodeHXB] = "华夏银行"
    CompanyBankCodes[CompanyBankCodeBOC] = "中国银行"
    CompanyBankCodes[CompanyBankCodeGDB] = "广发银行"
    CompanyBankCodes[CompanyBankCodeBOB] = "北京银行"
    CompanyBankCodes[CompanyBankCodeNBB] = "宁波银行"
    CompanyBankCodes[CompanyBankCodePSBC] = "邮储银行"

    TradeTypes = make(map[string]string)
    TradeTypes[TradeTypeJsApi] = "jsapi"
    TradeTypes[TradeTypeNative] = "扫码"
    TradeTypes[TradeTypeApp] = "app"
    TradeTypes[TradeTypeMobileWeb] = "h5"
}
