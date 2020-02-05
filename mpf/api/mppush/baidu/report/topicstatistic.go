/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/3 0003
 * Time: 22:30
 */
package report

import (
    "regexp"

    "github.com/a07061625/gompf/mpf/api/mppush"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 查询分类主题统计信息
type topicStatistic struct {
    mppush.BaseBaiDu
    topicId string // 分类主题标识
}

func (ts *topicStatistic) SetTopicId(topicId string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{1,128}$`, topicId)
    if match {
        ts.topicId = topicId
    } else {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "分类主题标识不合法", nil))
    }
}

func (ts *topicStatistic) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ts.topicId) == 0 {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "分类主题标识不能为空", nil))
    }
    ts.ReqData["topic_id"] = ts.topicId

    return ts.GetRequest()
}

func NewTopicStatistic() *topicStatistic {
    ts := &topicStatistic{mppush.NewBaseBaiDu(), ""}
    ts.ServiceUri = "/report/statistic_topic"
    return ts
}
