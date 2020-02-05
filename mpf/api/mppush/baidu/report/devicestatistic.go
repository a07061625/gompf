/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/3 0003
 * Time: 22:30
 */
package report

import (
    "github.com/a07061625/gompf/mpf/api/mppush"
    "github.com/valyala/fasthttp"
)

// 当前应用的设备统计信息
type deviceStatistic struct {
    mppush.BaseBaiDu
}

func (ds *deviceStatistic) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    return ds.GetRequest()
}

func NewDeviceStatistic() *deviceStatistic {
    ds := &deviceStatistic{mppush.NewBaseBaiDu()}
    ds.ServiceUri = "/report/statistic_device"
    return ds
}
