package process

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

// 复制审批流
type processCopy struct {
    dingtalk.BaseProvider
    corpId        string
    processCode   string // 审批码
    processName   string // 审批流名称
    processDesc   string // 审批流描述
    bizCategoryId string // 业务分类标识
}

func (pc *processCopy) SetProcessCode(processCode string) {
    if len(processCode) > 0 {
        pc.processCode = processCode
    } else {
        panic(mperr.NewDingTalkProvider(errorcode.DingTalkProviderParam, "审批码不合法", nil))
    }
}

func (pc *processCopy) SetProcessName(processName string) {
    if len(processName) > 0 {
        trueName := []rune(processName)
        pc.processName = string(trueName[:32])
    } else {
        panic(mperr.NewDingTalkProvider(errorcode.DingTalkProviderParam, "审批流名称不合法", nil))
    }
}

func (pc *processCopy) SetProcessDesc(processDesc string) {
    if len(processDesc) > 0 {
        pc.processDesc = processDesc
    } else {
        panic(mperr.NewDingTalkProvider(errorcode.DingTalkProviderParam, "审批流描述不合法", nil))
    }
}

func (pc *processCopy) SetBizCategoryId(bizCategoryId string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{1,64}$`, bizCategoryId)
    if match {
        pc.bizCategoryId = bizCategoryId
    } else {
        panic(mperr.NewDingTalkProvider(errorcode.DingTalkProviderParam, "业务分类标识不合法", nil))
    }
}

func (pc *processCopy) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(pc.processCode) == 0 {
        panic(mperr.NewDingTalkProvider(errorcode.DingTalkProviderParam, "审批码不能为空", nil))
    }
    if len(pc.processName) == 0 {
        panic(mperr.NewDingTalkProvider(errorcode.DingTalkProviderParam, "名称不能为空", nil))
    }
    if len(pc.processDesc) == 0 {
        panic(mperr.NewDingTalkProvider(errorcode.DingTalkProviderParam, "描述不能为空", nil))
    }
    if len(pc.bizCategoryId) == 0 {
        panic(mperr.NewDingTalkProvider(errorcode.DingTalkProviderParam, "业务分类标识不能为空", nil))
    }
    pc.ExtendData["process_code"] = pc.processCode
    pc.ExtendData["process_name"] = pc.processName
    pc.ExtendData["description"] = pc.processDesc
    pc.ExtendData["biz_category_id"] = pc.bizCategoryId

    pc.ReqUrl = dingtalk.UrlService + "/topapi/process/copy?access_token=" + dingtalk.NewUtil().GetProviderAuthorizeAccessToken(pc.corpId)

    reqBody := mpf.JsonMarshal(pc.ExtendData)
    client, req := pc.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewProcessCopy(corpId string) *processCopy {
    conf := dingtalk.NewConfig().GetProvider()
    pc := &processCopy{dingtalk.NewProvider(), "", "", "", "", ""}
    pc.corpId = corpId
    pc.ExtendData["agent_id"] = strconv.Itoa(conf.GetSuiteId())
    pc.ReqContentType = project.HTTPContentTypeJSON
    pc.ReqMethod = fasthttp.MethodPost
    return pc
}
