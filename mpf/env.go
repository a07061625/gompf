/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/23 0023
 * Time: 15:35
 */
package mpf

import (
    "flag"
    "log"
    "os"
    "regexp"
)

const (
    EnvTypeDev     = "dev"     // 环境类型-测试
    EnvTypeProduct = "product" // 环境类型-生产

    EnvServerTypeApi = "api" // 服务类型-API
    EnvServerTypeRpc = "rpc" // 服务类型-RPC

    GoEnvReqId            = "MP_REQ_ID"
    GoEnvDirRoot          = "MP_DIR_ROOT"
    GoEnvEnvType          = "MP_ENV_TYPE"
    GoEnvProjectTag       = "MP_PROJECT_TAG"
    GoEnvProjectModule    = "MP_PROJECT_MODULE"
    GoEnvProjectKey       = "MP_PROJECT_KEY"
    GoEnvProjectKeyModule = "MP_PROJECT_KEY_MODULE"
)

type env struct {
    dirRoot    string // 项目根目录
    dirConfigs string // 配置目录

    envType string // 环境类型

    projectTag       string // 项目标识
    projectModule    string // 项目模块
    projectKey       string // 项目代号
    projectKeyModule string // 项目模块代号

    serverHost   string
    serverPort   int
    serverDomain string
    serverType   string
}

func EnvDirRoot() string {
    return insEnv.dirRoot
}

func EnvDirConfigs() string {
    return insEnv.dirConfigs
}

func EnvType() string {
    return insEnv.envType
}

func EnvProjectTag() string {
    return insEnv.projectTag
}

func EnvProjectModule() string {
    return insEnv.projectModule
}

func EnvProjectKey() string {
    return insEnv.projectKey
}

func EnvProjectKeyModule() string {
    return insEnv.projectKeyModule
}

func EnvServerHost() string {
    return insEnv.serverHost
}

func EnvServerPort() int {
    return insEnv.serverPort
}

func EnvServerDomain() string {
    return insEnv.serverDomain
}

func EnvServerType() string {
    return insEnv.serverType
}

var (
    insEnv *env

    envType       = flag.String("mpet", EnvTypeProduct, "环境类型,只能是dev或product")
    projectTag    = flag.String("mppt", "", "项目标识,由小写字母和数字组成的3位长度字符串")
    projectModule = flag.String("mppm", "", "项目模块,由字母和数字组成的字符串")
)

func init() {
    flag.Parse()

    if (*envType != EnvTypeDev) && (*envType != EnvTypeProduct) {
        log.Fatalln("环境类型不支持")
    }
    match, _ := regexp.MatchString(`^[0-9a-z]{3}$`, *projectTag)
    if !match {
        log.Fatalln("项目标识不合法")
    }
    match, _ = regexp.MatchString(`^[0-9a-zA-Z]+$`, *projectModule)
    if !match {
        log.Fatalln("项目模块不合法")
    }

    insEnv = &env{}
    insEnv.envType = *envType
    insEnv.projectTag = *projectTag
    insEnv.projectModule = *projectModule
    insEnv.projectKey = *envType + *projectTag
    insEnv.projectKeyModule = *projectTag + *projectModule
    os.Setenv(GoEnvEnvType, insEnv.envType)
    os.Setenv(GoEnvProjectTag, insEnv.projectTag)
    os.Setenv(GoEnvProjectModule, insEnv.projectModule)
    os.Setenv(GoEnvProjectKey, insEnv.projectKey)
    os.Setenv(GoEnvProjectKeyModule, insEnv.projectKeyModule)
}
