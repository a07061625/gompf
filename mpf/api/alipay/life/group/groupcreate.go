/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/19 0019
 * Time: 22:17
 */
package group

import (
    "github.com/a07061625/gompf/mpf/api/alipay"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 用户分组创建接口
type groupCreate struct {
    alipay.BaseAliPay
    groupName string                   // 分组名称
    labelRule []map[string]interface{} // 标签规则
}

func (gc *groupCreate) SetGroupName(groupName string) {
    if (len(groupName) > 0) && (len(groupName) <= 30) {
        gc.groupName = groupName
    } else {
        panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "分组名称不合法", nil))
    }
}

func (gc *groupCreate) SetLabelRule(labelRule []map[string]interface{}) {
    if len(labelRule) > 0 {
        gc.labelRule = labelRule
    } else {
        panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "标签规则不合法", nil))
    }
}

func (gc *groupCreate) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(gc.groupName) == 0 {
        panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "分组名称不能为空", nil))
    }
    if len(gc.labelRule) == 0 {
        panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "标签规则不能为空", nil))
    }
    gc.BizContent["name"] = gc.groupName
    gc.BizContent["label_rule"] = gc.labelRule

    return gc.GetRequest()
}

func NewGroupCreate(appId string) *groupCreate {
    gc := &groupCreate{alipay.NewBase(appId), "", make([]map[string]interface{}, 0)}
    gc.SetMethod("alipay.open.public.group.create")
    return gc
}
