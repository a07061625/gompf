/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/31 0031
 * Time: 23:13
 */
package role

import (
    "regexp"
    "strconv"
    "strings"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/dingtalk"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 批量增加员工角色
type batchAdd struct {
    dingtalk.BaseCorp
    corpId   string
    agentTag string
    atType   string
    roleList []string // 角色列表
    userList []string // 用户列表
}

func (ba *batchAdd) SetRoleList(roleList []int) {
    if len(roleList) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "角色列表不能为空", nil))
    } else if len(roleList) > 20 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "角色不能超过20个", nil))
    }

    ba.roleList = make([]string, 0)
    for _, v := range roleList {
        if v > 0 {
            ba.roleList = append(ba.roleList, strconv.Itoa(v))
        }
    }
}

func (ba *batchAdd) SetUserList(userList []string) {
    if len(userList) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "用户列表不能为空", nil))
    } else if len(userList) > 100 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "用户不能超过100个", nil))
    }

    ba.userList = make([]string, 0)
    for _, v := range userList {
        match, _ := regexp.MatchString(project.RegexDigitAlpha, v)
        if match {
            ba.userList = append(ba.userList, v)
        }
    }
}

func (ba *batchAdd) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ba.roleList) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "角色列表不能为空", nil))
    }
    if len(ba.userList) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "用户列表不能为空", nil))
    }
    ba.ExtendData["roleIds"] = strings.Join(ba.roleList, ",")
    ba.ExtendData["userIds"] = strings.Join(ba.userList, ",")

    ba.ReqUrl = dingtalk.UrlService + "/topapi/role/addrolesforemps?access_token=" + dingtalk.NewUtil().GetAccessToken(ba.corpId, ba.agentTag, ba.atType)

    reqBody := mpf.JsonMarshal(ba.ExtendData)
    client, req := ba.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewBatchAdd(corpId, agentTag, atType string) *batchAdd {
    ba := &batchAdd{dingtalk.NewCorp(), "", "", "", make([]string, 0), make([]string, 0)}
    ba.corpId = corpId
    ba.agentTag = agentTag
    ba.atType = atType
    ba.ReqContentType = project.HTTPContentTypeJSON
    ba.ReqMethod = fasthttp.MethodPost
    return ba
}
