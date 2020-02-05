/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/2 0002
 * Time: 19:40
 */
package project

import (
    "regexp"
    "strconv"

    "github.com/a07061625/gompf/mpf/api/mpiot"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 修改项目
type projectModify struct {
    mpiot.BaseTencent
    projectId   int    // 项目ID
    projectName string // 项目名称
    projectDesc string // 项目描述
}

func (pm *projectModify) SetProjectId(projectId int) {
    if projectId > 0 {
        pm.projectId = projectId
    } else {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "项目ID不合法", nil))
    }
}

func (pm *projectModify) SetProjectName(projectName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, projectName)
    if match {
        pm.projectName = projectName
    } else {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "项目名称不合法", nil))
    }
}

func (pm *projectModify) SetProjectDesc(projectDesc string) {
    if len(projectDesc) > 0 {
        pm.projectDesc = projectDesc
    } else {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "项目描述不合法", nil))
    }
}

func (pm *projectModify) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if pm.projectId <= 0 {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "项目ID不能为空", nil))
    }
    if len(pm.projectName) == 0 {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "项目名称不能为空", nil))
    }
    if len(pm.projectDesc) == 0 {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "项目描述不能为空", nil))
    }
    pm.ReqData["ProjectId"] = strconv.Itoa(pm.projectId)
    pm.ReqData["ProjectName"] = pm.projectName
    pm.ReqData["ProjectDesc"] = pm.projectDesc

    return pm.GetRequest()
}

func NewProjectModify() *projectModify {
    pm := &projectModify{mpiot.NewBaseTencent(), 0, "", ""}
    pm.ReqHeader["X-TC-Action"] = "ModifyProject"
    return pm
}
