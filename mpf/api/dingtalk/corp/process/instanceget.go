package process

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/dingtalk"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 获取单个审批实例
type instanceGet struct {
    dingtalk.BaseCorp
    corpId     string
    agentTag   string
    instanceId string // 审批实例id
}

func (ig *instanceGet) SetInstanceId(instanceId string) {
    if len(instanceId) > 0 {
        ig.instanceId = instanceId
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "实例id不合法", nil))
    }
}

func (ig *instanceGet) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ig.instanceId) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "实例id不能为空", nil))
    }
    ig.ExtendData["process_instance_id"] = ig.instanceId

    ig.ReqUrl = dingtalk.UrlService + "/topapi/processinstance/get?access_token=" + dingtalk.NewUtil().GetCorpAccessToken(ig.corpId, ig.agentTag)

    reqBody := mpf.JsonMarshal(ig.ExtendData)
    client, req := ig.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewInstanceGet(corpId, agentTag string) *instanceGet {
    ig := &instanceGet{dingtalk.NewCorp(), "", "", ""}
    ig.corpId = corpId
    ig.agentTag = agentTag
    ig.ReqContentType = project.HTTPContentTypeJSON
    ig.ReqMethod = fasthttp.MethodPost
    return ig
}
