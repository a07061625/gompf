package extcontact

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/dingtalk"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 获取外部联系人列表
type extContactList struct {
    dingtalk.BaseCorp
    corpId   string
    agentTag string
    atType   string
}

func (ecl *extContactList) SetOffset(offset int) {
    if offset >= 0 {
        ecl.ExtendData["offset"] = offset
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "偏移量不合法", nil))
    }
}

func (ecl *extContactList) SetSize(size int) {
    if (size > 0) && (size <= 100) {
        ecl.ExtendData["size"] = size
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "分页大小不合法", nil))
    }
}

func (ecl *extContactList) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    ecl.ReqURI = dingtalk.UrlService + "/topapi/extcontact/list?access_token=" + dingtalk.NewUtil().GetAccessToken(ecl.corpId, ecl.agentTag, ecl.atType)

    reqBody := mpf.JSONMarshal(ecl.ExtendData)
    client, req := ecl.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewExtContactList(corpId, agentTag, atType string) *extContactList {
    ecl := &extContactList{dingtalk.NewCorp(), "", "", ""}
    ecl.corpId = corpId
    ecl.agentTag = agentTag
    ecl.atType = atType
    ecl.ExtendData["offset"] = 0
    ecl.ExtendData["size"] = 10
    ecl.ReqContentType = project.HTTPContentTypeJSON
    ecl.ReqMethod = fasthttp.MethodPost
    return ecl
}
