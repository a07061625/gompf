/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/19 0019
 * Time: 22:28
 */
package group

import (
    "regexp"

    "github.com/a07061625/gompf/mpf/api/alipay"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 用户分组修改接口
type groupModify struct {
    alipay.BaseAliPay
    groupId   string                   // 分组ID
    groupName string                   // 分组名称
    labelRule []map[string]interface{} // 标签规则
}

func (gm *groupModify) SetGroupId(groupId string) {
    match, _ := regexp.MatchString(`^[0-9]{1,10}$`, groupId)
    if match {
        gm.groupId = groupId
    } else {
        panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "分组ID不合法", nil))
    }
}

func (gm *groupModify) SetGroupName(groupName string) {
    if (len(groupName) > 0) && (len(groupName) <= 30) {
        gm.groupName = groupName
    } else {
        panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "分组名称不合法", nil))
    }
}

func (gm *groupModify) SetLabelRule(labelRule []map[string]interface{}) {
    if len(labelRule) > 0 {
        gm.labelRule = labelRule
    } else {
        panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "标签规则不合法", nil))
    }
}

func (gm *groupModify) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(gm.groupId) == 0 {
        panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "分组ID不能为空", nil))
    }
    if (len(gm.groupName) == 0) && (len(gm.labelRule) == 0) {
        panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "分组名称和标签规则不能都为空", nil))
    }
    gm.BizContent["group_id"] = gm.groupId
    if len(gm.groupName) > 0 {
        gm.BizContent["name"] = gm.groupName
    }
    if len(gm.labelRule) > 0 {
        gm.BizContent["label_rule"] = gm.labelRule
    }

    return gm.GetRequest()
}

func NewGroupModify(appId string) *groupModify {
    gm := &groupModify{alipay.NewBase(appId), "", "", make([]map[string]interface{}, 0)}
    gm.SetMethod("alipay.open.public.group.modify")
    return gm
}
