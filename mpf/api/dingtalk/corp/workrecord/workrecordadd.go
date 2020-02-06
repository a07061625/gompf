/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/1 0001
 * Time: 1:05
 */
package workrecord

import (
    "regexp"
    "time"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/dingtalk"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 发起待办
type workRecordAdd struct {
    dingtalk.BaseCorp
    corpId     string
    agentTag   string
    atType     string
    userId     string                   // 用户ID
    createTime int64                    // 待办时间
    title      string                   // 标题
    url        string                   // 跳转链接
    formList   []map[string]interface{} // 表单列表
}

func (wra *workRecordAdd) SetUserId(userId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, userId)
    if match {
        wra.userId = userId
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "用户ID不合法", nil))
    }
}

func (wra *workRecordAdd) SetCreateTime(createTime int64) {
    if createTime > time.Now().Unix() {
        wra.createTime = 1000 * createTime
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "待办时间不合法", nil))
    }
}

func (wra *workRecordAdd) SetTitle(title string) {
    if len(title) > 0 {
        wra.title = title
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "标题不合法", nil))
    }
}

func (wra *workRecordAdd) SetUrl(url string) {
    match, _ := regexp.MatchString(project.RegexUrlHttp, url)
    if match {
        wra.url = url
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "跳转链接不合法", nil))
    }
}

func (wra *workRecordAdd) SetFormList(formList []map[string]interface{}) {
    if len(formList) > 0 {
        wra.formList = formList
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "表单列表不合法", nil))
    }
}

func (wra *workRecordAdd) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(wra.userId) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "用户ID不能为空", nil))
    }
    if wra.createTime <= 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "待办时间不能为空", nil))
    }
    if len(wra.title) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "标题不能为空", nil))
    }
    if len(wra.url) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "跳转链接不能为空", nil))
    }
    if len(wra.formList) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "表单列表不能为空", nil))
    }
    wra.ExtendData["userid"] = wra.userId
    wra.ExtendData["create_time"] = wra.createTime
    wra.ExtendData["title"] = wra.title
    wra.ExtendData["url"] = wra.url
    wra.ExtendData["formItemList"] = wra.formList

    wra.ReqUrl = dingtalk.UrlService + "/topapi/workrecord/add?access_token=" + dingtalk.NewUtil().GetAccessToken(wra.corpId, wra.agentTag, wra.atType)

    reqBody := mpf.JsonMarshal(wra.ExtendData)
    client, req := wra.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewWorkRecordAdd(corpId, agentTag, atType string) *workRecordAdd {
    wra := &workRecordAdd{dingtalk.NewCorp(), "", "", "", "", 0, "", "", make([]map[string]interface{}, 0)}
    wra.corpId = corpId
    wra.agentTag = agentTag
    wra.atType = atType
    wra.ReqContentType = project.HttpContentTypeJson
    wra.ReqMethod = fasthttp.MethodPost
    return wra
}
