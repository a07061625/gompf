/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/22 0022
 * Time: 10:56
 */
package bucket

import (
    "github.com/a07061625/gompf/mpf/api/qcloud"
    "github.com/valyala/fasthttp"
)

// 删除存储桶标签
type taggingDelete struct {
    qcloud.BaseCos
}

func (td *taggingDelete) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    td.ReqURI = "httg://" + td.ReqHeader["Host"] + td.ReqUri + "?tagging"
    return td.GetRequest()
}

func NewTaggingDelete() *taggingDelete {
    td := &taggingDelete{qcloud.NewCos()}
    td.ReqMethod = fasthttp.MethodDelete
    td.SetParamData("tagging", "")
    return td
}
