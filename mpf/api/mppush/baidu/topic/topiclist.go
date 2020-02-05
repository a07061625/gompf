/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/3 0003
 * Time: 22:30
 */
package topic

import (
    "strconv"

    "github.com/a07061625/gompf/mpf/api/mppush"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 查询分类主题列表
type topicList struct {
    mppush.BaseBaiDu
}

func (tl *topicList) SetStart(start int) {
    if start >= 0 {
        tl.ReqData["start"] = strconv.Itoa(start)
    } else {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "起始索引位置不合法", nil))
    }
}

func (tl *topicList) SetLimit(limit int) {
    if (limit > 0) && (limit <= 20) {
        tl.ReqData["limit"] = strconv.Itoa(limit)
    } else {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "记录条数不合法", nil))
    }
}

func (tl *topicList) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    return tl.GetRequest()
}

func NewTopicList() *topicList {
    tl := &topicList{mppush.NewBaseBaiDu()}
    tl.ServiceUri = "/topic/query_list"
    tl.ReqData["start"] = "0"
    tl.ReqData["limit"] = "20"
    return tl
}
