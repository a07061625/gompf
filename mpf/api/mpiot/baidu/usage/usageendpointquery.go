/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/2 0002
 * Time: 9:05
 */
package usage

import (
    "regexp"
    "time"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/mpiot"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 查询特定实例某个时间段内的使用量
type usageEndpointQuery struct {
    mpiot.BaseBaiDu
    endpointName string // endpoint名称
    startDay     string // 开始日期(包含)
    endDay       string // 结束日期(不包含)
}

func (ueq *usageEndpointQuery) SetEndpointName(endpointName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, endpointName)
    if match {
        ueq.endpointName = endpointName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "endpoint名称不合法", nil))
    }
}

func (ueq *usageEndpointQuery) SetQueryTime(startTime, endTime int) {
    if startTime <= 1000000000 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "开始时间不合法", nil))
    } else if endTime <= startTime {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "结束时间必须大于开始时间", nil))
    }

    st := time.Unix(int64(startTime), 0)
    et := time.Unix(int64(endTime), 0)
    startDay := st.Format("2006-01-02")
    endDay := et.Format("2006-01-02")
    if startDay == endDay {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "开始时间和结束时间不能在同一天", nil))
    }

    ueq.startDay = startDay
    ueq.endDay = endDay
}

func (ueq *usageEndpointQuery) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ueq.endpointName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "endpoint名称不能为空", nil))
    }
    if len(ueq.startDay) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "开始日期不能为空", nil))
    }
    ueq.ServiceUri = "/v1/endpoint/" + ueq.endpointName + "/usage-query"
    ueq.ReqData["start"] = ueq.startDay
    ueq.ReqData["end"] = ueq.endDay
    ueq.ReqUrl = ueq.GetServiceUrl() + "?" + mpf.HttpCreateParams(ueq.ReqData, "none", 1)

    return ueq.GetRequest()
}

func NewUsageEndpointQuery() *usageEndpointQuery {
    ueq := &usageEndpointQuery{mpiot.NewBaseBaiDu(), "", "", ""}
    ueq.ReqMethod = fasthttp.MethodPost
    return ueq
}
