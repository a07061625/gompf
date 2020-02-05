/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/1 0001
 * Time: 1:05
 */
package user

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/dingtalk"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 根据unionid获取userid
type unionId2UserId struct {
    dingtalk.BaseCorp
    corpId   string
    agentTag string
    atType   string
    unionId  string // 唯一标识
}

func (uu *unionId2UserId) SetUnionId(unionId string) {
    if len(unionId) > 0 {
        uu.unionId = unionId
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "唯一标识不合法", nil))
    }
}

func (uu *unionId2UserId) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(uu.unionId) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "唯一标识不能为空", nil))
    }
    uu.ReqData["unionid"] = uu.unionId
    uu.ReqData["access_token"] = dingtalk.NewUtil().GetAccessToken(uu.corpId, uu.agentTag, uu.atType)
    uu.ReqUrl = dingtalk.UrlService + "/user/getUseridByUnionid?" + mpf.HttpCreateParams(uu.ReqData, "none", 1)

    return uu.GetRequest()
}

func NewUnionId2UserId(corpId, agentTag, atType string) *unionId2UserId {
    uu := &unionId2UserId{dingtalk.NewCorp(), "", "", "", ""}
    uu.corpId = corpId
    uu.agentTag = agentTag
    uu.atType = atType
    return uu
}
