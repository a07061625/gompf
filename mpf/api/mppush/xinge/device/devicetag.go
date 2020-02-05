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

// 标签操作
type deviceTag struct {
    mppush.BaseXinGe
    operatorType    int    // 操作类型
    servicePlatform string // 平台类型
    seq             int64  //请求ID
}

func (dt *deviceTag) SetOperatorType(operatorType int) {
    if (operatorType > 0) && (operatorType <= 10) {
        dt.operatorType = operatorType
    } else {
        panic(mperr.NewPushXinGe(errorcode.PushXinGeParam, "操作类型不合法", nil))
    }
}

func (dt *deviceTag) SetServicePlatform(servicePlatform string) {
    switch servicePlatform {
    case mppush.XinGePlatformTypeAndroid:
        dt.servicePlatform = servicePlatform
    case mppush.XinGePlatformTypeIOS:
        dt.servicePlatform = servicePlatform
    default:
        panic(mperr.NewPushXinGe(errorcode.PushXinGeParam, "平台类型不合法", nil))
    }
}

func (dt *deviceTag) SetSeq(seq int64) {
    if seq > 0 {
        dt.seq = seq
    } else {
        panic(mperr.NewPushXinGe(errorcode.PushXinGeParam, "请求ID不合法", nil))
    }
}

func (dt *deviceTag) SetTokenList(tokenList []string) {
    if len(tokenList) == 0 {
        panic(mperr.NewPushXinGe(errorcode.PushXinGeParam, "设备列表不合法", nil))
    } else if len(tokenList) > 20 {
        panic(mperr.NewPushXinGe(errorcode.PushXinGeParam, "设备列表不能超过20个", nil))
    }

    tokens := make([]string, 0)
    for _, v := range tokenList {
        if len(v) > 0 {
            tokens = append(tokens, v)
        }
    }
    if len(tokens) > 0 {
        dt.ExtendData["token_list"] = tokens
    }
}

func (dt *deviceTag) SetTagList(tagList []string) {
    if len(tagList) == 0 {
        panic(mperr.NewPushXinGe(errorcode.PushXinGeParam, "标签列表不合法", nil))
    } else if len(tagList) > 20 {
        panic(mperr.NewPushXinGe(errorcode.PushXinGeParam, "标签列表不能超过20个", nil))
    }

    tags := make([]string, 0)
    for _, v := range tagList {
        if len(v) > 0 {
            tags = append(tags, v)
        }
    }
    if len(tags) > 0 {
        dt.ExtendData["tag_list"] = tags
    }
}

func (dt *deviceTag) SetTagTokenList(tagTokenList map[string]string) {
    if len(tagTokenList) == 0 {
        panic(mperr.NewPushXinGe(errorcode.PushXinGeParam, "标签设备对应列表不能为空", nil))
    } else if len(tagTokenList) > 20 {
        panic(mperr.NewPushXinGe(errorcode.PushXinGeParam, "标签设备对应列表不能超过20个", nil))
    }

    tagTokens := make([][]string, 0)
    for k, v := range tagTokenList {
        eMap := make([]string, 0)
        eMap = append(eMap, k)
        eMap = append(eMap, v)
        tagTokens = append(tagTokens, eMap)
    }
    dt.ExtendData["tag_token_list"] = tagTokens
}

func (dt *deviceTag) SetOpType(opType string) {
    switch opType {
    case "qq":
        dt.ExtendData["op_type"] = opType
    case "rtx":
        dt.ExtendData["op_type"] = opType
    case "email":
        dt.ExtendData["op_type"] = opType
    case "other":
        dt.ExtendData["op_type"] = opType
    default:
        panic(mperr.NewPushXinGe(errorcode.PushXinGeParam, "操作人员类型不合法", nil))
    }
}

func (dt *deviceTag) SetOpId(opId string) {
    if len(opId) > 0 {
        dt.ExtendData["op_id"] = opId
    } else {
        panic(mperr.NewPushXinGe(errorcode.PushXinGeParam, "接口操作人员id不合法", nil))
    }
}

func (dt *deviceTag) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if dt.operatorType <= 0 {
        panic(mperr.NewPushXinGe(errorcode.PushXinGeParam, "操作类型不能为空", nil))
    }
    if len(dt.servicePlatform) == 0 {
        panic(mperr.NewPushXinGe(errorcode.PushXinGeParam, "平台类型不能为空", nil))
    }
    if dt.seq <= 0 {
        panic(mperr.NewPushXinGe(errorcode.PushXinGeParam, "请求ID不能为空", nil))
    }
    dt.ExtendData["operator_type"] = dt.operatorType
    dt.ExtendData["platform"] = dt.servicePlatform
    dt.ExtendData["seq"] = dt.seq

    return dt.GetRequest()
}

func NewDeviceTag(platform string) *deviceTag {
    dt := &deviceTag{mppush.NewBaseXinGe(platform), 0, "", 0}
    dt.ServiceUri = "device/tag"
    return dt
}
