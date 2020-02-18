package cspace

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

// 授权用户访问企业自定义空间
type customSpaceGrant struct {
    dingtalk.BaseCorp
    corpId    string
    agentTag  string
    domain    string   // 域名
    grantType string   // 权限类型
    userId    string   // 用户ID
    grantPath string   // 授权路径
    fileList  []string // 文件列表
}

func (csg *customSpaceGrant) SetDomain(domain string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{1,10}$`, domain)
    if match {
        csg.domain = domain
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "域名不合法", nil))
    }
}

func (csg *customSpaceGrant) SetGrantType(grantType string) {
    if (grantType == "add") || (grantType == "download") {
        csg.grantType = grantType
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "权限类型不合法", nil))
    }
}

func (csg *customSpaceGrant) SetUserId(userId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, userId)
    if match {
        csg.userId = userId
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "用户ID不合法", nil))
    }
}

func (csg *customSpaceGrant) SetGrantPath(grantPath string) {
    if len(grantPath) > 0 {
        csg.grantPath = grantPath
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "授权路径不合法", nil))
    }
}

func (csg *customSpaceGrant) SetFileList(fileList []string) {
    if len(fileList) > 0 {
        csg.fileList = fileList
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "文件列表不能为空", nil))
    }
}

func (csg *customSpaceGrant) SetDuration(duration int) {
    if (duration > 0) && (duration <= 3600) {
        csg.ReqData["duration"] = strconv.Itoa(duration)
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "权限有效时间不合法", nil))
    }
}

func (csg *customSpaceGrant) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(csg.userId) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "用户ID不能为空", nil))
    }
    if csg.grantType == "add" {
        if len(csg.grantPath) == 0 {
            panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "授权路径不能为空", nil))
        }
        csg.ReqData["path"] = csg.grantPath
    } else if csg.grantType == "download" {
        if len(csg.fileList) == 0 {
            panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "文件列表不能为空", nil))
        }
        csg.ReqData["fileids"] = strings.Join(csg.fileList, ",")
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "权限类型不能为空", nil))
    }

    csg.ReqData["userid"] = csg.userId
    csg.ReqData["type"] = csg.grantType
    if len(csg.domain) > 0 {
        csg.ReqData["domain"] = csg.domain
        csg.ReqData["access_token"] = dingtalk.NewUtil().GetCorpAccessToken(csg.corpId, csg.agentTag)
    } else {
        agentInfo := dingtalk.NewConfig().GetCorp(csg.corpId).GetAgentInfo(csg.agentTag)
        csg.ReqData["agent_id"] = agentInfo["id"]
        csg.ReqData["access_token"] = dingtalk.NewUtil().GetProviderAuthorizeAccessToken(csg.corpId)
    }

    csg.ReqUrl = dingtalk.UrlService + "/cspace/grant_custom_space?" + mpf.HTTPCreateParams(csg.ReqData, "none", 1)

    return csg.GetRequest()
}

func NewCustomSpaceGrant(corpId, agentTag string) *customSpaceGrant {
    csg := &customSpaceGrant{dingtalk.NewCorp(), "", "", "", "", "", "", make([]string, 0)}
    csg.corpId = corpId
    csg.agentTag = agentTag
    csg.ReqData["duration"] = "30"
    return csg
}
