/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/31 0031
 * Time: 23:33
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

// 获取离职员工通讯录字段信息
type contactList struct {
    dingtalk.BaseCorp
    corpId   string
    agentTag string
    userList []string // 员工列表
}

func (cl *contactList) SetUserList(userList []string) {
    if len(userList) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "员工列表不能为空", nil))
    } else if len(userList) > 20 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "员工不能超过20个", nil))
    }

    cl.userList = make([]string, 0)
    for _, v := range userList {
        match, _ := regexp.MatchString(project.RegexDigitAlpha, v)
        if match {
            cl.userList = append(cl.userList, v)
        }
    }
}

func (cl *contactList) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(cl.userList) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "员工列表不能为空", nil))
    }
    cl.ExtendData["userid_list"] = strings.Join(cl.userList, ",")

    cl.ReqUrl = dingtalk.UrlService + "/topapi/smartwork/hrm/employee/listcontact?access_token=" + dingtalk.NewUtil().GetProviderAuthorizeAccessToken(cl.corpId)

    reqBody := mpf.JsonMarshal(cl.ExtendData)
    client, req := cl.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewContactList(corpId, agentTag string) *contactList {
    cl := &contactList{dingtalk.NewCorp(), "", "", make([]string, 0)}
    cl.corpId = corpId
    cl.agentTag = agentTag
    cl.ReqContentType = project.HttpContentTypeJson
    cl.ReqMethod = fasthttp.MethodPost
    return cl
}
