/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/31 0031
 * Time: 23:52
 */
package employee

import (
    "regexp"
    "strings"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/dingtalk"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 获取员工花名册字段信息
type employeeList struct {
    dingtalk.BaseCorp
    corpId   string
    agentTag string
    atType   string
    userList []string // 员工列表
}

func (el *employeeList) SetUserList(userList []string) {
    if len(userList) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "员工列表不能为空", nil))
    } else if len(userList) > 20 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "员工不能超过20个", nil))
    }

    el.userList = make([]string, 0)
    for _, v := range userList {
        match, _ := regexp.MatchString(project.RegexDigitAlpha, v)
        if match {
            el.userList = append(el.userList, v)
        }
    }
}

func (el *employeeList) SetFieldList(fieldList []string) {
    if len(fieldList) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "花名册字段列表不能为空", nil))
    } else if len(fieldList) > 20 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "花名册字段不能超过20个", nil))
    }

    fields := make([]string, 0)
    for _, v := range fieldList {
        if len(v) > 0 {
            fields = append(fields, v)
        }
    }
    if len(fields) > 0 {
        el.ExtendData["field_filter_list"] = strings.Join(fields, ",")
    }
}

func (el *employeeList) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(el.userList) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "员工列表不能为空", nil))
    }
    el.ExtendData["userid_list"] = strings.Join(el.userList, ",")

    el.ReqUrl = dingtalk.UrlService + "/topapi/smartwork/hrm/employee/list?access_token=" + dingtalk.NewUtil().GetAccessToken(el.corpId, el.agentTag, el.atType)

    reqBody := mpf.JSONMarshal(el.ExtendData)
    client, req := el.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewEmployeeList(corpId, agentTag, atType string) *employeeList {
    el := &employeeList{dingtalk.NewCorp(), "", "", "", make([]string, 0)}
    el.corpId = corpId
    el.agentTag = agentTag
    el.atType = atType
    el.ReqContentType = project.HTTPContentTypeJSON
    el.ReqMethod = fasthttp.MethodPost
    return el
}
