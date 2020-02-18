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

// 获取离职员工离职信息
type dimissionList struct {
    dingtalk.BaseCorp
    corpId   string
    agentTag string
    atType   string
    userList []string // 员工列表
}

func (dl *dimissionList) SetUserList(userList []string) {
    if len(userList) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "员工列表不能为空", nil))
    } else if len(userList) > 50 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "员工不能超过50个", nil))
    }

    dl.userList = make([]string, 0)
    for _, v := range userList {
        match, _ := regexp.MatchString(project.RegexDigitAlpha, v)
        if match {
            dl.userList = append(dl.userList, v)
        }
    }
}

func (dl *dimissionList) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(dl.userList) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "员工列表不能为空", nil))
    }
    dl.ExtendData["userid_list"] = strings.Join(dl.userList, ",")

    dl.ReqUrl = dingtalk.UrlService + "/topapi/smartwork/hrm/employee/listdimission?access_token=" + dingtalk.NewUtil().GetAccessToken(dl.corpId, dl.agentTag, dl.atType)

    reqBody := mpf.JsonMarshal(dl.ExtendData)
    dlient, req := dl.GetRequest()
    req.SetBody([]byte(reqBody))

    return dlient, req
}

func NewDimissionList(corpId, agentTag, atType string) *dimissionList {
    dl := &dimissionList{dingtalk.NewCorp(), "", "", "", make([]string, 0)}
    dl.corpId = corpId
    dl.agentTag = agentTag
    dl.atType = atType
    dl.ReqContentType = project.HTTPContentTypeJSON
    dl.ReqMethod = fasthttp.MethodPost
    return dl
}
