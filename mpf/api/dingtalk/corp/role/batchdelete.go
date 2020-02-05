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

// 批量删除员工角色
type batchDelete struct {
    dingtalk.BaseCorp
    corpId   string
    agentTag string
    atType   string
    roleList []string // 角色列表
    userList []string // 用户列表
}

func (bd *batchDelete) SetRoleList(roleList []int) {
    if len(roleList) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "角色列表不能为空", nil))
    } else if len(roleList) > 20 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "角色不能超过20个", nil))
    }

    bd.roleList = make([]string, 0)
    for _, v := range roleList {
        if v > 0 {
            bd.roleList = append(bd.roleList, strconv.Itoa(v))
        }
    }
}

func (bd *batchDelete) SetUserList(userList []string) {
    if len(userList) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "用户列表不能为空", nil))
    } else if len(userList) > 100 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "用户不能超过100个", nil))
    }

    bd.userList = make([]string, 0)
    for _, v := range userList {
        match, _ := regexp.MatchString(project.RegexDigitAlpha, v)
        if match {
            bd.userList = append(bd.userList, v)
        }
    }
}

func (bd *batchDelete) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(bd.roleList) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "角色列表不能为空", nil))
    }
    if len(bd.userList) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "用户列表不能为空", nil))
    }
    bd.ExtendData["roleIds"] = strings.Join(bd.roleList, ",")
    bd.ExtendData["userIds"] = strings.Join(bd.userList, ",")

    bd.ReqUrl = dingtalk.UrlService + "/topapi/role/removerolesforemps?access_token=" + dingtalk.NewUtil().GetAccessToken(bd.corpId, bd.agentTag, bd.atType)

    reqBody := mpf.JsonMarshal(bd.ExtendData)
    client, req := bd.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewBatchDelete(corpId, agentTag, atType string) *batchDelete {
    bd := &batchDelete{dingtalk.NewCorp(), "", "", "", make([]string, 0), make([]string, 0)}
    bd.corpId = corpId
    bd.agentTag = agentTag
    bd.atType = atType
    bd.ReqContentType = project.HttpContentTypeJson
    bd.ReqMethod = fasthttp.MethodPost
    return bd
}
