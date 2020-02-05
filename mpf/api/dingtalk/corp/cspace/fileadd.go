package cspace

import (
    "regexp"
    "strconv"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/dingtalk"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 新增文件到用户钉盘
type fileAdd struct {
    dingtalk.BaseCorp
    corpId   string
    agentTag string
    atType   string
    authCode string // 授权码
    mediaId  string // 媒体ID
    spaceId  string // 空间ID
    folderId string // 文件夹ID
    fileName string // 文件名,包含扩展名
}

func (fa *fileAdd) SetAuthCode(authCode string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, authCode)
    if match {
        fa.authCode = authCode
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "授权码不合法", nil))
    }
}

func (fa *fileAdd) SetMediaId(mediaId string) {
    if len(mediaId) > 0 {
        fa.mediaId = mediaId
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "媒体ID不合法", nil))
    }
}

func (fa *fileAdd) SetSpaceId(spaceId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, spaceId)
    if match {
        fa.spaceId = spaceId
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "空间ID不合法", nil))
    }
}

func (fa *fileAdd) SetFolderId(folderId string) {
    if len(folderId) > 0 {
        fa.folderId = folderId
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "文件夹ID不合法", nil))
    }
}

func (fa *fileAdd) SetFileName(fileName string) {
    if len(fileName) > 0 {
        fa.fileName = fileName
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "文件名不合法", nil))
    }
}

func (fa *fileAdd) SetOverwrite(overwrite int) {
    if (overwrite == 0) || (overwrite == 1) {
        fa.ReqData["overwrite"] = strconv.Itoa(overwrite)
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "同名文件覆盖标识不合法", nil))
    }
}

func (fa *fileAdd) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(fa.authCode) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "授权码不能为空", nil))
    }
    if len(fa.mediaId) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "媒体ID不能为空", nil))
    }
    if len(fa.spaceId) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "空间ID不能为空", nil))
    }
    if len(fa.folderId) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "文件夹ID不能为空", nil))
    }
    if len(fa.fileName) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "文件名不能为空", nil))
    }
    fa.ReqData["code"] = fa.authCode
    fa.ReqData["media_id"] = fa.mediaId
    fa.ReqData["space_id"] = fa.spaceId
    fa.ReqData["folder_id"] = fa.folderId
    fa.ReqData["name"] = fa.fileName
    fa.ReqData["access_token"] = dingtalk.NewUtil().GetAccessToken(fa.corpId, fa.agentTag, fa.atType)
    fa.ReqUrl = dingtalk.UrlService + "/cspace/add?" + mpf.HttpCreateParams(fa.ReqData, "none", 1)

    return fa.GetRequest()
}

func NewFileAdd(corpId, agentTag, atType string) *fileAdd {
    agentInfo := dingtalk.NewConfig().GetCorp(corpId).GetAgentInfo(agentTag)
    fa := &fileAdd{dingtalk.NewCorp(), "", "", "", "", "", "", "", ""}
    fa.corpId = corpId
    fa.agentTag = agentTag
    fa.atType = atType
    fa.ReqData["agent_id"] = agentInfo["id"]
    fa.ReqData["overwrite"] = "0"
    return fa
}
