/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/2 0002
 * Time: 19:40
 */
package project

import (
    "strconv"

    "github.com/a07061625/gompf/mpf/api/mpiot"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 获取项目列表
type projectList struct {
    mpiot.BaseTencent
}

func (pl *projectList) SetOffset(offset int) {
    if offset >= 0 {
        pl.ReqData["Offset"] = strconv.Itoa(offset)
    } else {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "偏移量不合法", nil))
    }
}

func (pl *projectList) SetLimit(limit int) {
    if limit > 0 {
        pl.ReqData["Limit"] = strconv.Itoa(limit)
    } else {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "每页个数不合法", nil))
    }
}

func (pl *projectList) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    return pl.GetRequest()
}

func NewProjectList() *projectList {
    pl := &projectList{mpiot.NewBaseTencent()}
    pl.ReqData["Offset"] = "0"
    pl.ReqData["Limit"] = "10"
    pl.ReqHeader["X-TC-Action"] = "GetProjectList"
    return pl
}
