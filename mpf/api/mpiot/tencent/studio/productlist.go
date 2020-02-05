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

// 获取产品列表
type productList struct {
    mpiot.BaseTencent
    projectId int    // 项目ID
    devStatus string // 开发状态
}

func (pl *productList) SetProjectId(projectId int) {
    if projectId > 0 {
        pl.projectId = projectId
    } else {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "项目ID不合法", nil))
    }
}

func (pl *productList) SetDevStatus(devStatus string) {
    if len(devStatus) > 0 {
        pl.devStatus = devStatus
    } else {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "开发状态不合法", nil))
    }
}

func (pl *productList) SetOffset(offset int) {
    if offset >= 0 {
        pl.ReqData["Offset"] = strconv.Itoa(offset)
    } else {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "偏移量不合法", nil))
    }
}

func (pl *productList) SetLimit(limit int) {
    if limit > 0 {
        pl.ReqData["Limit"] = strconv.Itoa(limit)
    } else {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "每页个数不合法", nil))
    }
}

func (pl *productList) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if pl.projectId <= 0 {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "项目ID不能为空", nil))
    }
    if len(pl.devStatus) == 0 {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "开发状态不能为空", nil))
    }
    pl.ReqData["ProjectId"] = strconv.Itoa(pl.projectId)
    pl.ReqData["DevStatus"] = pl.devStatus

    return pl.GetRequest()
}

func NewProductList() *productList {
    pl := &productList{mpiot.NewBaseTencent(), 0, ""}
    pl.ReqData["Offset"] = "0"
    pl.ReqData["Limit"] = "10"
    pl.ReqHeader["X-TC-Action"] = "GetStudioProductList"
    return pl
}
