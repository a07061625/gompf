/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/2 0002
 * Time: 19:40
 */
package project

import (
    "regexp"

    "github.com/a07061625/gompf/mpf/api/mpiot"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 新建项目
type projectCreate struct {
    mpiot.BaseTencent
    projectName string // 项目名称
    projectDesc string // 项目描述
}

func (pc *projectCreate) SetProjectName(projectName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, projectName)
    if match {
        pc.projectName = projectName
    } else {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "项目名称不合法", nil))
    }
}

func (pc *projectCreate) SetProjectDesc(projectDesc string) {
    if len(projectDesc) > 0 {
        pc.projectDesc = projectDesc
    } else {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "项目描述不合法", nil))
    }
}

func (pc *projectCreate) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(pc.projectName) == 0 {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "项目名称不能为空", nil))
    }
    if len(pc.projectDesc) == 0 {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "项目描述不能为空", nil))
    }
    pc.ReqData["ProjectName"] = pc.projectName
    pc.ReqData["ProjectDesc"] = pc.projectDesc

    return pc.GetRequest()
}

func NewProjectCreate() *projectCreate {
    pc := &projectCreate{mpiot.NewBaseTencent(), "", ""}
    pc.ReqHeader["X-TC-Action"] = "CreateProject"
    return pc
}
