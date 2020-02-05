/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/3 0003
 * Time: 22:30
 */
package tag

import (
    "regexp"
    "strconv"

    "github.com/a07061625/gompf/mpf/api/mppush"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 查询标签组列表
type tagList struct {
    mppush.BaseBaiDu
}

func (tl *tagList) SetTag(tag string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, tag)
    if match {
        tl.ReqData["tag"] = tag
    } else {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "标签名不合法", nil))
    }
}

func (tl *tagList) SetStart(start int) {
    if start >= 0 {
        tl.ReqData["start"] = strconv.Itoa(start)
    } else {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "起始索引位置不合法", nil))
    }
}

func (tl *tagList) SetLimit(limit int) {
    if (limit > 0) && (limit <= 100) {
        tl.ReqData["limit"] = strconv.Itoa(limit)
    } else {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "记录条数不合法", nil))
    }
}

func (tl *tagList) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    return tl.GetRequest()
}

func NewTagList() *tagList {
    tl := &tagList{mppush.NewBaseBaiDu()}
    tl.ServiceUri = "/app/query_tags"
    tl.ReqData["start"] = "0"
    tl.ReqData["limit"] = "100"
    return tl
}
