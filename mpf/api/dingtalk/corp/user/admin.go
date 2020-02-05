/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/1 0001
 * Time: 1:01
 */
package user

import (
    "github.com/a07061625/gompf/mpf/api/dingtalk"
    "github.com/valyala/fasthttp"
)

// 获取管理员列表
type admin struct {
    dingtalk.BaseCorp
    corpId   string
    agentTag string
    atType   string
}

func (a *admin) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    a.ReqUrl = dingtalk.UrlService + "/user/get_admin?access_token=" + dingtalk.NewUtil().GetAccessToken(a.corpId, a.agentTag, a.atType)

    return a.GetRequest()
}

func NewAdmin(corpId, agentTag, atType string) *admin {
    a := &admin{dingtalk.NewCorp(), "", "", ""}
    a.corpId = corpId
    a.agentTag = agentTag
    a.atType = atType
    return a
}
