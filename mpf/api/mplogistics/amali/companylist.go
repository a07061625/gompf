/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/16 0016
 * Time: 12:59
 */
package amali

import (
    "strconv"
    "strings"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/mplogistics"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 快递公司查询
type companyList struct {
    mplogistics.BaseAMAli
    expName string // 公司名称
    page    int    // 页数
    maxSize int    // 每页记录数
}

func (cl *companyList) SetExpName(expName string) {
    cl.ReqData["expName"] = strings.TrimSpace(expName)
}

func (cl *companyList) SetPage(page int) {
    if page > 0 {
        cl.ReqData["page"] = strconv.Itoa(page)
    } else {
        panic(mperr.NewLogisticsAMAli(errorcode.LogisticsAMAliParam, "页数必须大于0", nil))
    }
}

func (cl *companyList) SetMaxSize(maxSize int) {
    if maxSize > 0 {
        cl.ReqData["maxSize"] = strconv.Itoa(maxSize)
    } else {
        panic(mperr.NewLogisticsAMAli(errorcode.LogisticsAMAliParam, "每页记录数必须大于0", nil))
    }
}

func (cl *companyList) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    cl.ServiceUri = "/showapi_expressList?" + mpf.HTTPCreateParams(cl.ReqData, "none", 1)

    return cl.GetRequest()
}

func NewCompanyList() *companyList {
    cl := &companyList{mplogistics.NewBaseAMAli(), "", 0, 0}
    cl.ReqData["page"] = "1"
    cl.ReqData["maxSize"] = "100"
    return cl
}
