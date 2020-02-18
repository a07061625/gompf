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

// 更新审批流
type processSync struct {
    dingtalk.BaseProvider
    corpId            string
    srcProcessCode    string // 源审批码
    targetProcessCode string // 目标审批码
    processName       string // 审批流名称
    bizCategoryId     string // 业务分类标识
}

func (ps *processSync) SetSrcProcessCode(srcProcessCode string) {
    if len(srcProcessCode) > 0 {
        ps.srcProcessCode = srcProcessCode
    } else {
        panic(mperr.NewDingTalkProvider(errorcode.DingTalkProviderParam, "源审批码不合法", nil))
    }
}

func (ps *processSync) SetTargetProcessCode(targetProcessCode string) {
    if len(targetProcessCode) > 0 {
        ps.targetProcessCode = targetProcessCode
    } else {
        panic(mperr.NewDingTalkProvider(errorcode.DingTalkProviderParam, "目标审批码不合法", nil))
    }
}

func (ps *processSync) SetProcessName(processName string) {
    if len(processName) > 0 {
        trueName := []rune(processName)
        ps.processName = string(trueName[:32])
    } else {
        panic(mperr.NewDingTalkProvider(errorcode.DingTalkProviderParam, "审批流名称不合法", nil))
    }
}

func (ps *processSync) SetBizCategoryId(bizCategoryId string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{1,64}$`, bizCategoryId)
    if match {
        ps.bizCategoryId = bizCategoryId
    } else {
        panic(mperr.NewDingTalkProvider(errorcode.DingTalkProviderParam, "业务分类标识不合法", nil))
    }
}

func (ps *processSync) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ps.srcProcessCode) == 0 {
        panic(mperr.NewDingTalkProvider(errorcode.DingTalkProviderParam, "源审批码不能为空", nil))
    }
    if len(ps.targetProcessCode) == 0 {
        panic(mperr.NewDingTalkProvider(errorcode.DingTalkProviderParam, "目标审批码不能为空", nil))
    }
    if len(ps.processName) == 0 {
        panic(mperr.NewDingTalkProvider(errorcode.DingTalkProviderParam, "名称不能为空", nil))
    }
    if len(ps.bizCategoryId) == 0 {
        panic(mperr.NewDingTalkProvider(errorcode.DingTalkProviderParam, "业务分类标识不能为空", nil))
    }
    ps.ExtendData["src_process_code"] = ps.srcProcessCode
    ps.ExtendData["target_process_code"] = ps.targetProcessCode
    ps.ExtendData["process_name"] = ps.processName
    ps.ExtendData["biz_category_id"] = ps.bizCategoryId

    ps.ReqUrl = dingtalk.UrlService + "/topapi/process/sync?access_token=" + dingtalk.NewUtil().GetProviderAuthorizeAccessToken(ps.corpId)

    reqBody := mpf.JSONMarshal(ps.ReqData)
    client, req := ps.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewProcessSync(corpId string) *processSync {
    conf := dingtalk.NewConfig().GetProvider()
    ps := &processSync{dingtalk.NewProvider(), "", "", "", "", ""}
    ps.corpId = corpId
    ps.ExtendData["agent_id"] = strconv.Itoa(conf.GetSuiteId())
    ps.ReqContentType = project.HTTPContentTypeJSON
    ps.ReqMethod = fasthttp.MethodPost
    return ps
}
