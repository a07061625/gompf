package microapp

import (
    "github.com/a07061625/gompf/mpf/api/dingtalk"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/valyala/fasthttp"
)

// 获取应用列表
type appList struct {
    dingtalk.BaseCorp
    corpId   string
    agentTag string
    atType   string
}

func (al *appList) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    al.ReqUrl = dingtalk.UrlService + "/microapp/list?access_token=" + dingtalk.NewUtil().GetAccessToken(al.corpId, al.agentTag, al.atType)

    client, req := al.GetRequest()
    req.SetBody([]byte("[]"))

    return client, req
}

func NewAppList(corpId, agentTag, atType string) *appList {
    al := &appList{dingtalk.NewCorp(), "", "", ""}
    al.corpId = corpId
    al.agentTag = agentTag
    al.atType = atType
    al.ReqContentType = project.HTTPContentTypeJSON
    al.ReqMethod = fasthttp.MethodPost
    return al
}
