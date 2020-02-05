/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/20 0020
 * Time: 23:18
 */
package service

import (
    "github.com/a07061625/gompf/mpf/api/qcloud"
    "github.com/valyala/fasthttp"
)

// 获取所有存储空间列表
type serviceGet struct {
    qcloud.BaseCos
}

func (sg *serviceGet) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    sg.ReqUrl = "http://" + sg.ReqHeader["Host"] + sg.ReqUri
    return sg.GetRequest()
}

func NewServiceGet() *serviceGet {
    sg := &serviceGet{qcloud.NewCos()}
    sg.SetHeaderData("Host", "service.cos.myqcloud.com")
    return sg
}
