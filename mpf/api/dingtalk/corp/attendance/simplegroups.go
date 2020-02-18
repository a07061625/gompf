package attendance

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/dingtalk"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 企业考勤组详情
type simpleGroups struct {
    dingtalk.BaseCorp
    corpId   string
    agentTag string
}

func (sg *simpleGroups) SetOffset(offset int) {
    if offset >= 0 {
        sg.ExtendData["offset"] = offset
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "偏移量不合法", nil))
    }
}

func (sg *simpleGroups) SetSize(size int) {
    if size > 0 {
        if size > 10 {
            sg.ExtendData["size"] = 10
        } else {
            sg.ExtendData["size"] = size
        }
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "分页大小不合法", nil))
    }
}

func (sg *simpleGroups) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    sg.ReqUrl = dingtalk.UrlService + "/topapi/attendance/getsimplegroups?access_token=" + dingtalk.NewUtil().GetCorpAccessToken(sg.corpId, sg.agentTag)

    reqBody := mpf.JSONMarshal(sg.ExtendData)
    client, req := sg.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewSimpleGroups(corpId, agentTag string) *simpleGroups {
    sg := &simpleGroups{dingtalk.NewCorp(), "", ""}
    sg.corpId = corpId
    sg.agentTag = agentTag
    sg.ExtendData["offset"] = 0
    sg.ExtendData["size"] = 10
    sg.ReqContentType = project.HTTPContentTypeJSON
    sg.ReqMethod = fasthttp.MethodPost
    return sg
}
