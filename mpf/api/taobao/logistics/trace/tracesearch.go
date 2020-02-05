/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/17 0017
 * Time: 9:43
 */
package trace

import (
    "strconv"
    "strings"

    "github.com/a07061625/gompf/mpf/api/taobao"
    "github.com/a07061625/gompf/mpf/api/taobao/logistics"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 物流流转信息查询
type traceSearch struct {
    taobao.BaseTaoBao
    tid        int      // 淘宝交易号
    spiltFlag  int      // 拆单标识 0:不拆单 1:拆单
    subTidList []string // 子订单列表
}

func (ts *traceSearch) SetTid(tid int) {
    if tid > 0 {
        ts.tid = tid
    } else {
        panic(mperr.NewLogisticsTaoBao(errorcode.LogisticsTaoBaoParam, "淘宝交易号不合法", nil))
    }
}

func (ts *traceSearch) SetSpiltFlag(spiltFlag int) {
    if (spiltFlag == 0) || (spiltFlag == 1) {
        ts.spiltFlag = spiltFlag
    } else {
        panic(mperr.NewLogisticsTaoBao(errorcode.LogisticsTaoBaoParam, "拆单标识不合法", nil))
    }
}

func (ts *traceSearch) SetSubTidList(subTidList []int) {
    ts.subTidList = make([]string, 0)
    for _, v := range subTidList {
        if v > 0 {
            ts.subTidList = append(ts.subTidList, strconv.Itoa(v))
        }
    }
}

func (ts *traceSearch) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if ts.tid <= 0 {
        panic(mperr.NewLogisticsTaoBao(errorcode.LogisticsTaoBaoParam, "淘宝交易号不能为空", nil))
    }
    if ts.spiltFlag == 1 {
        if len(ts.subTidList) == 0 {
            panic(mperr.NewLogisticsTaoBao(errorcode.LogisticsTaoBaoParam, "子订单列表不能为空", nil))
        } else if len(ts.subTidList) > 50 {
            panic(mperr.NewLogisticsTaoBao(errorcode.LogisticsTaoBaoParam, "子订单列表数量不能超过50个", nil))
        }
        ts.ReqData["sub_tid"] = strings.Join(ts.subTidList, ",")
    }
    ts.ReqData["tid"] = strconv.Itoa(ts.tid)
    ts.ReqData["is_split"] = strconv.Itoa(ts.spiltFlag)

    return ts.GetRequest()
}

func NewTraceSearch() *traceSearch {
    ts := &traceSearch{taobao.NewBaseTaoBao(), 0, 0, make([]string, 0)}
    conf := logistics.NewConfigTaoBao()
    ts.AppKey = conf.GetAppKey()
    ts.AppSecret = conf.GetAppSecret()
    ts.spiltFlag = 0
    ts.SetMethod("taobao.logistics.trace.search")
    return ts
}
