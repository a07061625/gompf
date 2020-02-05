/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/3 0003
 * Time: 15:36
 */
package device

import (
    "github.com/a07061625/gompf/mpf/api/mppush"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 账号绑定与解绑
type accountOperateBatch struct {
    mppush.BaseXinGe
    operatorType    int    // 操作类型
    servicePlatform string // 平台类型
}

func (aob *accountOperateBatch) SetOperatorType(operatorType int) {
    if (operatorType > 0) && (operatorType <= 5) {
        aob.operatorType = operatorType
    } else {
        panic(mperr.NewPushXinGe(errorcode.PushXinGeParam, "操作类型不合法", nil))
    }
}

func (aob *accountOperateBatch) SetServicePlatform(servicePlatform string) {
    switch servicePlatform {
    case mppush.XinGePlatformTypeAndroid:
        aob.servicePlatform = servicePlatform
    case mppush.XinGePlatformTypeIOS:
        aob.servicePlatform = servicePlatform
    default:
        panic(mperr.NewPushXinGe(errorcode.PushXinGeParam, "平台类型不合法", nil))
    }
}

func (aob *accountOperateBatch) SetAccountList(accountList []map[string]interface{}) {
    if len(accountList) > 0 {
        aob.ExtendData["account_list"] = accountList
    } else {
        panic(mperr.NewPushXinGe(errorcode.PushXinGeParam, "账号列表不合法", nil))
    }
}

func (aob *accountOperateBatch) SetTokenList(tokenList []string) {
    if len(tokenList) == 0 {
        panic(mperr.NewPushXinGe(errorcode.PushXinGeParam, "设备列表不合法", nil))
    }

    tokens := make([]string, 0)
    for _, v := range tokenList {
        if len(v) > 0 {
            tokens = append(tokens, v)
        }
    }
    if len(tokens) > 0 {
        aob.ExtendData["token_list"] = tokens
    }
}

func (aob *accountOperateBatch) SetTokenAccounts(tokenAccounts []map[string]interface{}) {
    if len(tokenAccounts) > 0 {
        aob.ExtendData["token_accounts"] = tokenAccounts
    } else {
        panic(mperr.NewPushXinGe(errorcode.PushXinGeParam, "设备账号列表不合法", nil))
    }
}

func (aob *accountOperateBatch) SetOpType(opType string) {
    switch opType {
    case "qq":
        aob.ExtendData["op_type"] = opType
    case "rtx":
        aob.ExtendData["op_type"] = opType
    case "email":
        aob.ExtendData["op_type"] = opType
    case "other":
        aob.ExtendData["op_type"] = opType
    default:
        panic(mperr.NewPushXinGe(errorcode.PushXinGeParam, "操作人员类型不合法", nil))
    }
}

func (aob *accountOperateBatch) SetOpId(opId string) {
    if len(opId) > 0 {
        aob.ExtendData["op_id"] = opId
    } else {
        panic(mperr.NewPushXinGe(errorcode.PushXinGeParam, "接口操作人员id不合法", nil))
    }
}

func (aob *accountOperateBatch) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if aob.operatorType <= 0 {
        panic(mperr.NewPushXinGe(errorcode.PushXinGeParam, "操作类型不能为空", nil))
    }
    if len(aob.servicePlatform) == 0 {
        panic(mperr.NewPushXinGe(errorcode.PushXinGeParam, "平台类型不能为空", nil))
    }
    aob.ExtendData["operator_type"] = aob.operatorType
    aob.ExtendData["platform"] = aob.servicePlatform

    return aob.GetRequest()
}

func NewAccountOperateBatch(platform string) *accountOperateBatch {
    aob := &accountOperateBatch{mppush.NewBaseXinGe(platform), 0, ""}
    aob.ServiceUri = "device/account/batchoperate"
    return aob
}
