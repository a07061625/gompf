/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/1 0001
 * Time: 1:05
 */
package user

import (
    "strconv"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/dingtalk"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 获取部门用户userid列表
type departmentMember struct {
    dingtalk.BaseCorp
    corpId   string
    agentTag string
    atType   string
    departId int // 部门id
}

func (dm *departmentMember) SetDepartId(departId int) {
    if departId > 0 {
        dm.departId = departId
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "部门id不合法", nil))
    }
}

func (dm *departmentMember) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if dm.departId <= 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "部门id不能为空", nil))
    }
    dm.ReqData["deptId"] = strconv.Itoa(dm.departId)
    dm.ReqData["access_token"] = dingtalk.NewUtil().GetAccessToken(dm.corpId, dm.agentTag, dm.atType)
    dm.ReqURI = dingtalk.UrlService + "/user/getDeptMember?" + mpf.HTTPCreateParams(dm.ReqData, "none", 1)

    return dm.GetRequest()
}

func NewDepartmentMember(corpId, agentTag, atType string) *departmentMember {
    dm := &departmentMember{dingtalk.NewCorp(), "", "", "", 0}
    dm.corpId = corpId
    dm.agentTag = agentTag
    dm.atType = atType
    return dm
}
