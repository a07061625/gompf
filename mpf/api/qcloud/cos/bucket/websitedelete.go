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

// 删除存储桶的静态网站配置
type websiteDelete struct {
    qcloud.BaseCos
}

func (wd *websiteDelete) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    wd.ReqURI = "http://" + wd.ReqHeader["Host"] + wd.ReqUri + "?website"
    return wd.GetRequest()
}

func NewWebsiteDelete() *websiteDelete {
    wd := &websiteDelete{qcloud.NewCos()}
    wd.ReqMethod = fasthttp.MethodDelete
    wd.SetParamData("website", "")
    return wd
}
