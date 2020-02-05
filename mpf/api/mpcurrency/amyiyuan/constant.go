/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/17 0017
 * Time: 18:29
 */
package amyiyuan

const (
    BankTypeICBC = "ICBC"
    BankTypeBOC  = "BOC"
    BankTypeAB   = "ABCHINA"
    BankTypeCOMM = "BANKCOMM"
    BankTypeCCB  = "CCB"
    BankTypeCMB  = "CMBCHINA"
    BankTypeCEB  = "CEBBANK"
    BankTypeSPDB = "SPDB"
    BankTypeCIB  = "CIB"
    BankTypeZXB  = "ECITIC"
)

var (
    bankTypes map[string]string
)

func init() {
    bankTypes = make(map[string]string)
    bankTypes[BankTypeICBC] = "工商银行"
    bankTypes[BankTypeBOC] = "中国银行"
    bankTypes[BankTypeAB] = "农业银行"
    bankTypes[BankTypeCOMM] = "交通银行"
    bankTypes[BankTypeCCB] = "建设银行"
    bankTypes[BankTypeCMB] = "招商银行"
    bankTypes[BankTypeCEB] = "光大银行"
    bankTypes[BankTypeSPDB] = "浦发银行"
    bankTypes[BankTypeCIB] = "兴业银行"
    bankTypes[BankTypeZXB] = "中信银行"
}
