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

// 查询项目详情
type projectDescribe struct {
    mpiot.BaseTencent
    projectId int // 项目ID
}

func (pd *projectDescribe) SetProjectId(projectId int) {
    if projectId > 0 {
        pd.projectId = projectId
    } else {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "项目ID不合法", nil))
    }
}

func (pd *projectDescribe) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if pd.projectId <= 0 {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "项目ID不能为空", nil))
    }
    pd.ReqData["ProjectId"] = strconv.Itoa(pd.projectId)

    return pd.GetRequest()
}

func NewProjectDescribe() *projectDescribe {
    pd := &projectDescribe{mpiot.NewBaseTencent(), 0}
    pd.ReqHeader["X-TC-Action"] = "DescribeProject"
    return pd
}
