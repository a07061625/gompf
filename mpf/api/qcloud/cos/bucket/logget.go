/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/22 0022
 * Time: 14:39
 */
package bucket

import (
    "github.com/a07061625/gompf/mpf/api/qcloud"
    "github.com/valyala/fasthttp"
)

// 获取日志记录
type logGet struct {
    qcloud.BaseCos
}

func (lg *logGet) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    lg.ReqUrl = "http://" + lg.ReqHeader["Host"] + lg.ReqUri + "?logging"
    return lg.GetRequest()
}

func NewLogGet() *logGet {
    lg := &logGet{qcloud.NewCos()}
    lg.SetParamData("logging", "")
    return lg
}
