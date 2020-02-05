package extcontact

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/dingtalk"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 获取外部联系人标签列表
type labelGroupsList struct {
    dingtalk.BaseCorp
    corpId   string
    agentTag string
    atType   string
}

func (lgl *labelGroupsList) SetOffset(offset int) {
    if offset >= 0 {
        lgl.ExtendData["offset"] = offset
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "偏移量不合法", nil))
    }
}

func (lgl *labelGroupsList) SetSize(size int) {
    if (size > 0) && (size <= 100) {
        lgl.ExtendData["size"] = size
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "分页大小不合法", nil))
    }
}

func (lgl *labelGroupsList) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    lgl.ReqUrl = dingtalk.UrlService + "/topapi/extcontact/listlabelgroups?access_token=" + dingtalk.NewUtil().GetAccessToken(lgl.corpId, lgl.agentTag, lgl.atType)

    reqBody := mpf.JsonMarshal(lgl.ExtendData)
    client, req := lgl.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewLabelGroupsList(corpId, agentTag, atType string) *labelGroupsList {
    lgl := &labelGroupsList{dingtalk.NewCorp(), "", "", ""}
    lgl.corpId = corpId
    lgl.agentTag = agentTag
    lgl.atType = atType
    lgl.ExtendData["offset"] = 0
    lgl.ExtendData["size"] = 10
    lgl.ReqContentType = project.HttpContentTypeJson
    lgl.ReqMethod = fasthttp.MethodPost
    return lgl
}
