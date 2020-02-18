/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/31 0031
 * Time: 22:33
 */
package employee

import (
    "strconv"
    "strings"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/dingtalk"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 查询企业在职员工列表
type onJobQuery struct {
    dingtalk.BaseCorp
    corpId     string
    agentTag   string
    atType     string
    statusList []string // 状态列表
}

func (ojq *onJobQuery) SetStatusList(statusList []int) {
    if len(statusList) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "状态列表不能为空", nil))
    }

    ojq.statusList = make([]string, 0)
    for _, v := range statusList {
        ojq.statusList = append(ojq.statusList, strconv.Itoa(v))
    }
}

func (ojq *onJobQuery) SetOffset(offset int) {
    if offset >= 0 {
        ojq.ExtendData["offset"] = offset
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "偏移量不合法", nil))
    }
}

func (ojq *onJobQuery) SetSize(size int) {
    if (size > 0) && (size <= 20) {
        ojq.ExtendData["size"] = size
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "分页大小不合法", nil))
    }
}

func (ojq *onJobQuery) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ojq.statusList) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "状态列表不能为空", nil))
    }
    ojq.ExtendData["status_list"] = strings.Join(ojq.statusList, ",")

    ojq.ReqUrl = dingtalk.UrlService + "/topapi/smartwork/hrm/employee/queryonjob?access_token=" + dingtalk.NewUtil().GetAccessToken(ojq.corpId, ojq.agentTag, ojq.atType)

    reqBody := mpf.JSONMarshal(ojq.ExtendData)
    client, req := ojq.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewOnJobQuery(corpId, agentTag, atType string) *onJobQuery {
    ojq := &onJobQuery{dingtalk.NewCorp(), "", "", "", make([]string, 0)}
    ojq.corpId = corpId
    ojq.agentTag = agentTag
    ojq.atType = atType
    ojq.ExtendData["offset"] = 0
    ojq.ExtendData["size"] = 10
    ojq.ReqContentType = project.HTTPContentTypeJSON
    ojq.ReqMethod = fasthttp.MethodPost
    return ojq
}
