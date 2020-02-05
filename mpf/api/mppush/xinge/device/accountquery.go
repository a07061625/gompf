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

// 账号-设备绑定查询
type accountQuery struct {
    mppush.BaseXinGe
    operatorType    int    // 操作类型
    servicePlatform string // 平台类型
}

func (aq *accountQuery) SetOperatorType(operatorType int) {
    if (operatorType > 0) && (operatorType <= 2) {
        aq.operatorType = operatorType
    } else {
        panic(mperr.NewPushXinGe(errorcode.PushXinGeParam, "操作类型不合法", nil))
    }
}

func (aq *accountQuery) SetServicePlatform(servicePlatform string) {
    switch servicePlatform {
    case mppush.XinGePlatformTypeAndroid:
        aq.servicePlatform = servicePlatform
    case mppush.XinGePlatformTypeIOS:
        aq.servicePlatform = servicePlatform
    default:
        panic(mperr.NewPushXinGe(errorcode.PushXinGeParam, "平台类型不合法", nil))
    }
}

func (aq *accountQuery) SetAccountList(accountList []map[string]interface{}) {
    if len(accountList) > 0 {
        aq.ExtendData["account_list"] = accountList
    } else {
        panic(mperr.NewPushXinGe(errorcode.PushXinGeParam, "账号列表不合法", nil))
    }
}

func (aq *accountQuery) SetTokenList(tokenList []string) {
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
        aq.ExtendData["token_list"] = tokens
    }
}

func (aq *accountQuery) SetOpType(opType string) {
    switch opType {
    case "qq":
        aq.ExtendData["op_type"] = opType
    case "rtx":
        aq.ExtendData["op_type"] = opType
    case "email":
        aq.ExtendData["op_type"] = opType
    case "other":
        aq.ExtendData["op_type"] = opType
    default:
        panic(mperr.NewPushXinGe(errorcode.PushXinGeParam, "操作人员类型不合法", nil))
    }
}

func (aq *accountQuery) SetOpId(opId string) {
    if len(opId) > 0 {
        aq.ExtendData["op_id"] = opId
    } else {
        panic(mperr.NewPushXinGe(errorcode.PushXinGeParam, "接口操作人员id不合法", nil))
    }
}

func (aq *accountQuery) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if aq.operatorType <= 0 {
        panic(mperr.NewPushXinGe(errorcode.PushXinGeParam, "操作类型不能为空", nil))
    }
    if len(aq.servicePlatform) == 0 {
        panic(mperr.NewPushXinGe(errorcode.PushXinGeParam, "平台类型不能为空", nil))
    }
    aq.ExtendData["operator_type"] = aq.operatorType
    aq.ExtendData["platform"] = aq.servicePlatform

    return aq.GetRequest()
}

func NewAccountQuery(platform string) *accountQuery {
    aq := &accountQuery{mppush.NewBaseXinGe(platform), 0, ""}
    aq.ServiceUri = "device/account/query"
    return aq
}
