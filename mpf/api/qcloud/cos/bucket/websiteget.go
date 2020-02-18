/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/22 0022
 * Time: 11:07
 */
package bucket

import (
    "github.com/a07061625/gompf/mpf/api/qcloud"
    "github.com/valyala/fasthttp"
)

// 获取存储桶的静态网站配置
type websiteGet struct {
    qcloud.BaseCos
}

func (wg *websiteGet) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    wg.ReqURI = "http://" + wg.ReqHeader["Host"] + wg.ReqUri + "?website"
    return wg.GetRequest()
}

func NewWebsiteGet() *websiteGet {
    wg := &websiteGet{qcloud.NewCos()}
    wg.SetParamData("website", "")
    return wg
}
