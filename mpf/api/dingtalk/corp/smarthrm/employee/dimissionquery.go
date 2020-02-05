/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/31 0031
 * Time: 22:33
 */
package employee

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/dingtalk"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 查询企业离职员工列表
type dimissionQuery struct {
    dingtalk.BaseCorp
    corpId   string
    agentTag string
    atType   string
}

func (dq *dimissionQuery) SetOffset(offset int) {
    if offset >= 0 {
        dq.ExtendData["offset"] = offset
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "偏移量不合法", nil))
    }
}

func (dq *dimissionQuery) SetSize(size int) {
    if (size > 0) && (size <= 50) {
        dq.ExtendData["size"] = size
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "分页大小不合法", nil))
    }
}

func (dq *dimissionQuery) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    dq.ReqUrl = dingtalk.UrlService + "/topapi/smartwork/hrm/employee/querydimission?access_token=" + dingtalk.NewUtil().GetAccessToken(dq.corpId, dq.agentTag, dq.atType)

    reqBody := mpf.JsonMarshal(dq.ExtendData)
    client, req := dq.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewDimissionQuery(corpId, agentTag, atType string) *dimissionQuery {
    dq := &dimissionQuery{dingtalk.NewCorp(), "", "", ""}
    dq.corpId = corpId
    dq.agentTag = agentTag
    dq.atType = atType
    dq.ExtendData["offset"] = 0
    dq.ExtendData["size"] = 10
    dq.ReqContentType = project.HttpContentTypeJson
    dq.ReqMethod = fasthttp.MethodPost
    return dq
}
