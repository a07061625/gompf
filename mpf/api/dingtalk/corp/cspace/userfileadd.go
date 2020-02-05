package cspace

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/dingtalk"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 发送钉盘文件给指定用户
type userFileAdd struct {
    dingtalk.BaseCorp
    corpId   string
    agentTag string
    atType   string
    userId   string // 用户ID
    mediaId  string // 媒体ID
    fileName string // 文件名,包含扩展名
}

func (ufa *userFileAdd) SetUserId(userId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, userId)
    if match {
        ufa.userId = userId
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "用户ID不合法", nil))
    }
}

func (ufa *userFileAdd) SetMediaId(mediaId string) {
    if len(mediaId) > 0 {
        ufa.mediaId = mediaId
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "媒体ID不合法", nil))
    }
}

func (ufa *userFileAdd) SetFileName(fileName string) {
    if len(fileName) > 0 {
        ufa.fileName = fileName
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "文件名不合法", nil))
    }
}

func (ufa *userFileAdd) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ufa.userId) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "用户ID不能为空", nil))
    }
    if len(ufa.mediaId) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "媒体ID不能为空", nil))
    }
    if len(ufa.fileName) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "文件名不能为空", nil))
    }
    ufa.ReqData["userid"] = ufa.userId
    ufa.ReqData["media_id"] = ufa.mediaId
    ufa.ReqData["file_name"] = ufa.fileName
    ufa.ReqData["access_token"] = dingtalk.NewUtil().GetAccessToken(ufa.corpId, ufa.agentTag, ufa.atType)
    ufa.ReqUrl = dingtalk.UrlService + "/cspace/add_to_single_chat?" + mpf.HttpCreateParams(ufa.ReqData, "none", 1)

    return ufa.GetRequest()
}

func NewUserFileAdd(corpId, agentTag, atType string) *userFileAdd {
    agentInfo := dingtalk.NewConfig().GetCorp(corpId).GetAgentInfo(agentTag)
    ufa := &userFileAdd{dingtalk.NewCorp(), "", "", "", "", "", ""}
    ufa.corpId = corpId
    ufa.agentTag = agentTag
    ufa.atType = atType
    ufa.ReqData["agent_id"] = agentInfo["id"]
    return ufa
}
