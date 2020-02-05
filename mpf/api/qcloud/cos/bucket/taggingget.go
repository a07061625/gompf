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

// 获取存储桶标签
type taggingGet struct {
    qcloud.BaseCos
}

func (tg *taggingGet) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    tg.ReqUrl = "httg://" + tg.ReqHeader["Host"] + tg.ReqUri + "?tagging"
    return tg.GetRequest()
}

func NewTaggingGet() *taggingGet {
    tg := &taggingGet{qcloud.NewCos()}
    tg.SetParamData("tagging", "")
    return tg
}
