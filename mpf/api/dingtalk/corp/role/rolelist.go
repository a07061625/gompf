/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/31 0031
 * Time: 22:33
 */
package role

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/dingtalk"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 获取角色列表
type roleList struct {
    dingtalk.BaseCorp
    corpId   string
    agentTag string
    atType   string
}

func (rl *roleList) SetOffset(offset int) {
    if offset >= 0 {
        rl.ExtendData["offset"] = offset
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "偏移量不合法", nil))
    }
}

func (rl *roleList) SetSize(size int) {
    if (size > 0) && (size <= 200) {
        rl.ExtendData["size"] = size
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "分页大小不合法", nil))
    }
}

func (rl *roleList) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    rl.ReqUrl = dingtalk.UrlService + "/topapi/role/list?access_token=" + dingtalk.NewUtil().GetAccessToken(rl.corpId, rl.agentTag, rl.atType)

    reqBody := mpf.JSONMarshal(rl.ExtendData)
    client, req := rl.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewRoleList(corpId, agentTag, atType string) *roleList {
    rl := &roleList{dingtalk.NewCorp(), "", "", ""}
    rl.corpId = corpId
    rl.agentTag = agentTag
    rl.atType = atType
    rl.ExtendData["offset"] = 0
    rl.ExtendData["size"] = 10
    rl.ReqContentType = project.HTTPContentTypeJSON
    rl.ReqMethod = fasthttp.MethodPost
    return rl
}
