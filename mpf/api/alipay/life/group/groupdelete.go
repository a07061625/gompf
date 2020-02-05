/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/19 0019
 * Time: 22:24
 */
package group

import (
    "regexp"

    "github.com/a07061625/gompf/mpf/api/alipay"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 用户分组删除接口
type groupDelete struct {
    alipay.BaseAliPay
    groupId string // 分组ID
}

func (gd *groupDelete) SetGroupId(groupId string) {
    match, _ := regexp.MatchString(`^[0-9]{1,10}$`, groupId)
    if match {
        gd.groupId = groupId
    } else {
        panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "分组ID不合法", nil))
    }
}

func (gd *groupDelete) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(gd.groupId) == 0 {
        panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "分组ID不能为空", nil))
    }
    gd.BizContent["group_id"] = gd.groupId

    return gd.GetRequest()
}

func NewGroupDelete(appId string) *groupDelete {
    gd := &groupDelete{alipay.NewBase(appId), ""}
    gd.SetMethod("alipay.open.public.group.delete")
    return gd
}
