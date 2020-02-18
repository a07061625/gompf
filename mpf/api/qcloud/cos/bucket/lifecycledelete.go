/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/22 0022
 * Time: 9:43
 */
package bucket

import (
    "github.com/a07061625/gompf/mpf/api/qcloud"
    "github.com/valyala/fasthttp"
)

// 删除存储桶的生命周期设置
type lifeCycleDelete struct {
    qcloud.BaseCos
}

func (lcd *lifeCycleDelete) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    lcd.ReqURI = "http://" + lcd.ReqHeader["Host"] + lcd.ReqUri + "?lifecycle"
    return lcd.GetRequest()
}

func NewLifeCycleDelete() *lifeCycleDelete {
    lcd := &lifeCycleDelete{qcloud.NewCos()}
    lcd.ReqMethod = fasthttp.MethodDelete
    lcd.SetParamData("lifecycle", "")
    return lcd
}
