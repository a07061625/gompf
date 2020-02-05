/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/2 0002
 * Time: 19:40
 */
package studio

import (
    "strconv"

    "github.com/a07061625/gompf/mpf/api/mpiot"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 搜索产品
type productSearch struct {
    mpiot.BaseTencent
    projectId   int    // 项目ID
    productName string // 产品名称
}

func (ps *productSearch) SetProjectId(projectId int) {
    if projectId > 0 {
        ps.projectId = projectId
    } else {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "项目ID不合法", nil))
    }
}

func (ps *productSearch) SetProductName(productName string) {
    if len(productName) > 0 {
        ps.productName = productName
    } else {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "产品名称不合法", nil))
    }
}

func (ps *productSearch) SetOffset(offset int) {
    if offset >= 0 {
        ps.ReqData["Offset"] = strconv.Itoa(offset)
    } else {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "偏移量不合法", nil))
    }
}

func (ps *productSearch) SetLimit(limit int) {
    if limit > 0 {
        ps.ReqData["Limit"] = strconv.Itoa(limit)
    } else {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "每页个数不合法", nil))
    }
}

func (ps *productSearch) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if ps.projectId <= 0 {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "项目ID不能为空", nil))
    }
    if len(ps.productName) == 0 {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "产品名称不能为空", nil))
    }
    ps.ReqData["ProjectId"] = strconv.Itoa(ps.projectId)
    ps.ReqData["ProjectName"] = ps.productName

    return ps.GetRequest()
}

func NewProductSearch() *productSearch {
    ps := &productSearch{mpiot.NewBaseTencent(), 0, ""}
    ps.ReqData["Offset"] = "0"
    ps.ReqData["Limit"] = "10"
    ps.ReqHeader["X-TC-Action"] = "SearchStudioProduct"
    return ps
}
