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

// 查询企业待入职员工列表
type preEntryQuery struct {
    dingtalk.BaseCorp
    corpId   string
    agentTag string
    atType   string
}

func (peq *preEntryQuery) SetOffset(offset int) {
    if offset >= 0 {
        peq.ExtendData["offset"] = offset
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "偏移量不合法", nil))
    }
}

func (peq *preEntryQuery) SetSize(size int) {
    if (size > 0) && (size <= 50) {
        peq.ExtendData["size"] = size
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "分页大小不合法", nil))
    }
}

func (peq *preEntryQuery) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    peq.ReqUrl = dingtalk.UrlService + "/topapi/smartwork/hrm/employee/querypreentry?access_token=" + dingtalk.NewUtil().GetAccessToken(peq.corpId, peq.agentTag, peq.atType)

    reqBody := mpf.JsonMarshal(peq.ExtendData)
    client, req := peq.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewPreEntryQuery(corpId, agentTag, atType string) *preEntryQuery {
    peq := &preEntryQuery{dingtalk.NewCorp(), "", "", ""}
    peq.corpId = corpId
    peq.agentTag = agentTag
    peq.atType = atType
    peq.ExtendData["offset"] = 0
    peq.ExtendData["size"] = 10
    peq.ReqContentType = project.HttpContentTypeJson
    peq.ReqMethod = fasthttp.MethodPost
    return peq
}
