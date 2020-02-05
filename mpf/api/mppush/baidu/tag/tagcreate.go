/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/3 0003
 * Time: 22:30
 */
package tag

import (
    "regexp"

    "github.com/a07061625/gompf/mpf/api/mppush"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 创建标签组
type tagCreate struct {
    mppush.BaseBaiDu
    tag string // 标签名
}

func (tc *tagCreate) SetTag(tag string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, tag)
    if match {
        tc.tag = tag
    } else {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "标签名不合法", nil))
    }
}

func (tc *tagCreate) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(tc.tag) == 0 {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "标签名不能为空", nil))
    }
    tc.ReqData["tag"] = tc.tag

    return tc.GetRequest()
}

func NewTagCreate() *tagCreate {
    tc := &tagCreate{mppush.NewBaseBaiDu(), ""}
    tc.ServiceUri = "/app/create_tag"
    return tc
}
